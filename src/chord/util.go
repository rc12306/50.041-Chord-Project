package chord

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"strconv"
)

// Between checks if identifier is in range (a, b)
func Between(nodeX, nodeA, nodeB int) bool {
	if nodeB < nodeA {
		nodeB += ringSize
	}
	if nodeX < nodeA {
		nodeX += ringSize
	}
	return nodeX > nodeA && nodeX < nodeB
}

// BetweenRightIncl checks if identifier is in range (a, b]
func BetweenRightIncl(nodeX, nodeA, nodeB int) bool {
	if nodeB < nodeA {
		nodeB += ringSize
	}
	if nodeX < nodeA {
		nodeX += ringSize
	}
	return nodeX > nodeA && nodeX <= nodeB
}

// BetweenLeftIncl checks if identifier is in range [a, b)
func BetweenLeftIncl(nodeX, nodeA, nodeB int) bool {
	if nodeB < nodeA {
		nodeB += ringSize
	}
	if nodeX < nodeA {
		nodeX += ringSize
	}
	return nodeX >= nodeA && nodeX < nodeB
}

// Hash provides the SHA-1 hashing required to get the identifiers for nodes and keys
func Hash(key string) int {
	hash := sha1.New()
	hash.Write([]byte(key))
	result := hash.Sum(nil)
	return int(binary.BigEndian.Uint64(result) % ringSize)
}

// PrintNode prints the node info in a formatted way
func (node *Node) PrintNode() {
	print := "=========================================================\n"
	print += "Identifier: " + strconv.Itoa(node.Identifier) + "\n"
	print += "IP: " + node.IP + "\n"
	if node.predecessor == nil {
		print += "Predecessor: nil \n"
	} else {
		print += "Predecessor: " + strconv.Itoa(node.predecessor.Identifier) + "\n"
	}
	print += "Successor: " + strconv.Itoa(node.successorList[0].Identifier) + "\n"
	print += "Successor List: "
	for _, successor := range node.successorList {
		if successor != nil {
			print += strconv.Itoa(successor.Identifier) + ", "
		}
	}
	print += "\nFinger Table: "
	for _, finger := range node.fingerTable {
		if finger != nil {
			print += strconv.Itoa(finger.Identifier) + ", "
		}
	}
	print += "\nHash Table:\n"
	for hashedFile, fileName := range node.hashTable {
		print += strconv.Itoa(hashedFile) + ": " + fileName + "\n"
	}
	print += "\n========================================================="

	fmt.Println(print)
}

func (node *Node) ReturnHash() map[int]string {
	return node.hashTable
}
