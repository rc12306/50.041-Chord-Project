package main

import (
	"chord/src/chord"
	"fmt"
	"time"
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
	time.Sleep(time.Second * 5)
	nodeA.PrintNode()
	nodeB.PrintNode()
	nodeC.PrintNode()
	nodeD.PrintNode()
	nodeE.PrintNode()
	nodeF.PrintNode()
	nodeG.PrintNode()
	nodeH.PrintNode()
	nodeI.PrintNode()
	nodeJ.PrintNode()
	var input string
	fmt.Scanln(&input)
}
