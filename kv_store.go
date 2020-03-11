package main

import "errors"

// kvStore stores key value pairs(key: hashed file name, value: file data)
type kvStore struct {
	storage map[string]string
}

// Get file using hashed file name as key from node
func Get(node *Node, key string) (string, error) {
	return "", errors.New("Unimplemented function Get()")
}

// Put file using hashed file name as key into node
func Put(node *Node, key string, value string) error {
	return errors.New("Unimplemented function Put()")
}

// Delete removes key-value pair from storage
func Delete(node *Node, key string) error {
	return errors.New("Unimplemented function Get()")
}
