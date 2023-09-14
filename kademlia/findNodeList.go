package kademlia

import "sync"

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
	// Find all not queried contacts from contacts
	notQueriedContacts := []Contact{}
	for _, c := range *contacts {
		if findNodeList.notQueried(c) {
			c.CalcDistance(target.ID)
			notQueriedContacts = append(notQueriedContacts, c)
		}
	}
	// Append and sort the contacts to the list of candidates
	findNodeList.candidates.Append(notQueriedContacts)
	findNodeList.candidates.Sort()
}

func (findNodeList *FindNodeList) notQueried(contact Contact) bool {
	for _, c := range findNodeList.queried{
		if c.ID == contact.ID {
			return false
		}
	}
	return true
}