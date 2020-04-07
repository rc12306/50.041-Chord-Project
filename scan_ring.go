package main

import (
	"fmt"
	"net"
	"os/exec"

	// "os/exec"

	// "errors"
	"log"
	// "net"
	"net/rpc"
	// "time"
)

// ChordNode is the node calling the node.go functions
//var ChordNode *Node

/*
	packet contains the format of how messages for sent to and fro nodes
	packetType:	ping, pong, query, answer
	msg:		Details of the packet type
	List: 		A list with remote node type(contains ID and IP address of node)
	senderIP:	sender IP address
*/
type Packet struct {
	PacketType string
	Msg        string
	MsgInt     int
	List       []*RemoteNode
	SenderIP   string
}

// RemoteNode contains information of other remote Chord Nodes
type RemoteNode struct {
	Identifier int
	IP         string
}

/*
	Functions to scan ip in network
*/
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

//http://play.golang.org/p/m8TNTtygK0
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

func NetworkIP() (string, []string) {
	// fmt.Println("Searching for IP of nodes in network ... ...")

	basicIP := GetOutboundIP()
	myIP := basicIP + "/24"
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
		// fmt.Println("sent: " + ip)
	}

	fmt.Println("Debug 1:", doneChan)

	alives := <-doneChan

	var ipSlice []string

	for _, addr := range alives {
		if addr.Ip != basicIP {
			ipSlice = append(ipSlice, addr.Ip)
		}
	}

	fmt.Println("Search completed!")

	return basicIP, ipSlice
}

/*
	node outside of chord ring (out_node) pings to the nodes in chord ring
	attempts to check IP of node inside of chord ring (in_node)
*/
func Ping(senderIP string, receiverIP string) bool {
	// try to handshake with other node
	// fmt.Println("Ping")
	client, err := rpc.Dial("tcp", receiverIP+":8081")
	if err != nil {
		// if handshake failed then the node is not even alive
		log.Println("Dialing:", err)
		return false
	}

	// Set up arguments
	payload := &Packet{"ping", "Are you alive?", 0, nil, senderIP}

	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Println("Connection error:", err)
		// fmt.Println(receiverIP, " not in chord ring")
		return false
	}
	// fmt.Println(reply.SenderIP + " is alive. ")
	client.Close()
	// fmt.Println(receiverIP, " in chord ring")
	return true
}

func CheckRing() []string {
	var ipInRing []string
	var myIp string
	var othersIp []string

	fmt.Println("Initiating IP scan ...")
	myIp, othersIp = NetworkIP()
	fmt.Println("Found IP in network: ", othersIp)

	for i := 0; i < len(othersIp); i++ {
		fmt.Println("Checking ", othersIp[i], " if in chord ring ...")
		checkIp := Ping(myIp, othersIp[i])

		if checkIp {
			ipInRing = append(ipInRing, othersIp[i])
			// fmt.Println(othersIp[i], " is in chord ring!")
		}
	}

	// fmt.Println(ipInRing, " are the IPs of nodes in chord ring!!!")

	return ipInRing
}

// func main() {
// 	ipInChordRing := CheckRing()
// 	fmt.Println(ipInChordRing)
// }
