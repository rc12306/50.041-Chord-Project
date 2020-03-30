package chord

import (
	"errors"
	"log"
	"strconv"
)

// get file using hashed file name as key from node
func (node *Node) get(hashedFile int) (string, error) {
	node.dataStoreLock.RLock()
	defer node.dataStoreLock.RUnlock()
	fileName, fileExists := node.hashTable[hashedFile]
	if !fileExists {
		return "", errors.New("File with identifier " + strconv.Itoa(hashedFile) + " does not exist in hash table")
	}
	return fileName, nil
}

// put file using hashed file name as key into node's hashtable
func (node *Node) put(key int, value string) error {
	node.dataStoreLock.Lock()
	defer node.dataStoreLock.Unlock()
	_, fileExists := node.hashTable[key]
	if fileExists {
		return errors.New("Error putting file " + value + " into hashtable: identifier " + strconv.Itoa(key) + " already exists in hash table")
	}
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
	keysToDelete := make([]int, end-start)
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
