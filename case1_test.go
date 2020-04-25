package main

import (
	"chord/src/chord"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

// Test1 creates the chord ring structure
func Test1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	fmt.Println("Starting test 1 ...")

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
	tDelay := time.Duration(myId*50) * time.Millisecond
	fmt.Println("\nWait for ", tDelay)
	time.Sleep(tDelay)
	fmt.Println("Node", myId, "has finished sleeping!")

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
		fmt.Println("\nCreating new ring at ", myIp)
		chord.ChordNode.CreateNodeAndJoin(nil)
		fmt.Println("New ring successfully created!")
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

		fmt.Println("\nJoining existing ring at ", Ip)
		chord.ChordNode.CreateNodeAndJoin(remoteNode)
		fmt.Println("Node ", myId, " successfully joined node!")
	}

	// Update new chord ring
	ipRing, ipNot := chord.CheckRing()
	fmt.Println("\nin RING: ", ipRing)
	fmt.Println("Outside: ", ipNot)

	// to measure timing
	fDelay := time.Duration(5) * time.Second
	fmt.Println("\nWaiting for other nodes to finish ... ...", fDelay)
	time.Sleep(fDelay)

	// fmt.Println("\nJoin ring test successful! \nPress Ctrl+C to end")
	// go func() {
	// 	<-c
	// 	wg.Done()
	// }()
	// wg.Wait()

	fmt.Println("\nTest completed.")
}
