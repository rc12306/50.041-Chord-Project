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
	// if hashFile does not exist, keyExists is False
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
	_, keyExists := node.hashTable[key]
	if keyExists {
		return errors.New("Error putting file " + value + " into hashtable: identifier " + strconv.Itoa(key) + " already exists in hash table")
	}
	// Add file in hash table if it does not exist
	node.hashTable[key] = value
	return nil
}

// delete removes key-value pair from storage
func (node *Node) delete(key int) error {
	node.dataStoreLock.Lock()
	defer node.dataStoreLock.Unlock()
	_, keyExists := node.hashTable[key]
	if keyExists {
		delete(node.hashTable, key)
		return nil
	} else {
		log.Print("Key with identifier" + strconv.Itoa(key) + "already does not exist in table")
		return errors.New("Key with identifier" + strconv.Itoa(key) + "does not exist in table")
	}
}

// check if key exists
func (node *Node) keyExists(key int) bool {
	_, keyExists := node.hashTable[key]
	return keyExists
}

// TransferKeys allow reassignment of keys on node join/fail for keys in (start, end]
func (node *Node) transferKeys(targetNode *RemoteNode, start int, end int) error {
	node.dataStoreLock.Lock()
	defer node.dataStoreLock.Unlock()
	keysToDelete := make([]int, 0)
	for keyIdentifier, fileName := range node.hashTable {
		if BetweenRightIncl(keyIdentifier, start, end) {
			fmt.Println("Transferring key", keyIdentifier, "to Node", targetNode.Identifier)
			// keyExists, err := targetNode.keyExistsRPC(keyIdentifier)
			// if err != nil {
			// return err
			// } else if !keyExists {
			err := targetNode.putRPC(keyIdentifier, fileName)
			if err != nil {
				return errors.New("Unable to transfer file " + fileName + " to Node " + strconv.Itoa(targetNode.Identifier))
			}
			// }
			keysToDelete = append(keysToDelete, keyIdentifier)
		}
	}
	for _, key := range keysToDelete {
		fmt.Println("Deleting key", key)
		delete(node.hashTable, key)
	}
	return nil
}
