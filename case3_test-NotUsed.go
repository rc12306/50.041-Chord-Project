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
	"sort"
)

func ranDelay() {
	rand.Seed(time.Now().UnixNano())
	ranVal := time.Millisecond * time.Duration(rand.Intn(2000))
	fmt.Println("Delay for ", ranVal)
	time.Sleep(ranVal)
}

func initRing() ([]string, string, int) {
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
		fmt.Println("Node ", myId, " successfully joined ring!")
	}

	// Update new chord ring
	ipRing, ipNot := chord.CheckRing()
	fmt.Println("\nIn Ring: ", ipRing)
	fmt.Println("Outside Ring: ", ipNot)

	return ipRing, myIp, myId
}

func addFiles(allIpLen int, lastIp string, myIp string) bool {
	// Add files into the chord ring
	fileSlice := [5]string{"a", "b", "c", "d", "e"}
	if (lastIp == myIp) {
		fmt.Println("\nAdding files into Node ", myIp)
		for _, file := range fileSlice {
			chord.ChordNode.AddFile(file)
		}
	} else {
		time.Sleep(time.Duration(allIpLen * 5)*time.Second)
	}
	return true
}

func searchFiles() bool {
	// Search for files
	searchSlice := [5]string{"a", "b", "c", "d", "e"}
	fmt.Println("\nTesting ... \nSearching for files in the ring ...")
	for _, search := range searchSlice {
		fmt.Println("\nSearching for ", search)
		chord.ChordNode.FindFile(search)
	}
	fmt.Println("\nInitial search completed.")
	return true
}

func Test3(t *testing.T) {
	rand.Seed(time.Now().Unix())

	var wg sync.WaitGroup
	wg.Add(1)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	// Create ring
	ipRing, myIp, myId := initRing()

	// Add files to one node in chord ring (e.g. node with biggest identifier)
	allIp := append(ipRing, myIp)
	sort.Strings(allIp)

	DELAY_CONST := len(allIp) * 5
	time.Sleep(time.Duration(DELAY_CONST)*time.Second)

	addDone := addFiles(len(allIp), allIp[len(allIp)-1], myIp)

	searchDone := false
	if addDone {
		searchDone = searchFiles()
	}

	if searchDone {
		hashTable := chord.ChordNode.ReturnHash()
		// Remove node storing the file
		if (len(hashTable) != 0) {
			fmt.Println("Node will leave network.")
			chord.ChordNode.ShutDown()
			chord.ChordNode = &chord.Node{}
			chord.ChordNode.Identifier = -1
			fmt.Print("Left chord network (" + myIp + ") as " + fmt.Sprint(myId) + ".\n")
		} else { // else, search again
			time.Sleep(time.Duration(DELAY_CONST)*time.Second)
			_, othersIp := chord.NetworkIP()
			fmt.Println("IPs in network: ", othersIp)
			if (len(othersIp) == len(allIp)) {
				fmt.Println("One node has left the network!")
				time.Sleep(time.Duration(DELAY_CONST)*time.Second)
				searchSlice2 := [5]string{"a", "b", "c", "d", "e"}
				fmt.Println("\nRestart search...")
				for _, search := range searchSlice2 {
					fmt.Println("\nSearching for ", search)
					chord.ChordNode.FindFile(search)
				}
			}
		}
	}

	go func() {
		<-c
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("\nTest completed.")
}
