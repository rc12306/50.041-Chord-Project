package chord

import (
	"fmt"
	"time"
)

const keySize = 8
const fingerTableUpdateRate = 50 * time.Millisecond
const checkPredUpdateRate = 100 * time.Millisecond

// Node refers to Chord Node
type Node struct {
	identifier    int
	predecessor   *Node
	fingerTable   []*Node
	successorList []*Node
	stop          chan bool
	fail          bool
	hashTable     map[string]string
}

// create new Chord ring
func (node *Node) create() error {
	node.predecessor = nil
	node.successorList = make([]*Node, keySize)
	node.successorList[0] = node
	node.hashTable = make(map[string]string)
	return nil
}

// join Chord ring which has remoteNode inside
func (node *Node) join(remoteNode *Node) error {
	// node.predecessor = nil
	successor, err := remoteNode.findSuccessor(node.identifier)
	if err != nil {
		return err
	}
	node.successorList = make([]*Node, keySize)
	node.successorList[0] = successor
	node.hashTable = make(map[string]string)
	node.updateSuccessorList(0)
	return nil
}

// notifies node of remote node's existence so that node can change predecessor to remoteNode
func (node *Node) notify(remoteNode *Node) {
	if node.predecessor == nil || Between(remoteNode.identifier, node.predecessor.identifier, node.identifier) {
		node.predecessor = remoteNode
		node.TransferKeys(remoteNode, remoteNode.identifier, node.identifier)
	}
}

// called periodically - updates ring structure
func (node *Node) stabilise() {
	/*
		[SUCESSOR POINTER]
		asks successor for successor's predecessor p
		decides whether p should be n's successor instead (happens when node p joined the system recently)
		notifies node n's successor of n's existence so that successor can change predecessor to n (done only if the successor knows of no closer predecessor than n)
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
			x := node.successorList[0].predecessor
			if x != nil && Between(x.identifier, node.identifier, node.successorList[0].identifier) {
				node.successorList[0] = x
			}
			node.successorList[0].notify(node)
			node.updateSuccessorList(0)
		case <-node.stop:
			ticker.Stop()
			return
		}
	}
}

func (node *Node) noReply() bool {
	return node.fail
}

// called periodically - updates finger table entries
// allows new nodes to initialise their finger tables and existing nodes to incorporate new nodes into their finger tables
func (node *Node) fixFingers() {
	// initialisation of finger table can be improved
	node.fingerTable = make([]*Node, keySize)
	next := 0
	ticker := time.NewTicker(fingerTableUpdateRate)
	for {
		select {
		case <-ticker.C:
			// updateFinger
			nextNode := 2 ^ (next - 1)
			node.fingerTable[next], _ = node.findSuccessor((node.identifier + nextNode) % keySize)
			next = (next + 1) % keySize
		case <-node.stop:
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
			if node.predecessor != nil && node.predecessor.noReply() {
				node.predecessor = nil
			}
		case <-node.stop:
			ticker.Stop()
			return
		}

	}

}

// find successor of node/key with identifier id i.e. smallest node with identifier >= id
func (node *Node) findSuccessor(id int) (*Node, error) {
	if BetweenRightIncl(id, node.identifier, node.successorList[0].identifier) {
		return node.successorList[0], nil
	}
	// get closest preceding node to id in the finger table of this node
	n, _ := node.findPredecessor(id)
	if n == node {
		return node, nil
	}
	return n.findSuccessor(id)

	/*
		if a node fails during the find successor procedure:
			after timeout, lookup proceeds by trying next best predecessor among nodes in the finger table and successor list
	*/
}

// find highest predecessor of node/key with identifier id i.e. largest node with identifier < id
func (node *Node) findPredecessor(id int) (*Node, error) {
	/*
		searches finger table (and successor list) for most immediate predecessor of id
	*/
	closestPred := node
	for i := len(node.fingerTable) - 1; i >= 0; i-- {
		fingerEntry := node.fingerTable[i]
		if fingerEntry != nil && Between(fingerEntry.identifier, node.identifier, id) {
			closestPred = fingerEntry
			break
		}
	}
	for i := len(node.successorList) - 1; i >= 0; i-- {
		successorListEntry := node.successorList[i]
		if successorListEntry != nil &&
			successorListEntry.identifier > closestPred.identifier &&
			Between(successorListEntry.identifier, node.identifier, id) {
			closestPred = successorListEntry
			break
		}
	}
	return closestPred, nil
}

func (node *Node) updateSuccessorList(firstLiveSuccessor int) {
	if node.successorList[firstLiveSuccessor].noReply() {
		firstLiveSuccessor++
		if firstLiveSuccessor == len(node.successorList) {
			fmt.Println("All nodes have failed")
		} else {
			node.updateSuccessorList(firstLiveSuccessor)
		}
	} else {
		newSuccessorList := node.successorList[firstLiveSuccessor].successorList
		copyList := make([]*Node, len(newSuccessorList))
		if len(newSuccessorList) > 1 {
			copy(copyList[1:], newSuccessorList)
		}
		copyList[0] = node.successorList[firstLiveSuccessor]
		node.successorList = copyList
	}
}

// CreateNodeAndJoin helps initialise nodes and add them to the network for testing
func CreateNodeAndJoin(identifier int, joinNode *Node) (newNode *Node) {
	node := Node{
		identifier: identifier,
	}
	if joinNode == nil {
		node.create()
	} else {
		node.join(joinNode)
	}
	go node.stabilise()
	go node.fixFingers()
	go node.checkPredecessor()
	return &node
}
