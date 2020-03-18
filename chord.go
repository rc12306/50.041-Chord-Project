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
	"net/rpc"
	"strings"
)

var LISTENING_PORT int = 8081

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

// LISTEN
func node_listen(port int, quit_listen chan bool) {

	for {
		select {
		case <-quit_listen:
			return
		default:

			addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+fmt.Sprint(port))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("IP address: ")
			inbound, err := net.ListenTCP("tcp", addy)

			if err != nil {
				log.Fatal(err)
			}

			listener := new(chord.Listener)
			rpc.Register(listener)
			rpc.Accept(inbound)
			fmt.Println("Accepted")

		}
	}

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

	quit_listen := make(chan bool)

	// Main CLI
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

				if node.Identifier != -1 {
					fmt.Print("Node already exists.\n>>>")
					break
				}

				node = chord.CreateNodeAndJoin(ID, IP_str, "")
				node.PrintNode()
				fmt.Print("Created chord network (" + IP_str + ") as " + fmt.Sprint(ID) + "." + "\n>>>")

				go node_listen(LISTENING_PORT, quit_listen)

			case "p": // PRINT NODE DATA
				if node.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}

				node.PrintNode()
				fmt.Print("\n>>>")

			case "j": // JOIN A NETWORK

				if node.Identifier != -1 {
					fmt.Print("Node already exists.\n>>>")
					break
				}

				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)\n>>>")
					break
				}

				remoteNode_IP_str := inputs[1]
				remoteNode_IP := fmt.Sprint(ip2Long(remoteNode_IP_str))
				remoteNode_ID := hash(remoteNode_IP)

				node = chord.CreateNodeAndJoin(ID, IP_str, remoteNode_IP_str)

				fmt.Print("Joined chord network (" + IP_str + ") as " + fmt.Sprint(ID) + ". remoteNode is " + fmt.Sprint(remoteNode_ID) + "." + "\n>>>")

				go node_listen(LISTENING_PORT, quit_listen)

			case "l": // LEAVE A NETWORK
				if node.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}

				quit_listen <- true

				node = &newNode
				fmt.Print("Left chord network (" + IP_str + ") as " + fmt.Sprint(ID) + "." + "\n>>>")

			case "f": // FIND A FILE BY FILENAME
				if node.Identifier == -1 {
					fmt.Print("Invalid node.\n>>>")
					break
				}

				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)\n>>>")
					break
				}
				fmt.Print("Node (" + IP_str + ") " + fmt.Sprint(ID) + " searching for: \n")

				inputs = inputs[1:]
				filename := strings.Join(inputs, " ")
				filename_hash := fmt.Sprint(hash(filename))
				fmt.Print("	" + filename + " (" + filename_hash + ")")

				// Step 1: Check own hash table.
				// Step 2: Forward query. | cObtain query.
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
