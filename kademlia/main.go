package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("main starting...")
	contact := CreateMyContact()
	network := Network{&contact}
	node := Kademlia{NewRoutingTable(contact), &network, map[KademliaID][]byte{}}
	go network.Listen(&node)
	if !IsBootstrap() {
		node.JoinNetwork()
	}

	// temp code to send pings to bootstrap
	for {
		c := NewContact(NewKademliaIDString(BOOTSTRAP_ID), BOOTSTRAP_IP)
		msg := "test"
		network.SendStoreReqMessage(&c, GetHash([]byte(msg)), []byte(msg))
		time.Sleep(time.Second*30)
	}

}
