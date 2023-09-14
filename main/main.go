package main

import (
	"fmt"
	"kadlab/kademlia"
	"time"
)

func main() {
	fmt.Println("main starting...")

	contact := kademlia.CreateMyContact()
	network := kademlia.CreateNetwork(&contact)
	node := kademlia.CreateKademlia(&contact, &network)
	// node := Kademlia{NewRoutingTable(contact), &network, map[KademliaID][]byte{}}
	go network.Listen(&node)
	if !kademlia.IsBootstrap() {
		node.JoinNetwork()
	}

	go kademlia.CLIServer(&node)
	// temp code to send pings to bootstrap
	for {
		c := kademlia.NewContact(kademlia.NewKademliaIDString(kademlia.BOOTSTRAP_ID), kademlia.BOOTSTRAP_IP)
		msg := "test"
		network.SendStoreReqMessage(&c, kademlia.GetHash([]byte(msg)), []byte(msg))
		time.Sleep(time.Second * 30)
	}

}
