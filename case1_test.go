package main

import (
	"chord/src/chord"
	"fmt"
	"testing"
	"net/rpc"
	"strings"
	"time"
)

// Test1 creates the chord ring structure
func Test1(t *testing.T) {
	fmt.Println("Starting test ...")

	// init User IP info
	fmt.Println("Gathering machine data ...")
	myIp := chord.GetOutboundIP()
	myIpStr := fmt.Sprint(ip2Long(myIp))
	myId := chord.Hash(myIpStr)
	fmt.Println("IP: ", myIp)
	fmt.Println("ID: ", myId)

	fmt.Println("Creating node ...")
	chord.ChordNode = &chord.Node{
		Identifier: myId,
		IP:         myIp,
	}
	go node_listen(myIp)
	chord.ChordNode.CreateNodeAndJoin(nil)

	// Scan for IPs in network & chord ring
	_, othersIp := NetworkIP()
	fmt.Println("Found IP in network: ", othersIp)

	for _, s := range othersIp {
		Ip := s
		if (strings.HasSuffix(fmt.Sprintf(Ip), "1")) {
			fmt.Println("Invalid node: ", Ip)
			continue
		}
		IpStr := fmt.Sprint(ip2Long(Ip))
		Id := chord.Hash(IpStr)

		fmt.Println("Creating remote node from ", Ip)
		remoteNode := &chord.RemoteNode{
			Identifier: Id,
			IP: Ip,
		}

		time.Sleep(Id * time.Millisecond)
		_, err := rpc.Dial("tcp", remoteNode.IP+":8081")
		fmt.Println("Error here: ", err)
		if err == nil {
			chord.ChordNode.CreateNodeAndJoin(remoteNode)
			fmt.Println("IP in Ring: ", CheckRing())
			break
		}
	}

	// if len(ipInChord) > 0 {
	// 	remoteIp := ipInChord[0]
	// 	fmt.Println("Attempting to join chord ring via node: ", remoteIp)
	//
	// 	// Generate remoteNode
	// 	remoteIpStr := fmt.Sprint(ip2Long(remoteIp))
	// 	remoteId := chord.Hash(remoteIpStr)
	// 	remoteNode := &chord.RemoteNode{
	// 		Identifier: remoteId,
	// 		IP:         remoteIp,
	// 	}
	//
	// 	chord.ChordNode.CreateNodeAndJoin(remoteNode)
	// 	fmt.Println("Successfully joined chord ring!")
	// } else {
	// 	fmt.Println("No existing chord ring!")
	//
	// 	fmt.Println("Creating chord ring ...")
	// 	chord.ChordNode.CreateNodeAndJoin(nil)
	//
	// 	fmt.Println("Sucessfully created chord ring!")
	// }



}
