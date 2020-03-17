package chord

import (
	"crypto/sha1"
	"encoding/binary"
)

// Between checks if identifier is in range (a, b)
func Between(nodeX, nodeA, nodeB int) bool {
	return false
}

// BetweenRightIncl checks if identifier is in range (a, b]
func BetweenRightIncl(nodeX, nodeA, nodeB int) bool {
	return false
}

// BetweenLeftIncl checks if identifier is in range [a, b)
func BetweenLeftIncl(nodeX, nodeA, nodeB int) bool {
	return false
}

func hash(key string) int {
	hash := sha1.New()
	hash.Write([]byte(key))
	result := hash.Sum(nil)
	return int(binary.BigEndian.Uint64(result))
}
