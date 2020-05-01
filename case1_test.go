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

func StaticDelay(val int, len string) {
	if len == "milliseconds" {
		delay := time.Millisecond * time.Duration(val*10)

		fmt.Println("\nSleep for ", delay)
		time.Sleep(delay)
		fmt.Println("Node awoken!!!")
	} else {
		delay := time.Second * time.Duration(val)

		fmt.Println("\nSleep for ", delay)
		time.Sleep(delay)
		fmt.Println("Node awoken!!!")
	}
}

func GenRand(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	ranVal := rand.Intn(max-min+1) + min
	return ranVal
}

// gnerate random delay
func randomDelay(min int, max int) {
	ranVal := GenRand(min, max) * GenRand(min, max)
	fmt.Println("Delayed for ", ranVal, "ms ... ...")
	time.Sleep(time.Millisecond * time.Duration(ranVal))
}

func InitNode() (string, int) {
	// get machine data
	fmt.Println("\nGathering machine data ...")
	myIp := chord.GetOutboundIP()
	myId := chord.Hash(myIp)
	fmt.Println("IP: ", myIp)
	fmt.Println("ID: ", myId)

	// set no node
	chord.ChordNode = &chord.Node{
		Identifier: -1,
	}
	fmt.Println("Finished data gathering!!!")
	return myIp, myId
}

// Initialise & create/join ring
func InitRing(myIp string, myId int) {
	fmt.Println("\n------------------------------------------------------------------------------")
	fmt.Println("Test 1: creating/join chord ring ...")
	start1 := time.Now()

	// scan for ring
	fmt.Println("\nScanning for ring ...")
	ipInRing, _ := chord.CheckRing()
	ringSize := len(ipInRing)
	fmt.Println("Ring scan completed!\nNodes in ring: ", ipInRing, "\nRing size: ", ringSize)

	// init node
	fmt.Println("\nCreating node ...")
	chord.ChordNode = &chord.Node{
		Identifier: myId,
		IP:         myIp,
	}
	fmt.Println("\nActivating node ...")
	go node_listen(myIp)

	// create / join ring
	if ringSize == 0 {
		// Ring does NOT exists => CREATE ring
		fmt.Println("\nRing does NOT exists!\nCreating new ring at ", myIp)
		chord.ChordNode.CreateNodeAndJoin(nil)
		fmt.Println("New ring successfully created!")
	} else {
		// Ring EXISTS => JOIN ring
		fmt.Println("\nRing does EXISTS!")
		remoteIp := ipInRing[0]
		remoteId := chord.Hash(remoteIp)
		remoteNode := &chord.RemoteNode{
			Identifier: remoteId,
			IP:         remoteIp,
		}

		chord.ChordNode.IP = myIp
		chord.ChordNode.Identifier = myId

		fmt.Println("Joining ring via ", remoteId, "(", remoteIp, ")")
		chord.ChordNode.CreateNodeAndJoin(remoteNode)
		fmt.Println("Node ", myId, " successfully joined ring!")
	}

	end1 := time.Now()
	duration1 := end1.Sub(start1)
	fmt.Println("Test 1 COMPLETED!!!\nDuration ", duration1)
	fmt.Println("------------------------------------------------------------------------------\n ")
	chord.ChordNode.PrintNode()
}

// Test1 creates the chord ring structure
func Test1(t *testing.T) {
	fmt.Println("Starting test 1 ...")

	// create waitGroup for user defined pause
	var wg sync.WaitGroup
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	// get machine data
	myIp, myId := InitNode()

	StaticDelay(myId, "milliseconds")

	// test create/join ring function
	InitRing(myIp, myId)

	StaticDelay(5, "")

	// Update new chord ring
	fmt.Println("\nAll nodes completed test 1\nChecking chord ring details ...")
	ipRing, ipNot := chord.CheckRing()
	fmt.Println("\nin RING: ", ipRing)
	fmt.Println("Outside: ", ipNot, "\n ")
	chord.ChordNode.PrintNode()

	// wait for exit
	fmt.Println("\nWaiting to exit ...\nPress crtl+c to continue")
	wg.Add(1)
	go func() {
		<-c
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("\nTest completed.")
}
