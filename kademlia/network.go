package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Network struct {
	myContact Contact
}

func (network *Network) Listen(kademlia *Kademlia) {
	port, err := net.Listen(PROTOCOL, ":" + PORT)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		connection, err := port.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		decoder := gob.NewDecoder(connection)
		var rpc RPC
	
		// Decode the received data into the struct
		err = decoder.Decode(&rpc)
		if err != nil {
			log.Println(err)
			continue
		}

		//TODO handle rpc in own thread, MUST add mutex then	
		network.handleRPC(rpc, kademlia)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("Pining", contact.Address)
	network.sendRCP(contact.Address, RPC{PING, network.myContact, *contact.ID, KademliaID{}, nil, nil})
}

func (network *Network) SendPongMessage(contact *Contact) {
	fmt.Println("Pinong", contact.Address)
	network.sendRCP(contact.Address, RPC{PONG, network.myContact, *contact.ID, KademliaID{}, nil, nil})
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	fmt.Println("Finding contact at", contact.Address)
	network.sendRCP(contact.Address, RPC{FIND_NODE_REQ, network.myContact, *contact.ID, KademliaID{}, nil, nil})
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) {
	fmt.Println("Finding data at", contact.Address)
	network.sendRCP(contact.Address, RPC{FIND_VALUE_REQ, network.myContact, *contact.ID, *NewKademliaID(hash), nil, nil})
}

func (network *Network) SendStoreMessage(contact *Contact, hash string, data []byte) {
	fmt.Println("Storing at", contact.Address)
	network.sendRCP(contact.Address, RPC{FIND_VALUE_REQ, network.myContact, *contact.ID, *NewKademliaID(hash), data, nil})
}

func  (network *Network) sendRCP(address string, rcp RPC) {
	connection, err := net.Dial(PROTOCOL, fmt.Sprintf("%s:%s", address, PORT))
	if err != nil {
		log.Println(err)
		return
	}
	encoder := gob.NewEncoder(connection)
	
	// Encode and send the struct
	err = encoder.Encode(rcp)
	if err != nil {
		log.Println(err)
	}

	connection.Close()
}

func (network *Network) handleRPC(rpc RPC, kademlia *Kademlia) {
	switch rpc.Type {
	case PING:
		fmt.Println("Pinged", rpc.Sender.Address)
		kademlia.table.AddContact(rpc.Sender)
		network.SendPongMessage(&rpc.Sender)
	case PONG:
		fmt.Println("Ponged", rpc.Sender.Address)
		kademlia.table.AddContact(rpc.Sender)
	case STORE_REQ:
		fmt.Println("Store request", rpc.Sender.Address)
	case STORE_RSP:
		fmt.Println("Store response", rpc.Sender.Address)
	case FIND_NODE_REQ:
		fmt.Println("Find node request", rpc.Sender.Address)
	case FIND_NODE_RSP:
		fmt.Println("Find node response", rpc.Sender.Address)
	case FIND_VALUE_REQ:
		fmt.Println("Find value request", rpc.Sender.Address)
	case FIND_VALUE_RSP:
		fmt.Println("Find value response", rpc.Sender.Address)
	case UNDEFINED:
		fallthrough
	default:
		log.Println("ERROR: undefined RPC type")
	}
}