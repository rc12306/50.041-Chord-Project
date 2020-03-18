package chord

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"strconv"
)

// Between checks if identifier is in range (a, b)
func Between(nodeX, nodeA, nodeB int) bool {
	if nodeB < nodeA && nodeX > nodeB {
		nodeB += ringSize
	} else if nodeB < nodeA && nodeX < nodeB {
		nodeA -= ringSize
	}
	return nodeX > nodeA && nodeX < nodeB
}

// BetweenRightIncl checks if identifier is in range (a, b]
func BetweenRightIncl(nodeX, nodeA, nodeB int) bool {
	if nodeB < nodeA && nodeX > nodeB {
		nodeB += ringSize
	} else if nodeB < nodeA && nodeX < nodeB {
		nodeA -= ringSize
	}
	return nodeX > nodeA && nodeX <= nodeB
}

// BetweenLeftIncl checks if identifier is in range [a, b)
func BetweenLeftIncl(nodeX, nodeA, nodeB int) bool {
	if nodeB < nodeA && nodeX > nodeB {
		nodeB += ringSize
	} else if nodeB < nodeA && nodeX < nodeB {
		nodeA -= ringSize
	}
	return nodeX >= nodeA && nodeX < nodeB
}

func hash(key string) int {
	hash := sha1.New()
	hash.Write([]byte(key))
	result := hash.Sum(nil)
	return int(binary.BigEndian.Uint64(result))
}

// PrintNode prints the node info in a formatted way
func (node *Node) PrintNode() {
	print := "=========================================================\n"
	print += "Identifier: " + strconv.Itoa(node.identifier) + "\n"
	if node.predecessor == nil {
		print += "Predecessor: nil \n"
	} else {
		print += "Predecessor: " + strconv.Itoa(node.predecessor.identifier) + "\n"
	}
	print += "Successor: " + strconv.Itoa(node.successorList[0].identifier) + "\n"
	print += "Successor List: "
	for _, successor := range node.successorList {
		if successor != nil {
			print += strconv.Itoa(successor.identifier) + ", "
		}
	}
	print += "\nFinger Table: "
	for _, finger := range node.fingerTable {
		if finger != nil {
			print += strconv.Itoa(finger.identifier) + ", "
		}
	}
	print += "\nHash Table:\n"
	for key, value := range node.hashTable {
		keyIdentifier := hash(key)
		print += "\t" + key + " (" + strconv.Itoa(keyIdentifier) + "): " + value + "\n"
	}
	print += "\n========================================================="

	fmt.Println(print)
}
