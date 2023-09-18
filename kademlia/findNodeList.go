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

	fmt.Println("  Recived contacts which is not queried, not in candidates and not me:")
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

	fmt.Println("  New candidates list:")
	for _, c := range findNodeList.candidates.contacts {
		fmt.Println("    ", c.String())
	}
}

// Returns true if responses are the k closest nodes known, else false
func (findNodeList *FindNodeList) checkKClosest() bool {
	done := false

	// Sort response contacts
	respondedClosest := ContactCandidates{findNodeList.responded}
	respondedClosest.Sort()

	// if any candidates exist
	candidateExists := findNodeList.candidates.Len() > 0 
	
	for i, c := range respondedClosest.contacts {
		if i >= k {
			done = true
			break
		}
		if !candidateExists || c.Less(&findNodeList.candidates.contacts[0]) {
			break
		}
	}

	return done
}

// // Returns true if kClosest have k nodes, else false
// func (findNodeList *FindNodeList) updateKClosest(roundResponseContacts *[]Contact) bool {
// 	fmt.Println("updateKClosest:")
// 	fmt.Println("  Round contacts which responded:")
// 	for _, c := range *roundResponseContacts {
// 		fmt.Println("    ", c.String())
// 	}

// 	done := false

// 	// Append and sort round response contacts
// 	findNodeList.kClosest.contacts = append(findNodeList.kClosest.contacts, *roundResponseContacts...)
// 	findNodeList.kClosest.Sort()

// 	newKClosest := []Contact{}
// 	candidateExists := findNodeList.candidates.Len() > 0 
	
// 	fmt.Println("  New kClosest list:")
// 	for i, c := range findNodeList.kClosest.contacts {
// 		if i >= k {
// 			done = true
// 			break
// 		}
// 		if !candidateExists || c.Less(&findNodeList.candidates.contacts[0]) {
// 			newKClosest = append(newKClosest, c)
// 			fmt.Println("    ", c.String())
// 		}
// 	}
	
// 	findNodeList.kClosest.contacts = newKClosest
// 	findNodeList.kClosest.Sort()

// 	return done
// }

// func (findNodeList *FindNodeList) finalUpdateKClosest() {
// 	respondedClosest := ContactCandidates{findNodeList.responded}
// 	respondedClosest.Sort()

// 	for _, c := range respondedClosest.contacts {
// 		if findNodeList.kClosest.Len() >= k {
// 			break
// 		}
// 		if !alreadyIn(c, findNodeList.kClosest.contacts) {
// 			findNodeList.kClosest.Append([]Contact{c});
// 		}
// 	}

// 	findNodeList.kClosest.Sort()
// }

func alreadyIn(contact Contact, contactList []Contact) bool {
	for _, c := range contactList{
		if c.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}