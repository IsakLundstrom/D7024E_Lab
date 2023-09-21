package kademlia

import (
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	furthest := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	rt.AddContact(furthest)
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	closest := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002")
	rt.AddContact(closest)

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)

	closestBucketIndex := rt.getBucketIndex(closest.ID)
	furthestBucketIndex := rt.getBucketIndex(furthest.ID)

	if closestBucketIndex != 0 {
		t.Errorf("closest bucket index excpected [%d], got [%d]", 0, closestBucketIndex)
	}

	if furthestBucketIndex != len(rt.buckets)-1 {
		t.Errorf("closest bucket index excpected [%d], got [%d]", len(rt.buckets)-1, furthestBucketIndex)
	}

	if contacts[0].ID != closest.ID || contacts[0].Address != closest.Address {
		t.Errorf("closest should be [%s] but was [%s]", closest.String(), contacts[0].String())
	}
	if contacts[len(contacts)-1].ID != furthest.ID || contacts[len(contacts)-1].Address != furthest.Address {
		t.Errorf("furthest should be [%s] but was [%s]", furthest.String(), contacts[len(contacts)-1].String())
	}
}

func TestRoutingTableString(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))

	expectedString := "RoutingTable:\n"
	expectedString += "  Me: {ID:ffffffff00000000000000000000000000000000 Address:localhost:8000 distance:<nil>}\n"
	expectedString += "  Buckets:\n"
	expectedString += "    Bucket 0:\n"
	expectedString += "      List len: 1\n"

	if rt.String() != expectedString {
		t.Errorf("string output does not match expected")
	}
}