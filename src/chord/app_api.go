package chord

import (
	"errors"
	"fmt"
	"log"
	"strconv"
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
	nodeStored, err := node.findSuccessor(keyIdentifier)
	if err != nil {
		fmt.Println("Unable to find node that file", fileName, "is stored at:", err)
		return
	}
	fmt.Println(fileName, "is stored at Node", nodeStored.Identifier, "("+nodeStored.IP+")")
	fileRetrieved, err := nodeStored.getRPC(keyIdentifier)
	if err == nil {
		fmt.Println(fileRetrieved, "has been successfully retrieved")
	} else {
		// successor pointers may be wrong: pause for 1 second for stablise to correct them
		fmt.Println("Failed to retrieve file: retrying lookup")
		time.Sleep(time.Second)
		nodeStored, err := node.findSuccessor(keyIdentifier)
		if err != nil {
			fmt.Println("Retry lookup for", fileName, "was unable to find node that file is stored at:", err)
			return
		}
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
	nodeStored, err := node.findSuccessor(keyIdentifier)
	if err != nil {
		fmt.Println("Unable to find node to add file", fileName, "in:", err)
		return
	}
	fmt.Println(fileName, "to be stored at Node", nodeStored.Identifier, "("+nodeStored.IP+")")
	err = nodeStored.putRPC(keyIdentifier, fileName)
	if err == nil {
		fmt.Println(fileName, "has been successfully put into Node", nodeStored.Identifier, "("+nodeStored.IP+")")
		fmt.Println("Repliciating key across nodes for", fileName)
		successorList, err := nodeStored.getSuccessorListRPC()
		if err != nil {
			fmt.Println("Failed to replicate file among successors: Could not get successor list")
			fmt.Println(err)
			return
		}
		replicationNodes := pruneList(nodeStored.IP, successorList[:replicationFactor])
		for index, successorNode := range replicationNodes {
			replicationErr := successorNode.putRPC(keyIdentifier, fileName)
			if replicationErr == nil {
				fmt.Println("Successfully replicated file", fileName, "in successor", index, "(Node", strconv.Itoa(successorNode.Identifier)+")", "of", nodeStored.Identifier)
			} else {
				fmt.Println("Failed to replicate file", fileName, "in successor", index, "(Node", strconv.Itoa(successorNode.Identifier)+")", "of", nodeStored.Identifier)
			}
		}
	} else {
		fmt.Println(err)
		fmt.Println("Failed to add", fileName, "into the Chord ring")
	}
}

// DelFile allows user to delete file into the Chord ring
func (node *Node) DelFile(fileName string) {
	keyIdentifier := Hash(fileName)
	nodeStored, err := node.findSuccessor(keyIdentifier)
	if err != nil {
		fmt.Println("Unable to find node to delete file", fileName, "from:", err)
		return
	}
	fmt.Println(fileName, "to be deleted at Node", nodeStored.Identifier, "("+nodeStored.IP+")")
	err = nodeStored.delRPC(keyIdentifier)
	if err == nil {
		fmt.Println(fileName, "has been successfully deleted from Node", nodeStored.Identifier, "("+nodeStored.IP+")")
		fmt.Println("Deleting key across nodes for", fileName)
		successorList, err := nodeStored.getSuccessorListRPC()
		if err != nil {
			fmt.Println("Failed to delete replicated files among successors: Could not get successor list")
			fmt.Println(err)
			return
		}
		replicationNodes := pruneList(nodeStored.IP, successorList[:replicationFactor])
		for index, successorNode := range replicationNodes {
			replicationErr := successorNode.delRPC(keyIdentifier)
			if replicationErr == nil {
				fmt.Println("Successfully deleted file", fileName, "in successor", index, "(Node", strconv.Itoa(successorNode.Identifier)+")", "of", nodeStored.Identifier)
			} else {
				fmt.Println("Failed to delete", fileName, "from successor", index, "(Node", strconv.Itoa(successorNode.Identifier)+")", "of", nodeStored.Identifier)
			}
		}
	} else {
		fmt.Println(err)
		fmt.Println("Failed to delete", fileName, "from the Chord ring")
	}
}
