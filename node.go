package main

import "errors"

// Node refers to Chord Node
type Node struct {
	identifier    []byte
	predecessor   *Node
	fingerTable   []*Node
	successorList []*Node
}

// create new Chord ring
func (node *Node) create() error {
	return errors.New("Unimplemented function create()")
}

// join Chord ring which has remoteNode inside
func (node *Node) join(remoteNode *Node) error {
	return errors.New("Unimplemented function join()")

}

// remoteNode thinks it may be the new predecessor of node
func (node *Node) notify(remoteNode *Node) {}

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
}

// called periodically - updates finger table entries
// allows new nodes to initialise their finger tables and existing nodes to incorporate new nodes into their finger tables
func (node *Node) fixFingers() {}

// called periodically - checks if predecessor has failed
func (node *Node) checkPredecessor() {
	/*
	   clear node’s predecessor pointer if n.predecessor has failed so that it can accept a new predecessor in notify()
	*/
}

// find successor of node with identifier id i.e. smallest node with identifier >= id
func (node *Node) findSuccessor(id []byte) (*Node, error) {
	return &Node{}, errors.New("Unimplemented function findSuccessor()")
}

// find closest predecessor of node with identifier id i.e. largest node with identifier < id
func (node *Node) findPredecessor(id []byte) (*Node, error) {
	/*
		searches finger table (and successor list) for most immediate predecessor of id
		if a node fails during the find successor procedure:
			after timeout, lookup proceeds by trying next best predecessor among nodes in the finger table and successor list

	*/
	return &Node{}, errors.New("Unimplemented function findPredecessor()")
}

// TransferKeys allow reassignment of keys on node join/fail ??
func (node *Node) TransferKeys(targetNode *Node, start []byte, end []byte) error {
	return errors.New("Unimplemented function TransferKeys()")
}
