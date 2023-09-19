package kademlia

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

type Network struct {
	myContact *Contact
}

func CreateNetwork(myContact *Contact) Network {
	return Network{myContact}
}

func (network *Network) Listen(kademlia *Kademlia) {
	port, err := net.Listen(PROTOCOL, ":"+PORT)
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
		network.handleReq(rpc, kademlia, connection) // TODO async
	}
}

func (network *Network) SendPingMessage(kademlia *Kademlia, contact *Contact) RPC {
	fmt.Println("Pinging", contact)
	return network.sendReq(kademlia, contact.Address, RPC{PING, *network.myContact, *contact.ID, nil, nil})
}

func (network *Network) SendStoreReqMessage(kademlia *Kademlia, contact *Contact, hash KademliaID, data []byte) RPC {
	fmt.Println("Storing", data, "with hash", hash, "at", contact)
	return network.sendReq(kademlia, contact.Address, RPC{STORE_REQ, *network.myContact, hash, data, nil})
}

func (network *Network) SendFindContactReqMessage(kademlia *Kademlia, contact Contact, target *KademliaID) RPC {
	fmt.Println("Requesting Find_node", target, "at", contact.String())
	return network.sendReq(kademlia, contact.Address, RPC{FIND_NODE_REQ, *network.myContact, *target, nil, nil})
}

func (network *Network) SendFindDataReqMessage(kademlia *Kademlia, contact *Contact, hash string) RPC {
	fmt.Println("Finding data at", contact)
	return network.sendReq(kademlia, contact.Address, RPC{FIND_VALUE_REQ, *network.myContact, *NewKademliaIDString(hash), nil, nil})
}

// func (network *Network) SendFindDataResMessage(contact *Contact, hash string) { //TODO function not done
// 	fmt.Println("Returning data to", contact)
// 	network.sendRsp(contact.Address, RPC{FIND_VALUE_RSP, *network.myContact, *contact.ID, *NewKademliaIDString(hash), nil, nil})
// }

func (network *Network) sendRsp(address string, sendRpc RPC, connection net.Conn) {
	encoder := gob.NewEncoder(connection)
	// Encode and send the struct
	err := encoder.Encode(sendRpc)
	if err != nil {
		log.Println(err)
	}
}

func (network *Network) sendReq(kademlia *Kademlia, address string, sendRpc RPC) RPC {
	connection, err := net.Dial(PROTOCOL, fmt.Sprintf("%s:%s", address, PORT))
	if err != nil {
		log.Println(err)
		return UndefinedRPC()
	}
	defer connection.Close()

	encoder := gob.NewEncoder(connection)

	// Encode and send the struct
	err = encoder.Encode(sendRpc)
	if err != nil {
		log.Println(err)
		return UndefinedRPC()
	}
	// Set timeout for reply
	timeout := time.Second * 3 //TODO change to a constant
	connection.SetReadDeadline(time.Now().Add(timeout))

	decoder := gob.NewDecoder(connection)
	var responseRpc RPC
	// Wait and decode the received data into the struct
	err = decoder.Decode(&responseRpc)
	if err != nil {
		log.Println(err)
		return UndefinedRPC()
	}

	// After getting response -> add responder to routing table
	responseRpc.Sender.CalcDistance(network.myContact.ID)
	kademlia.table.AddContact(responseRpc.Sender)

	return responseRpc
}

func (network *Network) handleReq(rpc RPC, kademlia *Kademlia, connection net.Conn) {
	defer connection.Close()

	// Handle request types of RCPs
	switch rpc.Type {
	case PING:
		fmt.Println("Pinged by", rpc.Sender.String(), "now sending back pong")
		network.sendRsp(rpc.Sender.Address, RPC{PONG, *network.myContact, *rpc.Sender.ID, nil, nil}, connection)

	case STORE_REQ:
		fmt.Println("Store request from", rpc.Sender.String())
		var storeStatus string

		_, exist := kademlia.store[rpc.TargetID.String()]
		if exist {
			storeStatus = "has"
		} else {
			kademlia.store[rpc.TargetID.String()] = string(rpc.Data)
			storeStatus = "ok"
		}

		fmt.Println("Status:", storeStatus, kademlia.store)
		network.sendRsp(rpc.Sender.Address, RPC{STORE_RSP, *network.myContact, *rpc.Sender.ID, []byte(storeStatus), nil}, connection)

	case FIND_NODE_REQ:
		fmt.Println("Find node", rpc.TargetID.String(), "request from", rpc.Sender.String(), "now responding with the k-closest nodes")
		kClosestNodes := kademlia.table.FindClosestContacts(&rpc.TargetID, k)
		network.sendRsp(rpc.Sender.Address, RPC{FIND_NODE_RSP, *network.myContact, rpc.TargetID, nil, kClosestNodes}, connection)

	case FIND_VALUE_REQ:
		fmt.Println("Find value request from", rpc.Sender.String())

	case UNDEFINED:
		log.Println("ERROR: undefined RPC type")

	default:
		log.Println("ERROR: undefined RPC type or not a request type")

	}

	// Calculate distance from my ID to senders ID and update table
	rpc.Sender.CalcDistance(kademlia.network.myContact.ID)
	kademlia.table.AddContact(rpc.Sender)
}
