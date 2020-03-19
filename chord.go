package main

import (
	"bufio"
	"chord/src/chord"
	// "crypto/sha1"     //hash()
	"encoding/binary" //ip2int()
	"fmt"
	"log" //GetOutboundIP()
	"net" //GetOutboundIP()
	"os"

	//"time"
	"bytes" //ip2Long()
	//"reflect" //testing
	"net/rpc"
	"strings"
)

var LISTENING_PORT int = 8081

/* --------------------------------DEPENDENCIES-----------------------*/

// Get preferred outbound ip of this machine as net.IP
// func GetOutboundIP() string {
// 	conn, err := net.Dial("udp", "8.8.8.8:80")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()

// 	localAddr := conn.LocalAddr().(*net.UDPAddr)

// 	return localAddr.IP.String()
// }

// convert net.IP to int
func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// convert string IP to int
func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

// hash a string into int
// func hash(key string) int {
// 	hash := sha1.New()
// 	hash.Write([]byte(key))
// 	result := hash.Sum(nil)
// 	return int(binary.BigEndian.Uint64(result))
// }

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

/* --------------------------------DEPENDENCIES-----------------------*/

func main() {

	// Init USER IP infos
	IP := chord.GetOutboundIP()       //String of IP
	IP_str := fmt.Sprint(ip2Long(IP)) //String of decimal IP
	ID := chord.Hash(IP_str)          //Hashed decimal IP

	chord.ChordNode = &chord.Node{
		Identifier: -1,
		Stop:       make(chan bool),
	}

	fmt.Println(
		`
		Welcome to CHORD!
		
		Create 	c     		: Create a new chord network.
		Join	j <id>		: Join the chord network by specifying id.
		Print	p      		: Print node info.
		Leave 	l     		: Leave the current chord network.
		Find	f <fname>	: Find a file.
		
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

				chord.ChordNode.IP = IP
				chord.ChordNode.Identifier = ID
				chord.ChordNode.CreateNodeAndJoin(nil)
				fmt.Print("Created chord network (" + IP + ") as " + fmt.Sprint(ID) + ".")

				go node_listen(IP)

				fmt.Print("\n>>>")

			case "p": // PRINT NODE DATA
				if chord.ChordNode.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}

				chord.ChordNode.PrintNode()
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

				remoteNode_IP := inputs[1]                              //String of IP
				remoteNode_IP_str := fmt.Sprint(ip2Long(remoteNode_IP)) //String of decimal IP
				remoteNode_ID := chord.Hash(remoteNode_IP_str)          //Hash of decimal IP
				remoteNode := &chord.RemoteNode{
					Identifier: remoteNode_ID,
					IP:         remoteNode_IP,
				}
				chord.ChordNode.IP = IP
				chord.ChordNode.Identifier = ID
				chord.ChordNode.CreateNodeAndJoin(remoteNode)

				fmt.Println("remoteNode is (" + remoteNode_IP + ") " + fmt.Sprint(remoteNode_ID) + ".")
				fmt.Println("Joined chord network (" + IP + ") as " + fmt.Sprint(ID) + ". ")

				go node_listen(IP)

				fmt.Print("\n>>>")

			case "l": // LEAVE A NETWORK
				if chord.ChordNode.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}
				chord.ChordNode.Stop <- true
				chord.ChordNode = &chord.Node{}
				chord.ChordNode.Identifier = -1
				fmt.Print("Left chord network (" + IP + ") as " + fmt.Sprint(ID) + "." + "\n>>>")

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
				filename_hash := fmt.Sprint(chord.Hash(filename))
				fmt.Print("	" + filename + " (" + filename_hash + ")")

				// Step 1: Check own hash table.
				// Step 2: Forward query. | Obtain query.
				// (Still waiting for connection.)

				fmt.Print("\n>>>")

			default:
				fmt.Print("Invalid input.\n>>>")
			}
		} else {
			fmt.Print(">>>")
		}
	}
}
