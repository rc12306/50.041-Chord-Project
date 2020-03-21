package chord

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const tableSize = 6
const ringSize = 64
const fingerTableUpdateRate = 50 * time.Millisecond
const checkPredUpdateRate = 100 * time.Millisecond

// Node refers to Chord Node
type Node struct {
	Identifier    int
	IP            string
	predecessor   *RemoteNode
	fingerTable   []*RemoteNode
	successorList []*RemoteNode
	stop          chan bool
	wg            sync.WaitGroup
	hashTable     map[string]string
}

// create new Chord ring
func (node *Node) create() error {
	node.predecessor = nil
	node.successorList = make([]*RemoteNode, tableSize)
	node.successorList[0] = &RemoteNode{Identifier: node.Identifier, IP: node.IP}
	node.hashTable = make(map[string]string)
	node.stop = make(chan bool)
	return nil
}

// join Chord ring which has remoteNode inside
func (node *Node) join(remoteNode *RemoteNode) error {
	node.predecessor = nil
	successor := remoteNode.FindSuccessorRPC(node.Identifier)
	node.successorList = make([]*RemoteNode, tableSize)
	node.successorList[0] = successor
	node.hashTable = make(map[string]string)
	node.updateSuccessorList(0)
	node.stop = make(chan bool)
	return nil
}

// notifies node of remote node's existence so that node can change predecessor to remoteNode
func (node *Node) notify(remoteNode *RemoteNode) {
	if node.predecessor == nil || Between(remoteNode.Identifier, node.predecessor.Identifier, node.Identifier) {
		node.predecessor = remoteNode
		node.TransferKeys(remoteNode, remoteNode.Identifier, node.Identifier)
		// if remoteNode.IP != node.IP {
		// }
	}
}

// called periodically - updates ring structure
func (node *Node) stabilise() {
	/*
		[SUCESSOR POINTER]
		asks successor for successor's predecessor p
		decides whether p should be n's successor instead (happens when node p joined the system recently)
		notifies node n's successor of p's existence so that successor can change predecessor to n (done only if the successor knows of no closer predecessor than n)
		[SUCCESSOR LIST]
		copies successor s's list, removing its last entry, and prepending s to it
		if it notices successor has failed:
			replaces it with the first live entry in its successor list
			reconciles successor list with new successor
	*/
	ticker := time.NewTicker(fingerTableUpdateRate)
	for {
		select {
		case <-ticker.C:
			// fmt.Println("Before get predecessor")
			x := node.successorList[0].GetPredecessorRPC()
			// fmt.Println("After get predecessor", x)
			if x != nil && (Between(x.Identifier, node.Identifier, node.successorList[0].Identifier) || node.IP == node.successorList[0].IP) {
				node.successorList[0] = x
			}
			node.successorList[0].NotifyRPC(&RemoteNode{IP: node.IP, Identifier: node.Identifier})
			node.updateSuccessorList(0)
		case <-node.stop:
			node.wg.Done()
			ticker.Stop()
			return
		}
	}
}

// called periodically - updates finger table entries
// allows new nodes to initialise their finger tables and existing nodes to incorporate new nodes into their finger tables
func (node *Node) fixFingers() {
	// initialisation of finger table can be improved
	node.fingerTable = make([]*RemoteNode, tableSize)
	next := 0
	ticker := time.NewTicker(fingerTableUpdateRate)
	for {
		select {
		case <-ticker.C:
			// updateFinger
			nextNode := int(math.Pow(2, float64(next)))
			closestSuccessor := node.findSuccessor((node.Identifier + nextNode) % ringSize)
			node.fingerTable[next] = closestSuccessor
			next = (next + 1) % tableSize
		case <-node.stop:
			node.wg.Done()
			ticker.Stop()
			return
		}
	}
}

// called periodically - checks if predecessor has failed
func (node *Node) checkPredecessor() {
	/*
	   clear node’s predecessor pointer if n.predecessor has failed so that it can accept a new predecessor in notify()
	*/
	ticker := time.NewTicker(checkPredUpdateRate)
	for {
		select {
		case <-ticker.C:
			if node.predecessor != nil && !node.predecessor.IsAliveRPC() {
				node.predecessor = nil
			}
		case <-node.stop:
			node.wg.Done()
			ticker.Stop()
			return
		}

	}

}

// find successor of node/key with identifier id i.e. smallest node with identifier >= id
func (node *Node) findSuccessor(id int) *RemoteNode {
	if BetweenRightIncl(id, node.Identifier, node.successorList[0].Identifier) {
		return node.successorList[0]
	}
	// get closest preceding node to id in the finger table of this node
	n, _ := node.findClosestPredecessor(id)
	if n.IP == node.IP {
		return n
	}
	return n.FindSuccessorRPC(id)

	/*
		if a node fails during the find successor procedure:
			after timeout, lookup proceeds by trying next best predecessor among nodes in the finger table and successor list
	*/
}

// find highest predecessor of node/key with identifier id i.e. largest node with identifier < id
func (node *Node) findClosestPredecessor(id int) (*RemoteNode, error) {
	/*
		searches finger table (and successor list) for most immediate predecessor of id
	*/
	closestPred := &RemoteNode{IP: node.IP, Identifier: node.Identifier}
	for i := len(node.fingerTable) - 1; i >= 0; i-- {
		fingerEntry := node.fingerTable[i]
		if fingerEntry != nil && Between(fingerEntry.Identifier, node.Identifier, id) {
			closestPred = fingerEntry
			break
		}
	}
	for i := len(node.successorList) - 1; i >= 0; i-- {
		successorListEntry := node.successorList[i]
		if successorListEntry != nil &&
			successorListEntry.Identifier > closestPred.Identifier &&
			Between(successorListEntry.Identifier, node.Identifier, id) {
			closestPred = successorListEntry
			break
		}
	}
	return closestPred, nil
}

func (node *Node) updateSuccessorList(firstLiveSuccessor int) {
	// fmt.Println(node.successorList, firstLiveSuccessor)
	if !node.successorList[firstLiveSuccessor].IsAliveRPC() {
		firstLiveSuccessor++
		if firstLiveSuccessor == len(node.successorList) {
			fmt.Println("All nodes have failed")
		} else {
			node.updateSuccessorList(firstLiveSuccessor)
		}
	} else {
		newSuccessorList := node.successorList[firstLiveSuccessor].GetSuccessorListRPC()
		copyList := make([]*RemoteNode, tableSize)

		copy(copyList[1:], newSuccessorList)
		copyList[0] = node.successorList[firstLiveSuccessor]
		// if newSuccessorList[0].IP != node.successorList[0].IP {
		// } else {
		// fmt.Println("Else")
		// copy(copyList, newSuccessorList)
		// }
		node.successorList = copyList
	}
}

// CreateNodeAndJoin helps initialise nodes and add them to the network for testing
func (node *Node) CreateNodeAndJoin(joinNode *RemoteNode) {
	if node.IP == "" {
		fmt.Println("IP of node has not been set")
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
