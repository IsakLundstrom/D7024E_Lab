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

func (findNodeList *FindNodeList) updateCandidates(target *Contact, contacts *[]Contact) {
	fmt.Println("Already queried list:")
	for _, c := range findNodeList.queried {
		fmt.Println("  ", c.String())
	}

	fmt.Println("Recived contacts which is not queried:")
	// Find all not queried contacts from contacts
	notQueriedContacts := []Contact{}
	for _, c := range *contacts {
		if !findNodeList.alreadyQueried(c) {
			c.CalcDistance(target.ID)
			notQueriedContacts = append(notQueriedContacts, c)
			fmt.Println("  ", c.String())
		}
	}
	// Append and sort the contacts to the list of candidates
	findNodeList.candidates.Append(notQueriedContacts)
	findNodeList.candidates.Sort()

	fmt.Println("New candidates list:")
	for _, c := range findNodeList.candidates.contacts {
		fmt.Println("  ", c.String())
	}
}

func (findNodeList *FindNodeList) alreadyQueried(contact Contact) bool {
	for _, c := range findNodeList.queried{
		if c.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}