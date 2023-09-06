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

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	sendRCP(contact.Address, RPC{PING, network.myContact, *contact.ID, KademliaID{}, nil, nil})
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	sendRCP(contact.Address, RPC{FIND_NODE_REQ, network.myContact, *contact.ID, KademliaID{}, nil, nil})
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) {
	sendRCP(contact.Address, RPC{FIND_VALUE_REQ, network.myContact, *contact.ID, *NewKademliaID(hash), nil, nil})
}

func (network *Network) SendStoreMessage(contact *Contact, hash string, data []byte) {
	sendRCP(contact.Address, RPC{FIND_VALUE_REQ, network.myContact, *contact.ID, *NewKademliaID(hash), data, nil})
}

func sendRCP(address string, msg RPC) {
	connection, err := net.Dial("tcp", address + ":4000")
	if err != nil {
		log.Println(err)
		return
	}
	encoder := gob.NewEncoder(connection)
	
	// Encode and send the struct
	err = encoder.Encode(msg)
    if err != nil {
        log.Println(err)
    }

	connection.Close()
}

func Server() {
	port, err := net.Listen("tcp", ":4000")
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
	
		// You can now use the receivedData as needed
		fmt.Println("Received: ", rpc)
	}
}
