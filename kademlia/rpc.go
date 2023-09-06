package main

import (
	"fmt"
	"log"
)

type RPC struct {
	Type 	RPCType
	Sender 	Contact
	Target 	KademliaID
	Hash	KademliaID
	Data	[]byte
	Nodes	[]Contact
}

type RPCType int64

const (
	UNDEFINED RPCType = iota
	PING
	PONG
	STORE_REQ
	STORE_RSP
	FIND_NODE_REQ
	FIND_NODE_RSP
	FIND_VALUE_REQ
	FIND_VALUE_RSP
)

func HandleRPC(rpc RPC) {
	switch rpc.Type {
	case PING:
		fmt.Println("Pinged")
	case PONG:
		fmt.Println("Ponged")
	case STORE_REQ:
		fmt.Println("Store request")
	case STORE_RSP:
		fmt.Println("Store response")
	case FIND_NODE_REQ:
		fmt.Println("Find node request")
	case FIND_NODE_RSP:
		fmt.Println("Find node response")
	case FIND_VALUE_REQ:
		fmt.Println("Find value request")
	case FIND_VALUE_RSP:
		fmt.Println("Find value response")
	case UNDEFINED:
		fallthrough
	default:
		log.Println("ERROR: undefined RPC type")
	}
}