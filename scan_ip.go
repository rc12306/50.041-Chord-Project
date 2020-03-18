package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
)

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	//fmt.Println(localAddr.IP)

	return localAddr.IP.String()
}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

type Pong struct {
	Ip    string
	Alive bool
}

func ping(pingChan <-chan string, pongChan chan<- Pong) {
	for ip := range pingChan {
		_, err := exec.Command("ping", "-c1", "-t1", ip).Output()
		var alive bool
		if err != nil {
			alive = false
		} else {
			alive = true
		}
		pongChan <- Pong{Ip: ip, Alive: alive}
	}
}

func receivePong(pongNum int, pongChan <-chan Pong, doneChan chan<- []Pong) {
	var alives []Pong
	for i := 0; i < pongNum; i++ {
		pong := <-pongChan
		//fmt.Println("received:", pong)

		if pong.Alive {
			alives = append(alives, pong)
		}
	}
	doneChan <- alives
}

func networkIP() (string, []string) {
	fmt.Println("Searching for IP of nodes in network ... ...")

	myIP := GetOutboundIP() + "/24"
	fmt.Println(myIP)
	hosts, _ := Hosts(myIP)
	concurrentMax := 100
	pingChan := make(chan string, concurrentMax)
	pongChan := make(chan Pong, len(hosts))
	doneChan := make(chan []Pong)

	for i := 0; i < concurrentMax; i++ {
		go ping(pingChan, pongChan)
	}

	go receivePong(len(hosts), pongChan, doneChan)

	for _, ip := range hosts {
		pingChan <- ip
		//fmt.Println("sent: " + ip)
	}

	alives := <-doneChan

	var ipSlice []string

	for _, addr := range alives {
		ipSlice = append(ipSlice, addr.Ip)
	}

	fmt.Println("Search completed!")

	return ipSlice[0], ipSlice
}

func main() {
	joinip, ipslice := networkIP()
	fmt.Println("\n", joinip)
	fmt.Println(ipslice, "\n")
}
