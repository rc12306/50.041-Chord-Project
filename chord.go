package main

import (
	"chord/src/chord"
	"fmt"
)

func main() {
	nodeA := chord.CreateNodeAndJoin(1, nil)
	nodeB := chord.CreateNodeAndJoin(8, nodeA)
	nodeC := chord.CreateNodeAndJoin(14, nodeA)
	nodeD := chord.CreateNodeAndJoin(21, nodeB)
	nodeE := chord.CreateNodeAndJoin(32, nodeD)
	nodeF := chord.CreateNodeAndJoin(38, nodeD)
	nodeG := chord.CreateNodeAndJoin(42, nodeC)
	nodeH := chord.CreateNodeAndJoin(48, nodeF)
	nodeI := chord.CreateNodeAndJoin(51, nodeH)
	nodeJ := chord.CreateNodeAndJoin(56, nodeB)
	fmt.Println(nodeE, nodeG, nodeI, nodeJ)

	var input string
	fmt.Scanln(&input)
}
