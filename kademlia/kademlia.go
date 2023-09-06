package main

import (
	"fmt"
)

type Kademlia struct {
	table 	RoutingTable
	network Network
}

func (kademlia *Kademlia) JoinNetwork() {
	fmt.Println("TODO Joining...")
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
