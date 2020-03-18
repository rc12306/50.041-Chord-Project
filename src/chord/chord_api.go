package chord

import (
	"fmt"
	"strconv"
	// "errors"
	"log"
	// need both net and net/rpc because node acts as sender and receiver
	// "net"
	"net/rpc"
	// "time"
)

/*
	packet contains the format of how messages for sent to and fro nodes
	packetType:	ping, pong, query, answer
	msg:		message from node depending on packet type
	clientIP:	client IP address
	clientIP:	client port
*/
type Packet struct {
	PacketType string
	Msg        string
	SenderIP   string
	// senderPort 	string
}

type Listener int

/*
	Handles the receival of message
	Different action for different packet type
*/
func (l *Listener) Receive(payload *Packet, reply *Packet) error {
	// Check what packet type it is
	switch packetType := payload.PacketType; packetType {
	case "ping":
		fmt.Println("Receive ping from " + payload.SenderIP)
		// reply to ping
		*reply = Packet{"pong", "alive", ""}
		return nil
	case "pong":
		fmt.Println("Receive pong from " + payload.SenderIP)
		// no reply
		// *reply = Packet{"",0,""}
		return nil
	case "query":
		fmt.Println("Receive query from " + payload.SenderIP)
		// call node file to do the search
		// node := query(payload.Msg)
		// *reply = Packet{"answer",node,""}
		return nil
	case "answer":
		fmt.Printf("File is in node %d", payload.Msg)
		// no reply
		// *reply = Packet{"",0,""}
		return nil
	default:
		return nil
	}

}

/*
	Sends a ping request to a node to check if it is alive
	Use by nodes
	Args:
		senderIP: 	sender IP (will it be given in node.go?)
		receiverIP:	receiver IP (must be given by node.go)
*/
func (node *Node) ping(senderIP string, receiverIP string) {
	// try to handshake with other node
	client, err := rpc.Dial("tcp", receiverIP+":8081")
	if err != nil {
		// if handshake failed then the node is not even alive
		log.Fatal("Dialing:", err)
	}

	// Set up arguments
	payload := &Packet{"ping", "Are you alive?", senderIP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Fatal("Connection error:", err)
	}

	fmt.Println(reply.SenderIP + " is alive. ")

}

/*
	TCP guarentees a reply so pong has not been implemented
*/
func pong() {

}

/*
	Use by node to query a file
	Args:
		id:	hash of filename (key for lookup)
	return id of node who holds the file
*/
func query(id int) int {
	// search closest successor
	closestPred, _ := FindSuccessor(id)
	return closestPred.identifier
}

func querySuccessorList() {

}

func handleQuery(id int, closestPredIP string) int {
	// query closest predecessor
	// get closestPred IP
	client, err := rpc.Dial("tcp", closestPredIP+":8081")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	// set up arguments
	payload := &Packet{"query", string(id), ".0.0.0.0"}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Fatal("Connection error:", err)
	}

	// fmt.Println("File is in node " + reply.Msg)
	ans, _ := strconv.Atoi(reply.Msg)
	return ans
}

/*
	have not implemented answer function because it is a recursive query
	and tcp ensures a reply
*/
func answer() {

}
