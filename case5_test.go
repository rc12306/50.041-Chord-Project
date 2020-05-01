package main

import (
	"chord/src/chord"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"testing"
)

// abrupt shutdown
func MultipleDeath(myIp string) {
	deathIp1 := IpSelect(myIp, "3")
	deathIp2 := IpSelect(myIp, "4")
	deathIp3 := IpSelect(myIp, "5")

	if myIp == deathIp1 || myIp == deathIp2 || myIp == deathIp3 {
		randomDelay(0, 50)
		//chord.ChordNode.ShutDown()
		chord.ChordNode = &chord.Node{}
		chord.ChordNode.Identifier = -1
		fmt.Print("\nNode FAILED!!!\n")
	} else {
		// create waitGroup for user defined pause
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		var wg sync.WaitGroup

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
}

// Test2 add& search for files
func Test5(t *testing.T) {
	fmt.Println("Starting test 2 ...")

	// get machine data
	myIp, myId := InitNode()

	StaticDelay(myId, "milliseconds")

	// test 1 create/join ring function
	InitRing(myIp, myId)

	StaticDelay(20, "")

	// Update new chord ring
	fmt.Println("\nAll nodes completed test 1\nChecking chord ring details ...")
	ipRing, ipNot := chord.CheckRing()
	fmt.Println("\nin RING: ", ipRing)
	fmt.Println("Outside: ", ipNot)
	chord.ChordNode.PrintNode()

	StaticDelay(1, "")

	// test 2.1 add files
	AllAddFiles()

	StaticDelay(5, "")

	// test 2.2 search files
	MultipleDeath(myIp)

}
