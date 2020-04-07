package chord

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

// get file using hashed file name as key from node
func (node *Node) get(hashedFile int) (string, error) {
	node.dataStoreLock.RLock()
	defer node.dataStoreLock.RUnlock()
	// Check if key-value pair exists
	// value is the value of hashedFile while ok is a bool
	// if hashFile does not exist, ok is False
	value, keyExists := node.hashTable[hashedFile]
	if keyExists {
		return value, nil
	} else {
		return "", errors.New("File with identifier " + strconv.Itoa(hashedFile) + " does not exist in hash table")
	}
}

// put file using hashed file name as key into node's hashtable
func (node *Node) put(key int, value string) error {
	node.dataStoreLock.Lock()
	defer node.dataStoreLock.Unlock()
	// Check if key-value pair exists
	// keyExists is a bool
	// if hashFile does not exist, keyExists is False
	_, keyExists := node.hashTable[key]
	fmt.Println(keyExists)
	if keyExists {
		return errors.New("Error putting file " + value + " into hashtable: identifier " + strconv.Itoa(key) + " already exists in hash table")
	}
	// Add file in hash table if it does not exist
	node.hashTable[key] = value
	return nil
}

// delete removes key-value pair from storage
func (node *Node) delete(key int) {
	node.dataStoreLock.Lock()
	defer node.dataStoreLock.Unlock()
	_, keyExists := node.hashTable[key]
	if keyExists {
		delete(node.hashTable, key)
	} else {
		log.Print("Key with identifier", key, "already does not exist in table")
	}
}

// TransferKeys allow reassignment of keys on node join/fail
func (node *Node) transferKeys(targetNode *RemoteNode, start int, end int) error {
	node.dataStoreLock.Lock()
	defer node.dataStoreLock.Unlock()
	keysToDelete := make([]int, 0)
	for keyIdentifier, fileName := range node.hashTable {
		if BetweenLeftIncl(keyIdentifier, start, end) {
			err := targetNode.putRPC(keyIdentifier, fileName)
			if err != nil {
				keysToDelete = append(keysToDelete, keyIdentifier)
			} else {
				return errors.New("Unable to transfer file " + fileName)
			}
		}
	}
	for key := range keysToDelete {
		delete(node.hashTable, key)
	}
	return nil
}
