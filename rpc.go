package main

import "errors"

// structure not developed yet: look at grpc or golang's native net/rpc??

// GetSuccessorID gets id of sucessor of remote node through RPC?
func (node *Node) GetSuccessorID(id []byte) ([]byte, error) {
	return make([]byte, 0), errors.New("Unimplemented function SetSuccessorID()")

}

// SetPredecessorID sets id of predecessor of remote node through RPC?
func (node *Node) SetPredecessorID(id []byte) error {
	return errors.New("Unimplemented function SetPredecessorID()")
}

// SetSuccessorID sets id of sucessor of remote node through RPC?
func (node *Node) SetSuccessorID(id []byte) error {
	return errors.New("Unimplemented function SetSuccessorID()")
}

// FindSuccessor finds sucessor of id of remote node through RPC?
func (node *Node) FindSuccessor(id []byte) (*Node, error) {
	return &Node{}, errors.New("Unimplemented function FindSucessor()")
}

// Notify notifies remote node that portential thinks it may be the new predecessor of it through RPC?
func (node *Node) Notify(potential *Node) error {
	return errors.New("Unimplemented function Notify()")
}
