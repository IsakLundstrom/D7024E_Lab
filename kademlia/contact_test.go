package kademlia

import (
	"testing"
)

func TestContact(t *testing.T) {
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	contact2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "localhost:8000")

	contact1.CalcDistance(contact1.ID)
	contact2.CalcDistance(contact1.ID)
	if !contact1.distance.Equals(NewKademliaID("0000000000000000000000000000000000000000")) {
		t.Errorf("contact distance to itself was [%s] and not [%s]", contact1.distance, "0000000000000000000000000000000000000000")
	}
	if !contact2.distance.Equals(NewKademliaID("5555555500000000000000000000000000000000")) {
		t.Errorf("contact.CalcDistance failed!")
	}
	if contact1.Less(&contact1) {
		t.Errorf("contact.Less failed! contact distance was less than itself")
	}

	if contact2.Less(&contact1) {
		t.Errorf("distance [%s] was less than [%s]", contact2.distance, contact1.distance)
	}
}
