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
	//"reflect" //testing
	"bytes" //ip2Long()
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
	//fmt.Println(reflect.TypeOf(hash(IP_string)).String())

	// Intro Message
	fmt.Println(
		`
		Welcome to CHORD!
		
		Create 	c     		: Create a new chord network.
		Join	j <id>		: Join the chord network by specifying id.
		Leave 	l     		: Leave the current chord network.
		Find	f <fname>	: Find a file
		
		Your IP is : ` + IP_str)

	fmt.Print(">>>")

	for {
		// input := ""
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}

		// var input string
		// fmt.Scanln(&input)
		// fmt.Println(input)

		inputs := strings.Split(input, " ")

		if inputs[0] != "" {
			switch inputs[0] {
			case "c":

				node := chord.CreateNodeAndJoin(ID, nil)
				node.PrintNode()

				message := "Created chord network (" + IP_str + ") as " + fmt.Sprint(ID) + "."
				fmt.Print(message + "\n>>>")

				// CODE for networking

			case "j":

				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)" + "\n>>>")
				} else {

					remoteNode_IP_str := inputs[1]
					remoteNode_ID := hash(fmt.Sprint(ip2Long(remoteNode_IP_str)))

					//node := chord.CreateNodeAndJoin(ID, remoteNode_ID) //Don't know the node. Only knows IP
					//node.PrintNode()

					message := "Joined chord network (" + IP_str + ") as " + fmt.Sprint(ID) + ". "
					message2 := "remoteNode is " + fmt.Sprint(remoteNode_ID) + "."
					fmt.Print(message + message2 + "\n>>>")
				}

			case "l":

				message := "Left chord network (" + IP_str + ") as " + fmt.Sprint(ID) + "."
				fmt.Print(message + "\n>>>")

			case "f":

				if len(inputs) <= 1 {
					fmt.Print("Missing Variable(s)" + "\n>>>")
				} else {
					message := "Node (" + IP_str + ") " + fmt.Sprint(ID) + " "
					inputs = inputs[1:]
					message2 := "searching for: \n"
					fmt.Print(message + message2)

					filename := strings.Join(inputs, " ")
					filename_hash := fmt.Sprint(hash(filename))
					fmt.Print("	" + filename + " (" + filename_hash + ")\n>>>")
				}
			default:
				message := "Invalid input."
				fmt.Print(message + "\n>>>")
			}
		} else {
			fmt.Print(">>>")
		}
	}
	nodeA := chord.CreateNodeAndJoin(1, nil)
	// nodeB := chord.CreateNodeAndJoin(8, nodeA)
	// nodeC := chord.CreateNodeAndJoin(14, nodeA)
	// nodeD := chord.CreateNodeAndJoin(21, nodeB)
	// nodeE := chord.CreateNodeAndJoin(32, nodeD)
	// nodeF := chord.CreateNodeAndJoin(38, nodeD)
	// nodeG := chord.CreateNodeAndJoin(42, nodeC)
	// nodeH := chord.CreateNodeAndJoin(48, nodeF)
	// nodeI := chord.CreateNodeAndJoin(51, nodeH)
	// nodeJ := chord.CreateNodeAndJoin(56, nodeB)
	// time.Sleep(time.Second * 5)
	// nodeA.PrintNode()
	// nodeB.PrintNode()
	// nodeC.PrintNode()
	// nodeD.PrintNode()
	// nodeE.PrintNode()
	// nodeF.PrintNode()
	// nodeG.PrintNode()
	// nodeH.PrintNode()
	// nodeI.PrintNode()
	// nodeJ.PrintNode()

	// fmt.Println("\nPutting key 'hello' with value 'world' into distributed table...")
	// ans := chord.Hash("hello")
	// fmt.Println("Hashed key 'hello' has identifier", ans)
	// successor, _ := nodeA.FindSuccessor(ans)
	// successor.Put("hello", "world")
	// fmt.Println("Key 'hello' with value 'world' has been saved into Node", successor.Identifier)
	// // successor.PrintNode()
	// fmt.Println("Getting value of key 'hello'...")
	// nodeOfKey, _ := nodeA.FindSuccessor(ans)
	// value, _ := nodeOfKey.Get("hello")
	// fmt.Println("Value of key 'hello' is '" + value + "'")
	// var input string
	// fmt.Scanln(&input)
}
