package chord

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// CreateNodeAndJoin helps initialise nodes and add them to the network for testing
func (node *Node) CreateNodeAndJoin(joinNode *RemoteNode) error {
	if node.IP == "" {
		log.Fatal("IP of node has not been set")
		return errors.New("IP of node has not been set")
	}
	if joinNode == nil {
		node.create()
	} else {
		err := node.join(joinNode)
		if err != nil {
			return err
		}
	}
	node.wg.Add(3)
	go node.stabilise()
	go node.fixFingers()
	go node.checkPredecessor()
	return nil

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

// FindFile allows user to retrieve the file from the Chord ring
func (node *Node) FindFile(fileName string) {
	keyIdentifier := Hash(fileName)
	nodeStored := node.findSuccessor(keyIdentifier)
	fmt.Println(fileName, "is stored at Node", nodeStored.Identifier, "("+nodeStored.IP+")")
	fileRetrieved, err := nodeStored.getRPC(keyIdentifier)
	if err == nil {
		fmt.Println(fileRetrieved, "has been successfully retrieved")
	} else {
		// successor pointers may be wrong: pause for 1 second for stablise to correct them
		fmt.Println("Failed to retrieve file: retrying lookup")
		time.Sleep(time.Second)
		nodeStored := node.findSuccessor(keyIdentifier)
		fmt.Println("Retry lookup:", fileName, "is stored at Node", nodeStored.Identifier, "("+nodeStored.IP+")")
		fileRetrieved, err := nodeStored.getRPC(keyIdentifier)
		if err == nil {
			fmt.Println(fileRetrieved, "has been successfully retrieved")
		} else {
			fmt.Println("Retry lookup for", fileName, "has failed")
		}
	}
}

// AddFile allows user to add file into the Chord ring
func (node *Node) AddFile(fileName string) {
	keyIdentifier := Hash(fileName)
	nodeStored := node.findSuccessor(keyIdentifier)
	fmt.Println(fileName, "to be stored at Node", nodeStored.Identifier, "("+nodeStored.IP+")")
	err := nodeStored.putRPC(keyIdentifier, fileName)
	if err == nil {
		fmt.Println(fileName, "has been successfully put into Node", nodeStored.Identifier, "("+nodeStored.IP+")")
		fmt.Println("Repliciating key across nodes for", fileName)
		successorList, _ := nodeStored.getSuccessorListRPC()
		replicationNodes := pruneList(nodeStored.IP, successorList[:replicationFactor])
		for index, successorNode := range replicationNodes {
			replicationErr := successorNode.putRPC(keyIdentifier, fileName)
			if replicationErr == nil {
				fmt.Println("Successfully replicated file", fileName, "in successor", index, "of", nodeStored.Identifier)
			} else {
				fmt.Println("Failed to replicate file", fileName, "in successor", index, "of", nodeStored.Identifier)
			}
		}
	} else {
		fmt.Println(err)
		fmt.Println("Failed to add", fileName, "into the Chord ring")
	}
}
