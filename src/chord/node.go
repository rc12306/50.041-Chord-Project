package chord

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const tableSize = 6
const ringSize = 64
const replicationFactor = 3
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
	hashTable     map[int]string
	dataStoreLock sync.RWMutex
	successorLock sync.RWMutex
}

// create new Chord ring
func (node *Node) create() error {
	node.predecessor = nil
	node.successorList = make([]*RemoteNode, tableSize)
	node.successorList[0] = &RemoteNode{Identifier: node.Identifier, IP: node.IP}
	node.hashTable = make(map[int]string)
	node.stop = make(chan bool)
	return nil
}

// join Chord ring which has remoteNode inside
func (node *Node) join(remoteNode *RemoteNode) error {
	node.predecessor = nil
	successor := remoteNode.findSuccessorRPC(node.Identifier)
	node.successorList = make([]*RemoteNode, tableSize)
	node.successorList[0] = successor
	node.hashTable = make(map[int]string)
	node.successorLock.Lock()
	node.updateSuccessorList(0)
	node.successorLock.Unlock()
	node.stop = make(chan bool)
	return nil
}

// notifies node of remote node's existence so that node can change predecessor to remoteNode
func (node *Node) notify(remoteNode *RemoteNode) {
	if node.predecessor == nil || Between(remoteNode.Identifier, node.predecessor.Identifier, node.Identifier) {
		node.predecessor = remoteNode
		node.transferKeys(remoteNode, remoteNode.Identifier, node.Identifier)
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
			node.successorLock.Lock()
			x := node.successorList[0].getPredecessorRPC()
			if x != nil && (Between(x.Identifier, node.Identifier, node.successorList[0].Identifier) || node.IP == node.successorList[0].IP) {
				node.successorList[0] = x
			}
			if node.successorList[0].IP != node.IP {
				node.successorList[0].notifyRPC(&RemoteNode{IP: node.IP, Identifier: node.Identifier})
			}
			node.updateSuccessorList(0)
			node.successorLock.Unlock()
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
			if node.predecessor != nil && !node.predecessor.ping() {
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
	node.successorLock.RLock()
	if BetweenRightIncl(id, node.Identifier, node.successorList[0].Identifier) {
		return node.successorList[0]
	}
	node.successorLock.RUnlock()
	// get closest preceding node to id in the finger table of this node
	n, _ := node.findClosestPredecessor(id)
	if n.IP == node.IP {
		return n
	}
	return n.findSuccessorRPC(id)
}

// find highest predecessor of node/key with identifier id i.e. largest node with identifier < id
func (node *Node) findClosestPredecessor(id int) (*RemoteNode, error) {
	/*
		searches finger table for most immediate predecessor of id
	*/
	for i := len(node.fingerTable) - 1; i >= 0; i-- {
		fingerEntry := node.fingerTable[i]
		if fingerEntry != nil && Between(fingerEntry.Identifier, node.Identifier, id) {
			return fingerEntry, nil
		}
	}
	return &RemoteNode{IP: node.IP, Identifier: node.Identifier}, nil
}

// not locked: functions calling updateSuccessorList must ensure successorLock is held before calling function
func (node *Node) updateSuccessorList(firstLiveSuccessorIndex int) {
	firstLiveSuccessor := node.successorList[firstLiveSuccessorIndex]
	if !firstLiveSuccessor.ping() {
		firstLiveSuccessorIndex++
		if firstLiveSuccessorIndex == len(node.successorList) {
			fmt.Println("All nodes have failed")
		} else {
			node.updateSuccessorList(firstLiveSuccessorIndex)
		}
	} else {
		newSuccessorList := firstLiveSuccessor.getSuccessorListRPC()
		copyList := make([]*RemoteNode, tableSize)

		copy(copyList[1:], newSuccessorList)
		copyList[0] = firstLiveSuccessor
		node.successorList = copyList
	}
}
