package main

import (
	"fmt"
	// "errors"
	"log"
	// need both net and net/rpc because node acts as sender and receiver
	"net"
	"net/rpc"
	"time"
	"chord/src/chord"
	"sync"
	// "context"
)

var wg sync.WaitGroup

func server() {
	defer wg.Done()
	fmt.Println("Server started")
	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8081")
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

 func client() {
	defer wg.Done()
	fmt.Println("Client started")
	client, err := rpc.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}

	newpacket1 := chord.Packet{"ping","Are you alive?","0.0.0.0"}
	var reply1 chord.Packet

	err = client.Call("Listener.Receive", newpacket1, &reply1)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error")
	}

	log.Printf("Reply: Type: %s, msg: %s, IP: %s", reply1.PacketType, reply1.Msg, reply1.SenderIP)

	newpacket2 := chord.Packet{"pong","Yes","0.0.0.0"}
	var reply2 chord.Packet

	err = client.Call("Listener.Receive", newpacket2, &reply2)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error")
	}

	log.Printf("Reply: Type: %s, msg: %s, IP: %s", reply2.PacketType, reply2.Msg, reply2.SenderIP)

	newpacket3 := chord.Packet{"query","100","0.0.0.0"}
	var reply3 chord.Packet

	err = client.Call("Listener.Receive", newpacket3, &reply3)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error")
	}

	log.Printf("Reply: Type: %s, msg: %s, IP: %s", reply3.PacketType, reply3.Msg, reply3.SenderIP)

	newpacket4 := chord.Packet{"answer","1","0.0.0.0"}
	var reply4 chord.Packet

	err = client.Call("Listener.Receive", newpacket4, &reply4)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error")
	}

	log.Printf("Reply: Type: %s, msg: %s, IP: %s", reply4.PacketType, reply4.Msg, reply4.SenderIP)

	newpacket5 := chord.Packet{"others","blah","0.0.0.0"}
	var reply5 chord.Packet

	err = client.Call("Listener.Receive", newpacket5, &reply5)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error")
	}

	log.Printf("Reply: Type: %s, msg: %s, IP: %s", reply5.PacketType, reply5.Msg, reply5.SenderIP)

}

func main() {
	wg.Add(2)
	go server()
	time.Sleep(time.Second * 5)
	go client()

	wg.Wait()
	fmt.Println("terminating program")

}