package chord

import "errors"

// Get file using hashed file name as key from node
func (node *Node) get(hashedFile int) (string, error) {
	// TODO: check if key-value pair exists
	return node.hashTable[hashedFile], nil
}

// put file using hashed file name as key into node's hashtable
func (node *Node) put(key int, value string) error {
	// TODO: check if key already exists
	node.hashTable[key] = value
	return nil
}

// Delete removes key-value pair from storage
func (node *Node) delete(key int) error {
	_, keyExists := node.hashTable[key]
	if keyExists {
		delete(node.hashTable, key)
		return nil
	}
	return errors.New("Key already does not exist in table")
}

// TransferKeys allow reassignment of keys on node join/fail
func (node *Node) transferKeys(targetNode *RemoteNode, start int, end int) error {
	for keyIdentifier, fileName := range node.hashTable {
		if BetweenLeftIncl(keyIdentifier, start, end) {
			err := targetNode.putRPC(keyIdentifier, fileName)
			if err != nil {
				node.delete(keyIdentifier)
			} else {
				return errors.New("Unable to transfer file " + fileName)
			}
		}
	}
	return nil
}
