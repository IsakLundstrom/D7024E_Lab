package kademlia

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaIDString("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	furthest := NewContact(NewKademliaIDString("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	rt.AddContact(furthest)
	rt.AddContact(NewContact(NewKademliaIDString("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaIDString("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaIDString("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaIDString("1111111400000000000000000000000000000000"), "localhost:8002"))
	closest := NewContact(NewKademliaIDString("2111111400000000000000000000000000000000"), "localhost:8002")
	rt.AddContact(closest)

	contacts := rt.FindClosestContacts(NewKademliaIDString("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}

	closestBucketIndex := rt.getBucketIndex(closest.ID)
	furthestBucketIndex := rt.getBucketIndex(furthest.ID)

	if closestBucketIndex != 0 {
		t.Errorf("closest bucket index excpected [%d], got [%d]", 0, closestBucketIndex)
	}

	if furthestBucketIndex != len(rt.buckets)-1 {
		t.Errorf("closest bucket index excpected [%d], got [%d]", len(rt.buckets)-1, furthestBucketIndex)
	}

	if contacts[0].String() != closest.String() {
		t.Errorf("closest should be [%s] but was [%s]", closest.String(), contacts[0].String())
	}
	if contacts[len(contacts)-1].String() != furthest.String() {
		t.Errorf("furthest should be [%s] but was [%s]", furthest.String(), contacts[len(contacts)-1].String())
	}
}
