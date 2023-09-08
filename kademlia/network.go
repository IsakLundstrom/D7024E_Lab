package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Network struct {
	myContact 	*Contact
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
		network.handleRPC(rpc, kademlia) // TODO async
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("Pinging", contact)
	network.sendRCP(contact.Address, RPC{PING, *network.myContact, *contact.ID, KademliaID{}, nil, nil})
}

func (network *Network) SendPongMessage(contact *Contact) {
	fmt.Println("Ponging", contact)
	network.sendRCP(contact.Address, RPC{PONG, *network.myContact, *contact.ID, KademliaID{}, nil, nil})
}

func (network *Network) SendStoreReqMessage(contact *Contact, hash KademliaID, data []byte) {
	fmt.Println("Storing", data, "with hash", hash, "at", contact)
	network.sendRCP(contact.Address, RPC{STORE_REQ, *network.myContact, *contact.ID, hash, data, nil})
}

func (network *Network) SendStoreRspMessage(contact *Contact) {
	fmt.Println("Stored response from", contact)
	network.sendRCP(contact.Address, RPC{STORE_RSP, *network.myContact, *contact.ID, KademliaID{}, nil, nil})
}

func (network *Network) SendFindContactReqMessage(contact *Contact, target *KademliaID) {
	fmt.Println("Finding node", target, " at", contact)
	network.sendRCP(contact.Address, RPC{FIND_NODE_REQ, *network.myContact, *target, KademliaID{}, nil, nil})
}

func (network *Network) SendFindContactRspMessage(contact *Contact, target *KademliaID, nodes []Contact) {
	fmt.Println("Returing nodes to", contact)
	network.sendRCP(contact.Address, RPC{FIND_NODE_RSP, *network.myContact, *contact.ID, *target, nil, nil})
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) {
	fmt.Println("Finding data at", contact)
	network.sendRCP(contact.Address, RPC{FIND_VALUE_REQ, *network.myContact, *contact.ID, *NewKademliaIDString(hash), nil, nil})
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
	// Calculate distance from my ID to senders ID and update table
	rpc.Sender.CalcDistance(kademlia.network.myContact.ID)
	kademlia.table.AddContact(rpc.Sender)

	// Handle all types of RCPs
	switch rpc.Type {
	case PING:
		fmt.Println("Pinged", rpc.Sender)

		network.SendPongMessage(&rpc.Sender)

	case PONG:
		fmt.Println("Ponged", rpc.Sender)

	case STORE_REQ:
		fmt.Println("Store request", rpc.Sender)

		kademlia.store[rpc.Hash] = rpc.Data
		fmt.Println(kademlia.store) //TODO remove
		network.SendStoreRspMessage(&rpc.Sender)

	case STORE_RSP:
		fmt.Println("Store response", rpc.Sender)

	case FIND_NODE_REQ:
		fmt.Println("Find node request", rpc.Sender)

		kClosestNodes := kademlia.table.FindClosestContacts(&rpc.TargetID, k)
		network.SendFindContactRspMessage(&rpc.Sender, &rpc.TargetID, kClosestNodes)

	case FIND_NODE_RSP:
		fmt.Println("Find node response", rpc.Sender)
		// LookupChannel <- struct {rpc.Sender; rpc.Nodes} LYCKA TILL :)

	case FIND_VALUE_REQ:
		fmt.Println("Find value request", rpc.Sender)

	case FIND_VALUE_RSP:
		fmt.Println("Find value response", rpc.Sender)

	case UNDEFINED:
		fallthrough
	default:
		log.Println("ERROR: undefined RPC type")
	}
}

func contactInArray(contact Contact, list []Contact) bool {
    for _, c := range list {
        if c == contact {
            return true
        }
    }
    return false
}