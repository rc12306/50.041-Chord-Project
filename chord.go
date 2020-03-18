package main

import (
	"bufio"
	"chord/src/chord"
	"crypto/sha1"     //hash()
	"encoding/binary" //ip2int()
	"fmt"
	"log" //GetOutboundIP()
	"net" //GetOutboundIP()
	"os"
	//"time"
	"bytes" //ip2Long()
	//"reflect" //testing
	"strings"
)

/* --------------------------------DEPENDENCIES-----------------------*/

// Get preferred outbound ip of this machine as net.IP
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

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
func hash(key string) int {
	hash := sha1.New()
	hash.Write([]byte(key))
	result := hash.Sum(nil)
	return int(binary.BigEndian.Uint64(result))
}

/* --------------------------------DEPENDENCIES-----------------------*/

func main() {

	// Init USER IP info
	IP := GetOutboundIP()
	IP_str := GetOutboundIP().String()
	IP_int := ip2int(IP)
	ID := hash(fmt.Sprint(IP_int))

	var newNode chord.Node = chord.Node{
		Identifier: -1,
	}
	node := &newNode

	// Intro Message
	fmt.Println(
		`
		Welcome to CHORD!
		
		Create 	c     		: Create a new chord network.
		Join	j <id>		: Join the chord network by specifying id.
		Print	p      		: Print node info.
		Leave 	l     		: Leave the current chord network.
		Find	f <fname>	: Find a file
		
		Your IP is : ` + IP_str)

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

			// CREATE A NODE
			case "c":
				node = chord.CreateNodeAndJoin(ID, nil)
				node.PrintNode()

				message := "Created chord network (" + IP_str + ") as " + fmt.Sprint(ID) + "."
				fmt.Print(message + "\n>>>")
				// Step 1: Wait for connections.

			// PRINT NODE DATA
			case "p":
				if node.Identifier == -1 {
					fmt.Println("Invalid node.")
					fmt.Print("\n>>>")
					break
				}

				node.PrintNode()
				fmt.Print("\n>>>")
				// (Still waiting for connection.)

			// JOIN A NETWORK
			case "j":
				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)" + "\n>>>")
					break
				}

				remoteNode_IP_str := inputs[1]
				remoteNode_ID := hash(fmt.Sprint(ip2Long(remoteNode_IP_str)))

				//node := chord.CreateNodeAndJoin(ID, remoteNode_ID) //Don't know the node. Only knows IP
				//node.PrintNode()

				// Step 1: Request for successor from remoteNode. | Receive successor IP.
				// Step 2: Request for successor list from successor IP. | Receive successor list and predecessor.
				// Step 3: Process own information (successor list and predecessor).

				message := "Joined chord network (" + IP_str + ") as " + fmt.Sprint(ID) + ". "
				message2 := "remoteNode is " + fmt.Sprint(remoteNode_ID) + "."
				fmt.Print(message + message2 + "\n>>>")
				// Step 4: Wait for connections.

			// LEAVE A NETWORK
			case "l":
				if node.Identifier == -1 {
					fmt.Println("Invalid node.")
					fmt.Print("\n>>>")
					break
				}

				node = &newNode
				message := "Left chord network (" + IP_str + ") as " + fmt.Sprint(ID) + "."
				fmt.Print(message + "\n>>>")

			// FIND A FILE BY FILENAME
			case "f":
				if node.Identifier == -1 {
					fmt.Println("Invalid node.")
					fmt.Print("\n>>>")
					break
				}

				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)" + "\n>>>")
					break
				}
				message := "Node (" + IP_str + ") " + fmt.Sprint(ID) + " "
				inputs = inputs[1:]
				message2 := "searching for: \n"
				fmt.Print(message + message2)

				filename := strings.Join(inputs, " ")
				filename_hash := fmt.Sprint(hash(filename))
				fmt.Print("	" + filename + " (" + filename_hash + ")")

				// Step 1: Check own hash table.
				// Step 2: Forward query. | Obtain query.
				// (Still waiting for connection.)
				fmt.Print("\n>>>")

			default:
				message := "Invalid input."
				fmt.Print(message + "\n>>>")
			}
		} else {
			fmt.Print(">>>")
		}
	}
}
