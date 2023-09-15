package kademlia

import (
	"fmt"
	"sync"
)

type FindNodeList struct {
	candidates 	ContactCandidates
	queried   	[]Contact
	responded 	[]Contact
	mutex		sync.Mutex
}

func NewFindNodeList() FindNodeList {
	return FindNodeList{
		candidates: ContactCandidates{[]Contact{}}, 
		queried: []Contact{}, 
		responded: []Contact{},
	}
}

func (findNodeList *FindNodeList) updateCandidates(me *Contact, target *Contact, contacts *[]Contact) {
	fmt.Println("Response contacts:")
	for _, c := range *contacts {
		fmt.Println("  ", c.String())
	}

	fmt.Println("Already queried list:")
	for _, c := range findNodeList.queried {
		fmt.Println("  ", c.String())
	}

	fmt.Println("Recived contacts which is not queried, not in candidates and not me:")
	// Find all not queried contacts from contacts
	notQueriedAndNotCandidates := []Contact{}
	for _, c := range *contacts {
		if !alreadyIn(c, findNodeList.queried) && !alreadyIn(c, findNodeList.candidates.contacts) && !c.ID.Equals(me.ID) {
			c.CalcDistance(target.ID)
			notQueriedAndNotCandidates = append(notQueriedAndNotCandidates, c)
			fmt.Println("  ", c.String())
		}
	}
	// Append and sort the contacts to the list of candidates
	findNodeList.candidates.Append(notQueriedAndNotCandidates)
	findNodeList.candidates.Sort()

	fmt.Println("New candidates list:")
	for _, c := range findNodeList.candidates.contacts {
		fmt.Println("  ", c.String())
	}
}

func alreadyIn(contact Contact, contactList []Contact) bool {
	for _, c := range contactList{
		if c.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}