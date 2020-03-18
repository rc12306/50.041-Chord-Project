package chord

// structure not developed yet: look at grpc or golang's native net/rpc??
// methods that are available for remote access must be of the form "func (t *T) MethodName(argType T1, replyType *T2) error"

// RemoteNode refers to the the structure of other Chord Nodes
type RemoteNode struct {
	Identifier int
	IP         string
}

// GetSuccessorListRPC gets successor list of remote node through RPC
func (remoteNode *RemoteNode) GetSuccessorListRPC() []*RemoteNode {
	return querySuccessorList(remoteNode.IP)
}

// GetPredecessorRPC gets predecessor of remote node through RPC
func (remoteNode *RemoteNode) GetPredecessorRPC() *RemoteNode {
	return queryPredecessor(remoteNode.IP)[0]
}

// FindSuccessorRPC finds sucessor of id of remote node through RPC
func (remoteNode *RemoteNode) FindSuccessorRPC(id int) *RemoteNode {
	// connect to remote node and ask it to run findsuccessor()
	return query(id, remoteNode.IP)
}

// NotifyRPC notifies remote node that portential thinks it may be the new predecessor of it through RPC?
func (remoteNode *RemoteNode) NotifyRPC(potential *RemoteNode) {
	notify(remoteNode.IP, potential)
}

// NoReplyRPC checks if remoteNode is alive
func (remoteNode *RemoteNode) NoReplyRPC() bool {
	return ping(remoteNode.IP)
}

// GetRPC gets file using hashed file name as key from remote node
func (remoteNode *RemoteNode) GetRPC(key string) string {
	return queryValue(remoteNode.IP, key)
}

// PutRPC puts key-value pair into remote node's hash table through RPC
func (remoteNode *RemoteNode) PutRPC(key, value string) error {
	putKeyValue(remoteNode.IP, key, value)
	// TODO: implement error handling
	return nil
}
