package chord

import "errors"

// get file using hashed file name as key from node
func (node *Node) get(key string) (string, error) {
	// TODO: check if key-value pair exists
	return node.hashTable[key], nil
}

// put file using hashed file name as key into node's hashtable
func (node *Node) put(key string, value string) error {
	// TODO: check if key already exists
	node.hashTable[key] = value
	return nil
}

// delete removes key-value pair from storage
func (node *Node) delete(key string) error {
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
		keyIdentifier := hash(key)
		if BetweenLeftIncl(keyIdentifier, start, end) {
			err := targetNode.put(key, value)
			if err != nil {
				node.delete(key)
			} else {
				return errors.New("Unable to transfer key " + key)
			}
		}
	}
	return nil
}
