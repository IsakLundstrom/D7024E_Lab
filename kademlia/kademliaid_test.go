package kademlia

import (
	"testing"
)

func TestEquals(t *testing.T) {
	id1 := NewKademliaID("0000000000000000000000000000000000000000")
	id2 := NewKademliaID("0000000000000000000000000000000000000000")
	id3 := NewKademliaID("0000000000000000000000000000000000000001")

	if !id1.Equals(id2) {
		t.Errorf("equals excpected [%s], got [%s]", "00000000000000000000", id2)
	}

	if id1.Equals(id3) {
		t.Errorf("equals excpected [%s], got [%s]", "0000000000000000000000000000000000000000", id2)
	}

}

func TestRandomId(t *testing.T) {
	for i := 0; i < 10; i++ {
		id1 := NewRandomKademliaID()
		id2 := NewRandomKademliaID()

		if id1.Equals(id2) {
			t.Errorf("equals excpected different ones, got the same ones")
		}
	}
}

func TestLess(t *testing.T) {
	id1 := NewKademliaID("0000000000000000000000000000000000000000")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")

	if !id1.Less(id2) {
		t.Errorf("Less excpected id1 < id2, got id1 > id2")
	}
	if id2.Less(id1) {
		t.Errorf("Less excpected id1 < id2, got id1 > id2")
	}
}

func TestInverse(t *testing.T) {
	id1 := NewKademliaID("0000000000000000000000000000000000000000")
	id2 := NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")

	inverse := id1.InverseBitwise()

	if !inverse.Equals(id2) {
		t.Errorf("equals excpected different ones, got the same ones")
	}

}
