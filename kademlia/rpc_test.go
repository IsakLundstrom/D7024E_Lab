package kademlia

import (
	"testing"
)

func TestUndefinedRPC(t *testing.T) {
	rpc := RPC{UNDEFINED, Contact{}, KademliaID{}, nil, nil}
	if rpc.Type != UndefinedRPC().Type {
		t.Error("did not get RPCType.UNDEFINED")
	}
}
