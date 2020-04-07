package main

import (
	"chord/src/chord"
	"fmt"
	"testing"
	"time"
)

// Test1 creates the chord ring structure
func Test1(t *testing.T) {
	fmt.Println("Starting test ...")

	// init User IP info
	fmt.Println("Gathering machine data ...")
	myIp := chord.GetOutboundIP()
	myId := chord.Hash(myIp)
	fmt.Println("IP: ", myIp)
	fmt.Println("ID: ", myId)

	fmt.Println("Creating node ...")
	chord.ChordNode = &chord.Node{
		Identifier: myId,
		IP:         myIp,
	}
	go node_listen(myIp)

	// Scan for IPs in network & chord ring
	_, othersIp := chord.NetworkIP()
	fmt.Println("Found IP in network: ", othersIp)

	// Delay according to ID of node to avoid concurrency issues
	fmt.Println("Wait for ", time.Duration(myId*100)*time.Millisecond)
	time.Sleep(time.Duration(myId) * 100 * time.Millisecond)

	fmt.Println("Looking for other IPs in ring...")
	nodesInRing, _ := chord.CheckRing()
	fmt.Println("Current IPs in Ring: ", nodesInRing)
	if len(nodesInRing) < 1 {
		fmt.Println("Creating node at ", myIp)
		chord.ChordNode.CreateNodeAndJoin(nil)
	} else {
		Ip := nodesInRing[0]
		Id := chord.Hash(Ip)
		remoteNode := &chord.RemoteNode{
			Identifier: Id,
			IP:         Ip,
		}
		fmt.Println("Joining node at ", Ip)
		chord.ChordNode.CreateNodeAndJoin(remoteNode)
	}

	ipRing, ipNot := chord.CheckRing()
	fmt.Println("IPs in ring: ", ipRing)
	fmt.Println("IPs NOT in ring: ", ipNot)

	// // add remote nodes to current node
	// for _, s := range othersIp {
	// 	Ip := s
	// 	fmt.Println("Current node: ", Ip)
	//
	// 	// ignore IPs which end with 1
	// 	if string(fmt.Sprintf(Ip)[len(fmt.Sprintf(Ip))-1]) == "1" {
	// 		fmt.Println("Invalid node: ", Ip)
	// 		continue
	// 	}
	//
	// 	// break if all nodes are in ring
	// 	currentRing := CheckRing()
	// 	fmt.Println("Nodes in ring: ", currentRing)
	// 	if len(currentRing) == len(othersIp) - 1 {
	// 		break
	// 	}
	//
	// 	// ignore IPs already in ring
	// 	skipNode := false
	// 	for _, n := range currentRing {
	// 		if Ip == n {
	// 			skipNode = true
	// 			fmt.Println("Node already in ring: ", Ip)
	// 			break
	// 		}
	// 	}
	//
	// 	if skipNode == true {
	// 		continue
	// 	}
	//
	// 	IpStr := fmt.Sprint(ip2Long(Ip))
	// 	Id := chord.Hash(IpStr)
	//
	// 	fmt.Println("Creating remote node from ", Ip)
	// 	remoteNode := &chord.RemoteNode{
	// 		Identifier: Id,
	// 		IP: Ip,
	// 	}
	//
	// 	time.Sleep(time.Duration(Id) * 100 * time.Millisecond)
	// 	chord.ChordNode.CreateNodeAndJoin(remoteNode)
	// }

}
