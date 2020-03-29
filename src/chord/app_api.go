package chord

import (
	"fmt"
	"log"
)

// CreateNodeAndJoin helps initialise nodes and add them to the network for testing
func (node *Node) CreateNodeAndJoin(joinNode *RemoteNode) {
	if node.IP == "" {
		log.Fatal("IP of node has not been set")
	} else {
		if joinNode == nil {
			node.create()
		} else {
			node.join(joinNode)
		}
		node.wg.Add(3)
		go node.stabilise()
		go node.fixFingers()
		go node.checkPredecessor()
	}
}

// ShutDown stops all functions and waits for all of them to end before returning
func (node *Node) ShutDown() {
	// telling three functions to stop
	node.stop <- true
	node.stop <- true
	node.stop <- true
	// wait for all three functions to end properly
	node.wg.Wait()
}

func (node *Node) FindFile(fileName string) {
	keyIdentifier := Hash(fileName)
	nodeStored := node.findSuccessor(keyIdentifier)
	fmt.Println(fileName, "is stored at Node", nodeStored.Identifier, "("+nodeStored.IP+")")
	fileRetrieved, err := nodeStored.getRPC(keyIdentifier)
	if err == nil {
		fmt.Println(fileRetrieved, "has been successfully retrieved")
	} else {
		// TODO: retry lookup
	}

}

func (node *Node) AddFile(fileName string) {
	keyIdentifier := Hash(fileName)
	nodeStored := node.findSuccessor(keyIdentifier)
	fmt.Println(fileName, "to be stored at Node", nodeStored.Identifier, "("+nodeStored.IP+")")
	err := nodeStored.putRPC(keyIdentifier, fileName)
	if err == nil {
		fmt.Println(fileName, "has been successfully put into Node", nodeStored.Identifier, "("+nodeStored.IP+")")
	} else {
		// TODO: print error?
	}
}
