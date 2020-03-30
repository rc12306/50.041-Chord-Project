package chord

import (
	"errors"
)

// Get file using hashed file name as key from node
func (node *Node) get(hashedFile int) (string, error) {
	// Check if key-value pair exists
	// value is the value of hashedFile while ok is a bool
	// if hashFile does not exist, ok is False
	value, keyExists := node.hashTable[hashedFile]
	if keyExists {
		return value, nil
	} else {
		return "", errors.New("File does not exist in the table")
	}
}

// put file using hashed file name as key into node's hashtable
func (node *Node) put(key int, value string) error {
	// Check if key-value pair exists
	// value is the value of hashedFile while keyExists is a bool
	// if hashFile does not exist, keyExists is False
	value, keyExists := node.hashTable[key]
	if keyExists {
		return errors.New("File already exist in the table")
	} else {
		// Add file in hash table if it does not exist
		node.hashTable[key] = value
		return nil
	}
}

// Delete removes key-value pair from storage
func (node *Node) delete(key int) error {
	_, keyExists := node.hashTable[key]
	if keyExists {
		delete(node.hashTable, key)
		return nil
	}
	return errors.New("Key already does not exist in the table")
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
