package main

type RPC struct {
	Type 		RPCType
	Sender 		Contact
	TargetID 	KademliaID
	Hash		KademliaID
	Data		[]byte
	Nodes		[]Contact
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