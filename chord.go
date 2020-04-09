package main

import (
	"bufio"
	"chord/src/chord"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strings"
)

var LISTENING_PORT int = 8081

/* --------------------------------DEPENDENCIES-----------------------*/

// LISTEN
func node_listen(hostIP string) {
	// defer wg.Done()
	// log.Println("Server started")
	addy, err := net.ResolveTCPAddr("tcp", hostIP+":8081")
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("IP address: " + addy.String())
	inbound, err := net.ListenTCP("tcp", addy)

	if err != nil {
		log.Fatal(err)
	}

	listener := new(chord.Listener)
	rpc.Register(listener)
	rpc.Accept(inbound)
	return
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

/* --------------------------------DEPENDENCIES-----------------------*/

func main() {

	// Init USER IP infos
	IP := chord.GetOutboundIP() //String of IP
	ID := chord.Hash(IP)        //Hashed decimal IP

	chord.ChordNode = &chord.Node{
		Identifier: -1,
	}

	fmt.Println(
		`
		Welcome to CHORD!
		
		Scan    s   		: Scan network for available nodes.
		Init    i   		: Create and Join.
		Print	p      		: Print node info.
		Leave 	l     		: Leave the current chord network.
		Find	f <fname>	: Find a file.
		Add     a <fname>	: Add a file.

		Create 	c     		: Create a new chord network. (Deprecated)
		Join	j <id>		: Join the chord network by specifying id. (Deprecated)

		Your IP is : ` + IP)
	fmt.Print(">>>")

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}
		inputs := strings.Split(input, " ")

		if inputs[0] != "" {
			switch inputs[0] {

			case "c": // CREATE A NODE

				if chord.ChordNode.Identifier != -1 {
					fmt.Print("Node already exists.\n>>>")
					break
				}

				fmt.Println("Deprecated function!")

				chord.ChordNode.IP = IP
				chord.ChordNode.Identifier = ID

				go node_listen(IP)
				chord.ChordNode.CreateNodeAndJoin(nil)

				fmt.Print("Created chord network (" + IP + ") as " + fmt.Sprint(ID) + ".")

				fmt.Print("\n>>>")

			case "j": // JOIN A NETWORK

				if chord.ChordNode.Identifier != -1 {
					fmt.Print("Node already exists.\n>>>")
					break
				}

				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)\n>>>")
					break
				}

				fmt.Println("Deprecated function!")

				ipSlice, _ := CheckRing()
				fmt.Println(ipSlice)

				remoteNode_IP := inputs[1]

				_, found := Find(ipSlice, remoteNode_IP)
				if !found {
					fmt.Println("IP not in ring nework.")
					fmt.Print("\n>>>")
					break
				}

				remoteNode_ID := chord.Hash(remoteNode_IP)
				remoteNode := &chord.RemoteNode{
					Identifier: remoteNode_ID,
					IP:         remoteNode_IP,
				}

				chord.ChordNode.IP = IP
				chord.ChordNode.Identifier = ID

				go node_listen(IP)
				chord.ChordNode.CreateNodeAndJoin(remoteNode)

				fmt.Println("remoteNode is (" + remoteNode_IP + ") " + fmt.Sprint(remoteNode_ID) + ".")
				fmt.Println("Joined chord network (" + IP + ") as " + fmt.Sprint(ID) + ". ")

				fmt.Print("\n>>>")

			case "s": // SCAN IP
				ipSlice, _ := CheckRing()
				fmt.Println(ipSlice)
				fmt.Print("\n>>>")

			case "i": // INITIALISE - Create Node and Join

				if chord.ChordNode.Identifier != -1 {
					fmt.Print("Node already exists.\n>>>")
					break
				}

				ipSlice, _ := CheckRing()
				fmt.Println(ipSlice)
				ringSize := len(ipSlice)

				if ringSize == 0 {
					chord.ChordNode.IP = IP
					chord.ChordNode.Identifier = ID

					go node_listen(IP)
					chord.ChordNode.CreateNodeAndJoin(nil)

					fmt.Print("Created chord network (" + IP + ") as " + fmt.Sprint(ID) + ".")

					fmt.Print("\n>>>")
					break
				}

				remoteNode_IP := ipSlice[rand.Intn(ringSize)]

				remoteNode_ID := chord.Hash(remoteNode_IP)
				remoteNode := &chord.RemoteNode{
					Identifier: remoteNode_ID,
					IP:         remoteNode_IP,
				}

				chord.ChordNode.IP = IP
				chord.ChordNode.Identifier = ID

				go node_listen(IP)
				chord.ChordNode.CreateNodeAndJoin(remoteNode)

				fmt.Println("remoteNode is (" + remoteNode_IP + ") " + fmt.Sprint(remoteNode_ID) + ".")
				fmt.Println("Joined chord network (" + IP + ") as " + fmt.Sprint(ID) + ". ")

				fmt.Print("\n>>>")

			case "p": // PRINT NODE DATA
				if chord.ChordNode.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}

				chord.ChordNode.PrintNode()
				fmt.Print("\n>>>")

			case "l": // LEAVE A NETWORK
				if chord.ChordNode.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}

				// TODO: use shutdown()
				// TODO: close node_listen goroutine()
				chord.ChordNode.ShutDown()
				chord.ChordNode = &chord.Node{}
				chord.ChordNode.Identifier = -1
				fmt.Print("Left chord network (" + IP + ") as " + fmt.Sprint(ID) + ".")
				os.Exit(1)

			case "f": // FIND A FILE BY FILENAME
				if chord.ChordNode.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}

				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)\n>>>")
					break
				}
				fmt.Print("Node (" + IP + ") " + fmt.Sprint(ID) + " searching for: \n")

				inputs = inputs[1:]
				filename := strings.Join(inputs, " ")
				fmt.Print("	" + filename)

				// TODO: use FindFile to find file with correct filename
				chord.ChordNode.FindFile(filename)

				fmt.Print("\n>>>")

			case "a": // ADD A FILE BY FILENAME
				if chord.ChordNode.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}

				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)\n>>>")
					break
				}
				fmt.Print("Node (" + IP + ") " + fmt.Sprint(ID) + " adding file: \n")

				inputs = inputs[1:]
				filename := strings.Join(inputs, " ")
				fmt.Print("	" + filename)

				// TODO: use AddFile to add file with correct filename
				chord.ChordNode.AddFile(filename)

				fmt.Print("\n>>>")

			default:
				fmt.Print("Invalid input.\n>>>")
			}
		} else {
			fmt.Print(">>>")
		}
	}
}
