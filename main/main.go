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
	kad := kademlia.CreateKademlia(&contact, &network)
	// kad := Kademlia{NewRoutingTable(contact), &network, map[KademliaID][]byte{}}
	go network.Listen(&kad)
	if !kademlia.IsBootstrap() {
		kad.JoinNetwork()
	}

	kademlia.CLIServer(&kad)
	// temp code to send pings to bootstrap
	for {
		c := kademlia.NewContact(kademlia.NewKademliaIDString(kademlia.BOOTSTRAP_ID), kademlia.BOOTSTRAP_IP)
		// msg := "test"
		// network.SendStoreReqMessage(&c, kademlia.GetHash([]byte(msg)), []byte(msg))
		rpc := network.SendPingMessage(&kad, &c)
		if rpc.Type == kademlia.PONG {
			fmt.Println("PONGED by", rpc.Sender)
		} else {
			fmt.Println("Connection timedout...")
		}
		time.Sleep(time.Second * 60)
	}

}
