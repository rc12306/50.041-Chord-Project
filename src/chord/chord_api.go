package chord

import (
	"fmt"
	"errors"
	"log"
	// need both net and net/rpc because node acts as sender and receiver
	"net"
	"net/rpc"
	"time"
)

/*
	packet contains the format of how messages for sent to and fro nodes
	packetType:	ping, pong, query, answer
	msg:		message from node depending on packet type
	clientIP:	client IP address
	clientIP:	client port
*/
type Packet struct {
	packetType 	string 
	msg 		int
	senderIP 	string
	// senderPort 	string
}

type Listener int

/*
	Handles the receival of message
	Different action for different packet type
*/
func (l *Listener) Receive(payload *Packet, reply *Packet) error{
	// Check what packet type it is
	switch packetType := payload.packetType; packetType{
	case "ping":
		fmt.Println("Receive ping from " + payload.senderIP)
		// reply to ping
		*reply = Packet{"pong", 1, ""}
		return nil
	case "pong":
		fmt.Println("Receive pong from " + payload.senderIP)
		// no reply
		// *reply = Packet{"",0,""}
		return nil
	case "query":
		// call node file to do the search
		node := query(payload.msg)
		*reply = Packet{"answer",node,""}
		return nil
	case "answer":
		fmt.Printf("File is in node %d", payload.msg)
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
func ping(senderIP string, receiverIP string){
	// try to handshake with other node
	client, err := rpc.Dial("tcp", receiverIP + ":1234")
	if err != nil {
		// if handshake failed then the node is not even alive
		log.Fatal("Dialing:", err)
	}

	// Set up arguments
	payload := &Packet{"ping",0,senderIP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Fatal("Connection error:", err)
	}

	fmt.Println(reply.senderIP + " is alive. ")

}

/*
	TCP guarentees a reply so pong has not been implemented
*/
func pong(){

}

/*
	Use by node to query a file
	Args:
		id:	hash of filename (key for lookup)
	return id of node who holds the file
*/
func query(id int) int {
	// search closest successor
	closestPred, _ := findSuccessor(id)
	return closestPred.identifier
}

func handleQuery(id int, closestPredIP string) int{
	// query closest predecessor
	// get closestPred IP
	client, err := rpc.Dial("tcp", closestPredIP + ":1234")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	// set up arguments
	payload := &Packet{"query", id, senderIP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Fatal("Connection error:", err)
	}

	// fmt.Println("File is in node " + reply.msg)
	return reply.msg
}

/*
	have not implemented answer function because it is a recursive query
	and tcp ensures a reply
*/
func answer(){

}

/*
func main() {
	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:1234")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("IP address: ")
	inbound, err := net.ListenTCP("tcp", addy)
	
	if err != nil {
		log.Fatal(err)
	}

	listener := new(Listener)
	rpc.Register(listener)
	rpc.Accept(inbound)

	// Client
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}

	newpacket := Packet{"ping",0,"0.0.0.0"}
	var reply Packet

	err = client.Call("Listener.Recevie", newpacket, &reply)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Reply: Type: %s, msg: %d, IP: %s", reply.packetType, reply.msg, reply.senderIP)
}
*/