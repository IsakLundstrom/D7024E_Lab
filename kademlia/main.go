package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("main starting...")
	contact := CreateMyContact()
	network := Network{contact}
	node := Kademlia{*NewRoutingTable(contact), network}
	go network.Listen(&node)
	if !IsBootstrap() {
		node.JoinNetwork()
	}

	// temp code to send pings to bootstrap
	for {
		c := NewContact(NewKademliaID(BOOTSTRAP_ID), BOOTSTRAP_IP)
		network.SendPingMessage(&c)
		time.Sleep(time.Second*30)
	}

}
