package chord

// structure not developed yet: look at grpc or golang's native net/rpc??
// methods that are available for remote access must be of the form "func (t *T) MethodName(argType T1, replyType *T2) error"

// RemoteNode refers to the the structure of other Chord Nodes
type RemoteNode struct {
	Identifier int
	IP         string
}

// GetSuccessorList gets successor list of remote node through RPC
func (remoteNode *RemoteNode) GetSuccessorList() []*RemoteNode {
	return querySuccessorList(remoteNode.IP)
}

// GetPredecessor gets predecessor of remote node through RPC
func (remoteNode *RemoteNode) GetPredecessor() *RemoteNode {
	return queryPredecessor(remoteNode.IP)[0]
}

// FindSuccessor finds sucessor of id of remote node through RPC
func (remoteNode *RemoteNode) FindSuccessor(id int) *RemoteNode {
	// connect to remote node and ask it to run findsuccessor()
	return query(id, remoteNode.IP)
}

// Notify notifies remote node that portential thinks it may be the new predecessor of it through RPC?
func (remoteNode *RemoteNode) Notify(potential *RemoteNode) {
	notify(remoteNode.IP, potential)
}

// NoReply checks if remoteNode is alive
func (remoteNode *RemoteNode) NoReply() bool {
	return false
}

// Get file using hashed file name as key from remote node
func (remoteNode *RemoteNode) Get(key string) string {
	return queryValue(remoteNode.IP, key)
}

// Put key-value pair into remote node's hash table through RPC
func (remoteNode *RemoteNode) Put(key, value string) error {
	putKeyValue(remoteNode.IP, key, value)
	// TODO: implement error handling
	return nil
}
