package main

import (
  "testing"
  "fmt"
  "chord/src/chord"
)

// Test1 creates the chord ring structure
func Test1(t *testing.T) {
  // init User IP info
  IP_0 := chord.GetOutboundIP()
  IP_str_0 := fmt.Sprint(ip2Long(IP_0))
  ID_0 := chord.Hash(IP_str_0)

  // create initial node
  // all other nodes will join to this
  remoteNode := &chord.RemoteNode{
  	Identifier: ID_0,
  	IP:         IP_0,
  }
  go node_listen(IP_0)

  // add other IPs to first node
  _, othersIp := NetworkIP()
  for _, s := range othersIp {
    fmt.Println("Found IP in network: ", s)
    IP := s
    IP_str := fmt.Sprint(ip2Long(IP))
    ID := chord.Hash(IP_str)
    chord.ChordNode = &chord.Node{
      Identifier: ID,
      IP: IP,
    }
    chord.ChordNode.IP = IP
  	chord.ChordNode.Identifier = ID

  	go node_listen(IP)
		chord.ChordNode.CreateNodeAndJoin(remoteNode)
    //
		// fmt.Println("remoteNode is (" + remoteNode_IP + ") " + fmt.Sprint(remoteNode_ID) + ".")
		// fmt.Println("Joined chord network (" + IP + ") as " + fmt.Sprint(ID) + ". ")
  }
}

// TODO:
// fix binding error
