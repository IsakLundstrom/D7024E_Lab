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
	// add bootstrap node to routing table
	kademlia.table.AddContact(NewContact(NewKademliaID(BOOTSTRAP_ID), BOOTSTRAP_IP))

	// lookup on itself
	kademlia.LookupContact(&kademlia.table.me)
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
