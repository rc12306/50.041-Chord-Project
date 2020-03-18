package chord

import "errors"

// Get file using hashed file name as key from node
func (node *Node) Get(key string) (string, error) {
	// TODO: check if key-value pair exists
	return node.hashTable[key], nil
}

// Put file using hashed file name as key into node's hashtable
func (node *Node) Put(key string, value string) error {
	// TODO: check if key already exists
	node.hashTable[key] = value
	return nil
}

// Delete removes key-value pair from storage
func (node *Node) Delete(key string) error {
	_, keyExists := node.hashTable[key]
	if keyExists {
		delete(node.hashTable, key)
		return nil
	}
	return errors.New("Key " + key + " already does not exist in table")
}

// TransferKeys allow reassignment of keys on node join/fail
func (node *Node) TransferKeys(targetNode *Node, start int, end int) error {
	for key, value := range node.hashTable {
		keyIdentifier := Hash(key)
		if BetweenLeftIncl(keyIdentifier, start, end) {
			err := targetNode.Put(key, value)
			if err != nil {
				node.Delete(key)
			} else {
				return errors.New("Unable to transfer key " + key)
			}
		}
	}
	return nil
}
