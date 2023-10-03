package main

import (
	"fmt"
	"kadlab/kademlia"
	"log"
	"time"
)

func main() {
	fmt.Println("main starting...")

	contact := kademlia.CreateMyContact(kademlia.IP_PREFIX)
	network := kademlia.CreateNetwork(&contact)
	kad := kademlia.CreateKademlia(&network)

	fmt.Println("My contact:", contact.String())

	isBoot, err := kademlia.IsBootstrap(kademlia.IP_PREFIX)
	if err != nil {
		log.Fatal(err)
		return
	}
	if !isBoot {
		kad.JoinNetwork()
	}

	go kademlia.StartAPI(&kad)

	go network.Listen(&kad)
	kademlia.CLIServer(&kad)
	// temp code to send pings to bootstrap
	for {
		c := kademlia.NewContact(kademlia.NewKademliaID(kademlia.BOOTSTRAP_ID), kademlia.BOOTSTRAP_IP)
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
