package chord

import "errors"

// structure not developed yet: look at grpc or golang's native net/rpc??
// methods that are available for remote access must be of the form "func (t *T) MethodName(argType T1, replyType *T2) error"

// RemoteNode refers to the the structure of other Chord Nodes
type RemoteNode struct {
	Identifier int
	IP         string
}

// GetSuccessorList gets successor list of remote node through RPC
func (remoteNode *RemoteNode) GetSuccessorList() []*RemoteNode {
	return make([]*RemoteNode, 0)
}

// GetPredecessor gets predecessor of remote node through RPC
func (remoteNode *RemoteNode) GetPredecessor() *RemoteNode {
	return &RemoteNode{}
}

// FindSuccessor finds sucessor of id of remote node through RPC?
func (remoteNode *RemoteNode) FindSuccessor(id int) *RemoteNode {
	// connect to remote node and ask it to run findsuccessor()
	return &RemoteNode{}
}

// Notify notifies remote node that portential thinks it may be the new predecessor of it through RPC?
func (remoteNode *RemoteNode) Notify(potential *RemoteNode) error {
	return errors.New("Unimplemented function Notify()")
}

// NoReply checks if remoteNode is alive
func (remoteNode *RemoteNode) NoReply() bool {
	return false
}

// Get file using hashed file name as key from remote node
func (remoteNode *RemoteNode) Get(key string) string {
	return ""
}

// Put key-value pair into remote node's hash table through RPC
func (remoteNode *RemoteNode) Put(key, value string) error {
	return errors.New("Unimplemented function Put()")
}
