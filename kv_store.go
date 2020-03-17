package main

import "errors"

// get file using hashed file name as key from node
func (node *Node) get(key string) (string, error) {
	return "", errors.New("Unimplemented function Get()")
}

// put file using hashed file name as key into node's hashtable
func (node *Node) put(key string, value string) error {
	return errors.New("Unimplemented function Put()")
}

// delete removes key-value pair from storage
func (node *Node) delete(key string) error {
	return errors.New("Unimplemented function Get()")
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
				return errors.New("Unable to transfer key", key)
			}
		}
	}
	return nil
}
