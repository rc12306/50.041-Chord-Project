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
	fmt.Println("\nGathering machine data ...")
	myIp := chord.GetOutboundIP()
	myId := chord.Hash(myIp)
	fmt.Println("IP: ", myIp)
	fmt.Println("ID: ", myId)

	// Scan for IPs in network
	fmt.Println("\nScanning network ... ...")
	_, othersIp := chord.NetworkIP()
	fmt.Println("IPs in network: ", othersIp)

	// Delay according to ID of node to avoid concurrency issues
	tDelay := time.Duration(myId) * time.Second
	fmt.Println("\nWait for ", tDelay)
	time.Sleep(tDelay)

	fmt.Println("\nLooking for IPs in ring...")
	nodesInRing, _ := chord.CheckRing()
	fmt.Println("Current IPs in Ring: ", nodesInRing)

	fmt.Println("\nCreating node ...")
	chord.ChordNode = &chord.Node{
		Identifier: myId,
		IP:         myIp,
	}
	fmt.Println("\nActivating node ...")
	go node_listen(myIp)

	// Create/Join ring
	if len(nodesInRing) < 1 {
		// Chord ring NOT exists
		// Create new ring
		fmt.Println("Creating new ring at ", myIp)
		chord.ChordNode.CreateNodeAndJoin(nil)
	} else {
		// Chord ring exists
		// Join chord ring via a node in the ring
		Ip := nodesInRing[0]
		Id := chord.Hash(Ip)
		remoteNode := &chord.RemoteNode{
			Identifier: Id,
			IP:         Ip,
		}

		chord.ChordNode.IP = Ip
		chord.ChordNode.Identifier = Id

		fmt.Println("Joining existing ring at ", Ip)
		chord.ChordNode.CreateNodeAndJoin(remoteNode)
	}

	// Update new chord ring
	ipRing, ipNot := chord.CheckRing()
	fmt.Println("in RING: ", ipRing)
	fmt.Println("Outside: ", ipNot)

	// // Delay according to ID of node to avoid concurrency issues
	// eDelay := time.Duration(90-myId) * time.Second
	// fmt.Println("\nWait for ", eDelay)
	// time.Sleep(eDelay)

	// // Delay according to ensure that all nodes enters ring before test ends
	// eDelay := time.Duration(90) * time.Second
	// fmt.Println("Wait for ", eDelay)
	// time.Sleep(eDelay)

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
