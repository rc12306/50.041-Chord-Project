package chord

import (
	"errors"
	"log"

	// "net"
	"net/rpc"
	// "time"
)

// ChordNode is the node calling the node.go functions
var ChordNode *Node

/*
	Packet contains the format of how messages for sent to and fro nodes
	PacketType:	ping, pong, query, answer
	Msg:		Details of the packet type
	MsgInt:
	List: 		A list with remote node type(contains ID and IP address of node)
	senderIP:	sender IP address
*/
type Packet struct {
	PacketType string
	Msg        string
	MsgInt     int
	List       []*RemoteNode
	SenderIP   string
}

// RemoteNode contains information of other remote Chord Nodes
type RemoteNode struct {
	Identifier int
	IP         string
}

/*
	This is the format for rpc calls, an int, in this case is Listener needs to be there
	in order for the other node to call the Receive function in this node
*/
type Listener int

/*
	Handles the receival of message
	Different action for different packet type
	ping:				Ping checks if the node is alive, reply if alive
	findSuccessor:		Find succesor
	getSuccessorList:	A query to get the succesor lists from specified node
	getPredecessor:		A query to get predecessor info from specified node
	notify:				A notification to inform specified node to check it's predecessor

	getValue:			Lookup of which node the file resides in
	putKeyValue:		Put the file in the node

	default:			If the packet type does not fall in any of the catergories then ignore

*/
func (l *Listener) Receive(payload *Packet, reply *Packet) error {
	// Check what packet type it is
	switch packetType := payload.PacketType; packetType {
	case "ping":
		// fmt.Println("Receive ping from " + payload.SenderIP)
		// reply to ping
		*reply = Packet{"pong", "alive", 0, nil, ChordNode.IP}
		// *reply = Packet{"pong", "alive", nil, "127.0.0.1"}
		return nil
	case "findSuccessor":
		// fmt.Println("Receive query from " + payload.SenderIP)
		// Call node to do the search
		node := handleFindSuccessor(payload.MsgInt)
		// Form the packet for reply
		*reply = Packet{"answer", "", node.Identifier, []*RemoteNode{node}, ChordNode.IP}
		return nil
	case "getSuccessorList":
		// Call node to return succesor list
		nodes := handleGetSuccesorList()
		if nodes == nil {
			// reply with empty list
			*reply = Packet{"SuccesorList", "", 0, []*RemoteNode{}, ChordNode.IP}
		}
		var newList []*RemoteNode
		for _, v := range nodes {
			if v != nil {
				newList = append(newList, v)
			}
		}
		*reply = Packet{"SuccesorList", "", 0, newList, ChordNode.IP}
		return nil
	case "getPredecessor":
		// Call node to return it's predecessor
		node := handleGetPredecessor()
		if node == nil {
			*reply = Packet{"Predecessor", "", 0, []*RemoteNode{}, ChordNode.IP}
		} else {
			*reply = Packet{"Predecessor", "", 0, []*RemoteNode{node}, ChordNode.IP}
		}
		return nil
	case "notify":
		// Call node to make changes if necessory
		if payload.List != nil || len(payload.List) != 0 {
			handleNotify(payload.List[0])
		}
		return nil
	case "getValue":
		// Call node to get value from hashtable
		value := handleGetValue(payload.MsgInt)
		*reply = Packet{"Value", value, 0, nil, ChordNode.IP}
		return nil
	case "putKeyValue":
		// Call node to put file (and its identifier) into hashtable
		putSuccess := handlePutKeyValue(payload.MsgInt, payload.Msg)
		if putSuccess == nil {
			*reply = Packet{"Value", "Success", 0, nil, ChordNode.IP}
		} else {
			*reply = Packet{"Value", "File already exist in the table", 0, nil, ChordNode.IP}
		}
		return nil
	case "delKeyValue":
		// Call node to put file (and its identifier) into hashtable
		delSuccess := handleDelKeyValue(payload.MsgInt)
		if delSuccess == nil {
			*reply = Packet{"Value", "Success", 0, nil, ChordNode.IP}
		} else {
			*reply = Packet{"Value", "File has already been removed from hash table", 0, nil, ChordNode.IP}
		}
		return nil
	default:
		// Packet Pong, Answer, Value will enter this case
		return nil
	}

}

/*
	Sends a ping request to a node to check if it is alive
	Use by nodes
	Args:
		senderIP: 	sender IP (will it be given in node.go?)
		receiverIP:	receiver IP (must be given by node.go)
	return bool
*/
func (remoteNode *RemoteNode) ping() bool {
	// try to handshake with other node
	// fmt.Println("Ping")
	if remoteNode == nil {
		log.Printf("Remote node has not been set: unable to make RPC call")
		return false
	}
	client, err := rpc.Dial("tcp", remoteNode.IP+":8081")
	if err != nil {
		// if handshake failed then the node is not even alive
		log.Printf("Remote node has not started accepting connections: unable to make RPC call")
		return false
	}

	// Set up arguments
	payload := &Packet{"ping", "Are you alive?", 0, nil, ChordNode.IP}
	// payload := &Packet{"ping", "Are you alive?", nil, "127.0.0.1"}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Printf("Remote node closed connection abruptly: unable to complete RPC call")
		return false
	}
	// fmt.Println(reply.SenderIP + " is alive. ")
	client.Close()
	return true
}

/*
	TCP guarentees a reply so pong has not been implemented
*/
func pong() {

}

/*
	Set up query for either
		1. Node which file resides in
		2. Successor node
	Use by node
	Args:
		id:	hash of filename (key for lookup) OR
			hash of IP
	return id and ip of node who holds the file, error (if any)
*/
func (remoteNode *RemoteNode) findSuccessorRPC(id int) (*RemoteNode, error) {
	// connect to remote node and ask it to run findsuccessor()
	if remoteNode == nil {
		log.Printf("Remote node has not been set: unable to make RPC call")
		return nil, errors.New("Remote node has not been set: unable to make RPC call")
	}
	// query closest predecessor
	// get closestPred IP
	client, err := rpc.Dial("tcp", remoteNode.IP+":8081")
	if err != nil {
		log.Printf("Remote node has not started accepting connections: unable to make RPC call")
		return nil, errors.New("Remote node has not started accepting connections: unable to make RPC call")
	}

	// set up arguments
	payload := &Packet{"findSuccessor", "", id, nil, ChordNode.IP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Printf("Remote node closed connection abruptly: unable to complete RPC call")
		return nil, errors.New("Remote node closed connection abruptly: unable to complete RPC call")
	}
	client.Close()
	return reply.List[0], nil
}

func handleFindSuccessor(id int) *RemoteNode {
	// search closest successor
	closestPred := ChordNode.findSuccessor(id)
	return closestPred
}

/*
	Have not implemented answer function because it is a recursive query
	and tcp ensures a reply
*/
func answer() {

}

/*
	Set up query for getting successor list
	Use by node
	return successor list, error (if any)
*/
func (remoteNode *RemoteNode) getSuccessorListRPC() ([]*RemoteNode, error) {
	if remoteNode == nil {
		log.Printf("Remote node has not been set: unable to make RPC call")
		return nil, errors.New("Remote node has not been set: unable to make RPC call")
	}
	client, err := rpc.Dial("tcp", remoteNode.IP+":8081")
	if err != nil {
		log.Printf("Remote node has not started accepting connections: unable to make RPC call")
		return nil, errors.New("Remote node has not started accepting connections: unable to make RPC call")
	}

	// set up arguments
	payload := &Packet{"getSuccessorList", "Get successor list", 0, nil, ChordNode.IP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Printf("Remote node closed connection abruptly: unable to complete RPC call")
		return nil, errors.New("Remote node closed connection abruptly: unable to complete RPC call")
	}

	client.Close()
	if reply.List != nil {
		return reply.List, nil
	}
	return nil, errors.New("No successor list")
}

func handleGetSuccesorList() []*RemoteNode {
	// Call function in node.go
	return ChordNode.successorList
}

/*
	Set up query for getting predecessor
	Use by node
	Args:
		receiverIP:	IP of node you want to query
	return predecessor, error (if any)
*/
func (remoteNode *RemoteNode) getPredecessorRPC() (*RemoteNode, error) {
	if remoteNode == nil {
		log.Printf("Remote node has not been set: unable to make RPC call")
		return nil, errors.New("Remote node has not been set: unable to make RPC call")
	}
	client, err := rpc.Dial("tcp", remoteNode.IP+":8081")
	if err != nil {
		log.Printf("Remote node has not started accepting connections: unable to make RPC call")
		return nil, errors.New("Remote node has not started accepting connections: unable to make RPC call")
	}

	// set up arguments
	payload := &Packet{"getPredecessor", "Get predecessor", 0, nil, ChordNode.IP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Printf("Remote node closed connection abruptly: unable to complete RPC call")
		return nil, errors.New("Remote node closed connection abruptly: unable to complete RPC call")
	}

	client.Close()
	if len(reply.List) != 0 {
		return reply.List[0], nil
	}
	return nil, errors.New("No predecessors")
}

func handleGetPredecessor() *RemoteNode {
	// Call node to return predecessor
	return ChordNode.predecessor
}

/*
	Set up notification
	Use by node
	Args:
		receiverIP:	IP of node you want to notify
		potentialPred: node that might be the pred of node with receiverIP
*/
func (remoteNode *RemoteNode) notifyRPC(potentialPred *RemoteNode) {
	if remoteNode == nil {
		log.Printf("Remote node has not been set: unable to make RPC call")
		return
	}
	client, err := rpc.Dial("tcp", remoteNode.IP+":8081")
	if err != nil {
		log.Printf("Remote node has not started accepting connections: unable to make RPC call")
		return
	}

	// set up arguments
	// remoteNode := RemoteNode{ChordNode.Identifier, ChordNode.IP}
	payload := &Packet{"notify", "I am our predecessor.", 0, []*RemoteNode{potentialPred}, ChordNode.IP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Printf("Remote node closed connection abruptly: unable to complete RPC call")
		return
	}

	client.Close()
	return
}

func handleNotify(potentialPred *RemoteNode) {
	// Call node to make changes
	ChordNode.notify(potentialPred)
}

/*
	Set up getting files
	Use by node
	Args:
		key:	hashed value of file name
	return file value, error (if any)
*/
func (remoteNode *RemoteNode) getRPC(key int) (string, error) {
	if remoteNode == nil {
		log.Printf("Remote node has not been set: unable to make RPC call")
		return "", errors.New("Remote node has not been set: unable to make RPC call")
	}
	client, err := rpc.Dial("tcp", remoteNode.IP+":8081")
	if err != nil {
		log.Printf("Remote node has not started accepting connections: unable to make RPC call")
		return "", errors.New("Remote node has not started accepting connections: unable to make RPC call")
	}

	// set up arguments
	payload := &Packet{"getValue", "", key, nil, ChordNode.IP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Printf("Remote node closed connection abruptly: unable to complete RPC call")
		return "", errors.New("Remote node closed connection abruptly: unable to complete RPC call")
	}
	// Return error if key does not exist in hashtable
	if reply.Msg == "" {
		return "", errors.New("File does not exist in the table")
	}

	client.Close()
	return reply.Msg, nil
}

func handleGetValue(key int) string {
	// Call node to make changes
	value, err := ChordNode.get(key)
	if err != nil {
		// log.Println(err)
		return ""
	}
	return value
}

/*
	Set up adding files
	Use by node
	Args:
		key:	hashed value of file name
		value:	file name
	return error (if any)
*/
func (remoteNode *RemoteNode) putRPC(key int, value string) error {
	if remoteNode == nil {
		log.Printf("Remote node has not been set: unable to make RPC call")
		return errors.New("Remote node has not been set: unable to make RPC call")
	}
	client, err := rpc.Dial("tcp", remoteNode.IP+":8081")
	if err != nil {
		log.Printf("Remote node has not started accepting connections: unable to make RPC call")
		return errors.New("Remote node has not started accepting connections: unable to make RPC call")
	}

	// set up arguments
	payload := &Packet{"putKeyValue", value, key, nil, ChordNode.IP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Printf("Remote node closed connection abruptly: unable to complete RPC call")
		return errors.New("Remote node closed connection abruptly: unable to complete RPC call")
	}

	// Check that put was successful
	if reply.Msg != "Success" {
		//log.Printf("File already exists in the table")
		return errors.New("File already exists in the table")
	}

	client.Close()
	return nil
}

func handlePutKeyValue(key int, value string) error {
	// Call node to make changes
	err := ChordNode.put(key, value)
	if err != nil {
		//log.Println(err)
		return errors.New("File already exists in the table")
	}
	return nil
}

/*
	Set up adding files
	Use by node
	Args:
		key:	hashed value of file name
		value:	file name
	return error (if any)
*/
func (remoteNode *RemoteNode) delRPC(key int) error {
	if remoteNode == nil {
		log.Printf("Remote node has not been set: unable to make RPC call")
		return errors.New("Remote node has not been set: unable to make RPC call")
	}
	client, err := rpc.Dial("tcp", remoteNode.IP+":8081")
	if err != nil {
		log.Printf("Remote node has not started accepting connections: unable to make RPC call")
		return errors.New("Remote node has not started accepting connections: unable to make RPC call")
	}

	// set up arguments
	payload := &Packet{"delKeyValue", "", key, nil, ChordNode.IP}
	var reply Packet

	// and make an rpc call
	err = client.Call("Listener.Receive", payload, &reply)
	if err != nil {
		log.Printf("Remote node closed connection abruptly: unable to complete RPC call")
		return errors.New("Remote node closed connection abruptly: unable to complete RPC call")
	}

	// Check that put was successful
	if reply.Msg != "Success" {
		return errors.New("File has already been removed")
	}

	client.Close()
	return nil
}

func handleDelKeyValue(key int) error {
	// Call node to make changes
	err := ChordNode.delete(key)
	if err != nil {
		//log.Println(err)
		return errors.New("File has already been removed")
	}
	return nil
}
