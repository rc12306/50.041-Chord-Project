package main

import (
	"net"
	"os/exec"
	// "fmt"
	// "github.com/k0kubun/pp"
)

func Hosts1(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc1(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc1(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

type Pong1 struct {
	Ip    string
	Alive bool
}

func ping1(pingChan <-chan string, pongChan chan<- Pong1) {
	for ip := range pingChan {
		_, err := exec.Command("ping", "-c1", "-t1", ip).Output()
		var alive bool
		if err != nil {
			alive = false
		} else {
			alive = true
		}
		pongChan <- Pong1{Ip: ip, Alive: alive}
	}
}

func receivePong1(pongNum int, pongChan <-chan Pong1, doneChan chan<- []Pong1) {
	var alives []Pong1
	for i := 0; i < pongNum; i++ {
		pong1 := <-pongChan
		// fmt.Println("received:", pong)
		if pong1.Alive {
			alives = append(alives, pong1)
		}
	}
	doneChan <- alives
}

// func main() {
//   hosts, _ := Hosts("172.22.0.0/24")
//   concurrentMax := 100
//   pingChan := make(chan string, concurrentMax)
//   pongChan := make(chan Pong, len(hosts))
//   doneChan := make(chan []Pong)
//
//   for i := 0; i < concurrentMax; i++ {
//     go ping(pingChan, pongChan)
//   }
//
//   go receivePong(len(hosts), pongChan, doneChan)
//
//   for _, ip := range hosts {
//     pingChan <- ip
//     // fmt.Println("sent: " + ip)
//   }
//
//   alives := <-doneChan
//   pp.Println(alives)
// }
