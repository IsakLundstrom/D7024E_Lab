package kademlia

import (
	"fmt"
	"sync"
)

type FindNodeList struct {
	// kClosest 	ContactCandidates 	// nodes which have responded and have a dist < all nodes in candidates (aka dist < candidates[0])
	candidates 	ContactCandidates 	// current candidate nodes to request to (always sorted)
	queried   	[]Contact 			// all nodes which have been queried
	responded 	[]Contact 			// all nodes which have responded
	mutex		sync.Mutex
}

func NewFindNodeList() FindNodeList {
	return FindNodeList{
		// kClosest: ContactCandidates{[]Contact{}},
		candidates: ContactCandidates{[]Contact{}}, 
		queried: []Contact{}, 
		responded: []Contact{},
	}
}

func (findNodeList *FindNodeList) updateCandidates(me *Contact, target *Contact, contacts *[]Contact) {
	fmt.Println("updateCandidates:")
	fmt.Println("  Response contacts:")
	for _, c := range *contacts {
		fmt.Println("    ", c.String())
	}

	fmt.Println("  Already queried list:")
	for _, c := range findNodeList.queried {
		fmt.Println("    ", c.String())
	}
	
	fmt.Println("  Old candidates list:")
	for _, c := range findNodeList.candidates.contacts {
		fmt.Println("    ", c.String())
	}

	fmt.Println("  Recived contacts which is not queried, not in candidates and not me:")
	// Find all not queried contacts from contacts
	notQueriedAndNotCandidates := []Contact{}
	for _, c := range *contacts {
		if !alreadyIn(c, findNodeList.queried) && !alreadyIn(c, findNodeList.candidates.contacts) && !c.ID.Equals(me.ID) {
			c.CalcDistance(target.ID)
			notQueriedAndNotCandidates = append(notQueriedAndNotCandidates, c)
			fmt.Println("    ", c.String())
		}
	}
	// Append and sort the contacts to the list of candidates
	findNodeList.candidates.Append(notQueriedAndNotCandidates)
	findNodeList.candidates.Sort()

	fmt.Println("  New candidates list:")
	for _, c := range findNodeList.candidates.contacts {
		fmt.Println("    ", c.String())
	}
}

// Returns true if responses are the k closest nodes known, else false
func (findNodeList *FindNodeList) checkKClosest() bool {

	// Sort response contacts
	respondedClosest := ContactCandidates{findNodeList.responded}
	respondedClosest.Sort()

	candidateExists := findNodeList.candidates.Len() > 0 

	// if any candidate dosent exist and responded >= k, return true
	if !candidateExists && respondedClosest.Len() >= k {
		return true
	}

	// if we have >= k responses and all k responses have a dist < all in candidates (aka the kth responded dist < first candidate) 
	if respondedClosest.Len() >= k && respondedClosest.contacts[k - 1].Less(&findNodeList.candidates.contacts[0]) {
		return true
	}

	return false
}

func alreadyIn(contact Contact, contactList []Contact) bool {
	for _, c := range contactList{
		if c.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}