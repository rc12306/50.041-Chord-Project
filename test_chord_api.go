package main

import (
	"fmt"
	"strconv"

	// "errors"
	"log"
	// need both net and net/rpc because node acts as sender and receiver
	"chord/src/chord"
	"net"
	"net/rpc"
	"sync"
)

var wg sync.WaitGroup

func server(hostIP string) {
	defer wg.Done()
	log.Println("Server started")
	addy, err := net.ResolveTCPAddr("tcp", hostIP+":8081")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("IP address: " + addy.String())
	inbound, err := net.ListenTCP("tcp", addy)

	if err != nil {
		log.Fatal(err)
	}

	listener := new(chord.Listener)
	rpc.Register(listener)
	rpc.Accept(inbound)
	return

}

func client(testNumber int, receiverIP string) {
	defer wg.Done()
	log.Println("Client started")
	switch testNumber {
	case 1:
		log.Println("Test ping localhost:8081")
		alive := strconv.FormatBool(chord.ChordNode.Ping(receiverIP))
		log.Println("Reply from localhost " + alive)
	case 2:
		log.Println("Test pong localhost:8081")
		client, err := rpc.Dial("tcp", receiverIP+":8081")
		if err != nil {
			log.Fatal(err)
		}

		newpacket2 := chord.Packet{"pong", "Yes", nil, "127.0.0.1"}
		var reply2 chord.Packet

		err = client.Call("Listener.Receive", newpacket2, &reply2)
		if err != nil {
			fmt.Println("Pong error")
			log.Fatal(err)
		}
		// Reply should be a blank packet
		log.Printf("Reply: Type: %s, msg: %s, IP: %s", reply2.PacketType, reply2.Msg, reply2.SenderIP)

		client.Close()
	case 3:
		log.Println("Test query 100 from " + receiverIP)
		node := chord.ChordNode.Query(100, receiverIP)
		log.Printf("File/ Successor is found in node %s, with IP %s", node.Identifier, node.IP)
	case 4:
		log.Println("Test answer " + receiverIP + ":8081")
		client, err := rpc.Dial("tcp", receiverIP+":8081")
		if err != nil {
			log.Fatal(err)
		}

		newpacket4 := chord.Packet{"answer", "1", []*chord.RemoteNode{{1, "1.2.3.4"}}, "127.0.0.1"}
		var reply4 chord.Packet

		err = client.Call("Listener.Receive", newpacket4, &reply4)
		if err != nil {
			log.Fatal(err)
			fmt.Println("error")
		}
		// Reply should be a blank packet
		log.Printf("Reply: Type: %s, msg: %s, IP: %s", reply4.PacketType, reply4.Msg, reply4.SenderIP)
		client.Close()
	case 5:
		log.Println("Test query predecessor from " + receiverIP)
		nodeList := chord.ChordNode.QueryPredecessor(receiverIP)
		for i, node := range nodeList {
			log.Printf("%d: Found node %s, with IP %s", i, node.Identifier, node.IP)
		}
	case 6:
		log.Println("Test query predecessor from " + receiverIP)
		nodeList := chord.ChordNode.QuerySuccessorList(receiverIP)
		for i, node := range nodeList {
			log.Printf("%d: Found node %s, with IP %s", i, node.Identifier, node.IP)
		}
	case 7:
		log.Println("Test others " + receiverIP + ":8081")
		client, err := rpc.Dial("tcp", receiverIP+":8081")
		if err != nil {
			log.Fatal(err)
		}

		newpacket5 := chord.Packet{"others", "blah", nil, "127.0.0.1"}
		var reply5 chord.Packet

		err = client.Call("Listener.Receive", newpacket5, &reply5)
		if err != nil {
			log.Fatal(err)
			fmt.Println("error")
		}

		log.Printf("Reply: Type: %s, msg: %s, IP: %s", reply5.PacketType, reply5.Msg, reply5.SenderIP)
	default:
		return
	}
	return
}

func main() {
	wg.Add(5)

	// Find IP in the network
	myip, ipslice := chord.NetworkIP()
	fmt.Println("\n My IP addr: ", myip)
	fmt.Println("Other IP in network: ", ipslice, "\n")

	id := chord.Hash(myip)
	chord.ChordNode = &chord.Node{
		Identifier: id,
		IP:         myip,
	}

	/*
		id2 := chord.Hash("")
		other := &chord.RemoteNode{
			Identifier: id2,
			IP:         "",
		}
	*/

	go server(myip)

	chord.ChordNode.CreateNodeAndJoin(nil)
	// chord.ChordNode.CreateNodeAndJoin(other)
	chord.ChordNode.PrintNode()

	/*
		chord.ChordNode.CreateNodeAndJoin(1, nil)
		fmt.Printf("Node A's id is %d \n", nodeA.Identifier)
		chord.ChordNode = nodeA

		go server(nodeA.IP)
		time.Sleep(time.Second * 2)

			go client(1, nodeA.IP)
			time.Sleep(time.Second * 2)
			go client(2, nodeA.IP)
			time.Sleep(time.Second * 2)
			go client(4, nodeA.IP)
			time.Sleep(time.Second * 2)
			go client(7, nodeA.IP)
	*/

	wg.Wait()
	fmt.Println("terminating program")
}
