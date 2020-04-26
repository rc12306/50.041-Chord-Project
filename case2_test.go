package main

import (
	"chord/src/chord"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"testing"
	"time"
)

// only select 1 node to execute job
func IpSelect(myIp string) string {
	newIp := ""
	ipSlice := strings.Split(myIp, ".")
	ipLen := len(ipSlice)
	// fmt.Println("ipSlice:", ipSlice)
	// fmt.Println("ipLen:", ipLen)

	if ipLen > 0 {
		ipSlice[ipLen-1] = "2"
		//fmt.Println("ipSlice:", ipSlice)

		for i := 0; i < ipLen-1; i++ {
			newIp += ipSlice[i] + "."
			// fmt.Println("newIp:", newIp)
		}
		newIp += ipSlice[ipLen-1]
	}
	// fmt.Println("newIp:", newIp)
	return newIp
}

// Add files into ring by 1 node
func AddMyFiles(myIp string) {
	selIp := IpSelect(myIp)

	if myIp == selIp {
		fileSlice := []string{"a", "c", "e", "g", "i"}

		fmt.Println("\n------------------------------------------------------------------------------")
		fmt.Println("Test 2.1: Add", fileSlice, "into ring ...")
		start21 := time.Now()

		for _, file := range fileSlice {
			start22 := time.Now()
			fmt.Println("\nAdding file ", file, "...")
			chord.ChordNode.AddFile(file)
			end22 := time.Now()
			duration22 := end22.Sub(start22)
			fmt.Println("Each file took", duration22, "to add!")
		}
		end21 := time.Now()
		duration21 := end21.Sub(start21)
		fmt.Println("\nTest 2.1 COMPLETED!!!")
		fmt.Println("Duration:", duration21, "to add!")
		fmt.Println("------------------------------------------------------------------------------")
	} else {
		fmt.Println("\nWaiting for files to be added ...")
		StaticDelay(10, ".")
		fmt.Println("Done waiting!")
	}
}

// Search for files by all nodes
func SearchMyFiles() {
	fileSlice := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}

	fmt.Println("\n------------------------------------------------------------------------------")
	fmt.Println("Test 2.2: Search", fileSlice, "in ring ...")
	start21 := time.Now()

	for _, file := range fileSlice {
		start22 := time.Now()
		fmt.Println("\nSearching file ", file, "...")
		chord.ChordNode.FindFile(file)
		end22 := time.Now()
		duration22 := end22.Sub(start22)
		fmt.Println("Search took", duration22)
	}

	end21 := time.Now()
	duration21 := end21.Sub(start21)
	fmt.Println("\nTest 2.2 COMPLETED!!!")
	fmt.Println("Duration:", duration21, "to add!")
	fmt.Println("------------------------------------------------------------------------------")
}

// Test2 add& search for files
func Test2(t *testing.T) {
	fmt.Println("Starting test 2 ...")

	// create waitGroup for user defined pause
	var wg sync.WaitGroup
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	// get machine data
	myIp, myId := InitNode()

	StaticDelay(myId, "milliseconds")

	// test 1 create/join ring function
	InitRing(myIp, myId)

	StaticDelay(5, "")

	// Update new chord ring
	fmt.Println("\nAll nodes completed test 1\nChecking chord ring details ...")
	ipRing, ipNot := chord.CheckRing()
	fmt.Println("\nin RING: ", ipRing)
	fmt.Println("Outside: ", ipNot)
	chord.ChordNode.PrintNode()

	StaticDelay(1, "")

	// test 2.1 add files
	AddMyFiles(myIp)

	StaticDelay(5, "")

	// test 2.2 search files
	SearchMyFiles()

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
