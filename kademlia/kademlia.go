package kademlia

import (
	"fmt"
	"sync"
)

var k int = 20
var alpha int = 3
var lookupMutex sync.Mutex

// var LookupChannel = make(chan struct {concact Contact, list []Contact}, alpha) utkommenterat

type Kademlia struct {
	table   *RoutingTable
	network *Network
	store   map[KademliaID][]byte
}

func CreateKademlia(myContact *Contact, network *Network) Kademlia {
	return Kademlia{NewRoutingTable(*myContact), network, map[KademliaID][]byte{}}

}

func (kademlia *Kademlia) JoinNetwork() {
	fmt.Println("TODO Joining...")
	// add bootstrap node to routing table
	kademlia.table.AddContact(NewContact(NewKademliaIDString(BOOTSTRAP_ID), BOOTSTRAP_IP))

	// lookup on itself
	kademlia.LookupContact(&kademlia.table.me)
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	lookupMutex.Lock()
	defer lookupMutex.Unlock()

	shortList := ShortList{[]ShortListItem{}}
	// closestNode := Contact{} utkommenterat

	// Get k closet nodes in my routing table
	alphaClosestNodes := kademlia.table.FindClosestContacts(target.ID, alpha)

	// Send FIND_NODE to #alpha nodes
	for _, node := range alphaClosestNodes {
		// *kademlia.network.queried = append(*kademlia.network.queried, node)
		shortList.list = append(shortList.list, ShortListItem{node, true, false})
		go kademlia.network.SendFindContactReqMessage(&node, target.ID)
	}

	// contact, kContacts := <- LookupChannel

	// shortList.fill(kContacts)

	// // Update closestNode
	// if network.closestNode.ID !=nil || rpc.Sender.Less(network.closestNode) {
	// 	*network.closestNode = rpc.Sender
	// }

	// *network.queried.Fill()

	// //TODO if sender is node looked for set done! if sender == target
	// if *rpc.Sender.ID == rpc.TargetID {
	// 	fmt.Println("Found node at ", rpc.Sender)
	// 	// kademlia.queried = []Contact{}
	// }

	// count := 0
	// for _, node := range rpc.Nodes {
	// 	if count >= alpha {
	// 		break
	// 	}
	// 	if contactInArray(node, *network.queried) {
	// 		continue
	// 	}
	// 	*network.queried = append(*network.queried, node)
	// 	go network.SendFindContactReqMessage(&node, &rpc.TargetID)
	// 	count++
	// }

	//TODO we expect a answer, maybe handle it?
	return nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
