package main

import (
	"chord/src/chord"
	"fmt"
	"testing"
)

// Test1 creates the chord ring structure
func Test1(t *testing.T) {
	fmt.Println("Starting test ...")
	ipInChord := CheckRing()

	if len(ipInChord) > 0 {
		remoteIp := ipInChord[0]
		fmt.Println("Attempting to join chord ring via node: ", remoteIp)

		// Generate remoteNode
		remoteIpStr := fmt.Sprint(ip2Long(remoteIp))
		remoteId := chord.Hash(remoteIpStr)
		remoteNode := &chord.RemoteNode{
			Identifier: remoteId,
			IP:         remoteIp,
		}

		chord.ChordNode.CreateNodeAndJoin(remoteNode)
		fmt.Println("Successfully joined chord ring!")
	} else {
		fmt.Println("No existing chord ring!\nCreating new chord ring ...")
		chord.ChordNode.CreateNodeAndJoin(nil)
		fmt.Println("Sucessfully created chord ring!")
	}

	// // init User IP info
	// IP_0 := chord.GetOutboundIP()
	// IP_str_0 := fmt.Sprint(ip2Long(IP_0))
	// ID_0 := chord.Hash(IP_str_0)

	// // create initial node
	// // all other nodes will join to this
	// fmt.Println("Creating initial node...")
	// fmt.Println("IP: ", IP_0)
	// fmt.Println("ID: ", ID_0)
	// /*remoteNode := &chord.RemoteNode{
	// 	Identifier: ID_0,
	// 	IP:         IP_0,
	// }*/

	// fmt.Println("Listening to node ...")
	// go node_listen(IP_0)

	// // add other IPs to first node
	// _, othersIp := NetworkIP()
	// for _, s := range othersIp {
	// 	fmt.Println("Found IP in network: ", s)
	// 	IP := s
	// 	IP_str := fmt.Sprint(ip2Long(IP))
	// 	ID := chord.Hash(IP_str)
	// 	chord.ChordNode = &chord.Node{
	// 		Identifier: ID,
	// 		IP:         IP,
	// 	}
	// 	chord.ChordNode.IP = IP
	// 	chord.ChordNode.Identifier = ID

	// 	fmt.Println("IP: ", IP)
	// 	fmt.Println("ID: ", ID)

	// 	fmt.Println("Listening for nodes again ...")
	// 	go node_listen(IP)

	// 	fmt.Println("Attempting to join node ...")
	// 	chord.ChordNode.CreateNodeAndJoin(nil)
	// 	//
	// 	// fmt.Println("remoteNode is (" + remoteNode_IP + ") " + fmt.Sprint(remoteNode_ID) + ".")
	// 	// fmt.Println("Joined chord network (" + IP + ") as " + fmt.Sprint(ID) + ". ")
	//	}
}
