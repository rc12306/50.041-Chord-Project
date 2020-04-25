package main

import (
	"chord/src/chord"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

func ranDelay() {
	rand.Seed(time.Now().UnixNano())
	ranVal := time.Millisecond * time.Duration(rand.Intn(2000))
	fmt.Println("Delay for ", ranVal)
	time.Sleep(ranVal)
}

func Test3(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

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
	tDelay := time.Duration(myId*100) * time.Millisecond
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

	// Wait till all nodes have joint the chord ring
	eDelay := time.Duration(10) * time.Second
	fmt.Println("\nWait for ", eDelay)
	time.Sleep(eDelay)
	fmt.Println("Node", myId, "has finished sleeping!\nTest: add & search files")

	// Add files into the chord ring
	fileSlice := [5]string{"apple", "pear", "cat", "puppy", "shield"}
	fmt.Println("\nTesting ... \nAdding files into Chord Ring ...")
	for _, file := range fileSlice {
		fmt.Println("Adding: ", file)
		chord.ChordNode.AddFile(file)
	}
	// fmt.Println("Node ", myId, "successfully added files ", fileSlice, " into chord ring!!!")

	// Search for files
	searchSlice := [8]string{"apple", "irritating", "cat", "nothing", "shield", "hydra", "puppy", "pear"}
	fmt.Println("\nTesting ... \nSearching for files in the ring ...")
	for _, search := range searchSlice {
		// generates random delays b/w search
		ranDelay()
		fmt.Println("\nSearching for ", search)
		chord.ChordNode.FindFile(search)
	}
	fmt.Println("\nSearch test successful! \nPress Ctrl+C to end")

	go func() {
		<-c
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("\nTest completed.")
}
