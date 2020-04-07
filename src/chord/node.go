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
const fingerTableUpdateRate = 500 * time.Millisecond
const checkPredUpdateRate = 1000 * time.Millisecond

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
	successor, _ := remoteNode.findSuccessorRPC(node.Identifier)
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
			x, _ := node.successorList[0].getPredecessorRPC()
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
		node.successorLock.RUnlock()
		return node.successorList[0]
	}
	node.successorLock.RUnlock()
	// get closest preceding node to id in the finger table of this node
	n, _ := node.findClosestPredecessor(id)
	if n.IP == node.IP {
		return n
	}
	successor, _ := n.findSuccessorRPC(id)
	return successor
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
	node.dataStoreLock.RLock()
	defer node.dataStoreLock.RUnlock()
	firstLiveSuccessor := node.successorList[firstLiveSuccessorIndex]
	if !firstLiveSuccessor.ping() {
		firstLiveSuccessorIndex++
		if firstLiveSuccessorIndex == len(node.successorList) {
			fmt.Println("All nodes have failed")
		} else {
			node.updateSuccessorList(firstLiveSuccessorIndex)
		}
	} else {
		newSuccessorList, _ := firstLiveSuccessor.getSuccessorListRPC()
		copyList := make([]*RemoteNode, tableSize)

		copy(copyList[1:], newSuccessorList)
		copyList[0] = firstLiveSuccessor
		/*
			Handling duplicated keys:
				1. Get primary keys of node (i.e. non-replicated keys)
				2. Delete keys from old replication nodes
				3. Add keys to new replication nodes
		*/
		// get replicated keys
		replicatedKeys := make(map[int]string)
		for key, value := range node.hashTable {
			if node.predecessor == nil || BetweenRightIncl(key, node.predecessor.Identifier, node.Identifier) {
				replicatedKeys[key] = value
			}
		}
		// check for nodes that still remain as replciation node
		repeatedReplicationNodes := make([]*RemoteNode, 0)
		newReplicationNodes := pruneList(node.IP, copyList[:replicationFactor])
		oldReplicationNodes := pruneList(node.IP, node.successorList[:replicationFactor])
		// fmt.Println(newReplicationNodes, oldReplicationNodes)
		for _, newNode := range newReplicationNodes {
			if !containsNode(newNode, node.successorList[:replicationFactor]) {
				if node != nil {
					for key, value := range replicatedKeys {
						newNode.putRPC(key, value)
					}
				}
			} else {
				repeatedReplicationNodes = append(repeatedReplicationNodes, newNode)
			}
		}
		for _, oldNode := range oldReplicationNodes {
			if !containsNode(oldNode, repeatedReplicationNodes) {
				if node != nil {
					for key := range replicatedKeys {
						oldNode.delRPC(key)
					}
				}
			}
		}
		node.successorList = copyList
	}
}

func containsNode(node *RemoteNode, nodes []*RemoteNode) bool {
	if node == nil {
		return false
	}
	for _, replicationNode := range nodes {
		if replicationNode != nil && node.IP == replicationNode.IP {
			return true
		}
	}
	return false
}

func pruneList(myIP string, nodes []*RemoteNode) []*RemoteNode {
	uniqueList := make([]*RemoteNode, 0)
	for _, replicationNode := range nodes {
		if replicationNode != nil && replicationNode.IP != myIP && !containsNode(replicationNode, uniqueList) {
			uniqueList = append(uniqueList, replicationNode)
		}
	}
	return uniqueList
}
