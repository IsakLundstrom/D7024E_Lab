package kademlia

import (
	"fmt"
	"testing"
)

func TestFindNodeList(t *testing.T) {

	nodeList := NewFindNodeList()
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	target := NewKademliaID("2222222200000000000000000000000000000000")
	c2 := NewContact(NewKademliaID("1234567800000000000000000000000000000000"), "localhost:8002")
	c3 := NewContact(NewKademliaID("EEEEEEEE00000000000000000000000000000000"), "localhost:8002")
	c4 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "localhost:8003")
	c5 := NewContact(NewKademliaID("9999999900000000000000000000000000000000"), "localhost:8005")

	r1 := NewContact(NewKademliaID("229999AA00000000000000000000000000000000"), "localhost:8005")
	r2 := NewContact(NewKademliaID("229999BB00000000000000000000000000000000"), "localhost:8005")
	r3 := NewContact(NewKademliaID("229999CC00000000000000000000000000000000"), "localhost:8005")
	r1.CalcDistance(target)
	r2.CalcDistance(target)
	r3.CalcDistance(target)

	nodeList.responded = []Contact{r1, r2, r3}

	if !nodeList.checkKClosest(3) {
		t.Errorf("checkKClosest failed!, should have 3 responded and 0 candidates, but had [%d], and [%d]", len(nodeList.responded), len(nodeList.queried))
	}

	nodeList.updateCandidates(&me, target, &[]Contact{c2, c3, c4})

	if len(nodeList.candidates.contacts) != 3 {
		t.Errorf("updateCandidates failed! added 3 contacts but length was [%d]", len(nodeList.candidates.contacts))
	}

	nodeList.updateCandidates(&me, target, &[]Contact{me, c2})
	if len(nodeList.candidates.contacts) != 3 {
		t.Errorf("updateCandidates failed! was able to add \"me\" or a duplicate contact")
	}

	nodeList.updateCandidates(&me, target, &[]Contact{c5})
	if len(nodeList.candidates.contacts) != 4 {
		t.Errorf("updateCandidates failed! could not add new contact")
	}

	nodeList.responded = []Contact{r1, r2, r3}
	nodeList.queried = []Contact{c5}
	if !nodeList.checkKClosest(3) {
		t.Errorf("checkKClosest failed!")
		fmt.Println("contacts:")
		for _, c := range nodeList.candidates.contacts {
			fmt.Println("    ", c.String())
		}
		fmt.Println("responded:")
		for _, c := range nodeList.responded {
			fmt.Println("    ", c.String())
		}
	} //TODO

	fmt.Println(len(nodeList.candidates.contacts))
}
