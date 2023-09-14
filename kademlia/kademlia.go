package kademlia

import (
	"fmt"
	"sync"
	"time"
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
	fmt.Println("Joining started...")
	// add bootstrap node to routing table
	kademlia.table.AddContact(NewContact(NewKademliaIDString(BOOTSTRAP_ID), BOOTSTRAP_IP))

	// lookup on itself
	kademlia.LookupContact(&kademlia.table.me)
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	// lookupMutex.Lock()
	// defer lookupMutex.Unlock()
	rpcChannel := make(chan RPC) //TODO maybe add channel limit

	// shortList := ShortList{[]ShortListItem{}}
	findNodeList := NewFindNodeList()

	// kClosestNodes := []Contact{} 
	// closestNode := kademlia.table.me
	closestNode := NewContact(kademlia.table.me.ID.InverseBitwise(), "")
	closestNode.CalcDistance(kademlia.network.myContact.ID)
	fmt.Println("closestnode distance:", closestNode.distance)

	// Get alpha closet nodes in my routing table
	alphaClosestNodes := kademlia.table.FindClosestContacts(target.ID, alpha)

	// Send FIND_NODE to #alpha nodes
	fmt.Println("Send FIND_NODE to #alpha nodes")
	for _, node := range alphaClosestNodes {
		// shortList.list = append(shortList.list, ShortListItem{node, true, false})
		findNodeList.mutex.Lock()
		findNodeList.queried = append(findNodeList.queried, node)
		findNodeList.mutex.Unlock()
		go func(node Contact) {
			rpcResponse := kademlia.network.SendFindContactReqMessage(kademlia, node, target.ID)
			if rpcResponse.Type == FIND_NODE_RSP {
				rpcChannel <- rpcResponse
			}
		}(node)
	}

	
	fmt.Println("Start iterative process...")
	for {
		roundTimeout := 3 * time.Second
		roundEndTime := time.Now().Add(roundTimeout)
		// roundClosestNode := closestNode 
		foundCloserNode := false

		round: for {
			fmt.Println("Round started")

			select {
			case <- time.After(time.Until(roundEndTime)):
				break round // i dont know yet
			case rpcResponse := <- rpcChannel:
				findNodeList.mutex.Lock()
				findNodeList.responded = append(findNodeList.responded, rpcResponse.Sender)
				// Check if >= k have responded already
				if len(findNodeList.responded) >= 20 {
					break round
				}
				findNodeList.updateCandidates(target, &rpcResponse.Nodes)
				// Update closestNode
				if rpcResponse.Sender.Less(&closestNode) {
					closestNode = rpcResponse.Sender
					foundCloserNode = true
				}
				findNodeList.mutex.Unlock()
			}
			
		}

		if foundCloserNode {
			fmt.Println("New Round.", "closestnode:", closestNode.ID, " distance:", closestNode.distance)

			// New round
			findNodeList.mutex.Lock()
			min := findNodeList.candidates.Len()
			if alpha < min {
				min = alpha 
			} 
			nodes := findNodeList.candidates.GetContacts(min)
			for _, node := range nodes {
				findNodeList.queried = append(findNodeList.queried, node)
				go func(node Contact) {
					rpcResponse := kademlia.network.SendFindContactReqMessage(kademlia, node, target.ID)
					if rpcResponse.Type == FIND_NODE_RSP {
						rpcChannel <- rpcResponse
					}
				}(node)
			}

			findNodeList.candidates.contacts = findNodeList.candidates.contacts[min:] // remove called nodes
			findNodeList.mutex.Unlock()
			continue
		} else {
			fmt.Println("No closer node found in round")

			// No closer node found -> send find node to rest k nodes which have not already been queried
			findNodeList.mutex.Lock()
			rest := k - len(findNodeList.queried)
			min := findNodeList.candidates.Len()
			if rest < min {
				min = rest 
			} 
			nodes := findNodeList.candidates.GetContacts(min)
			for _, node := range nodes {
				findNodeList.queried = append(findNodeList.queried, node) // this call shouldnt matter
				go func(node Contact) {
					rpcResponse := kademlia.network.SendFindContactReqMessage(kademlia, node, target.ID)
					if rpcResponse.Type == FIND_NODE_RSP {
						rpcChannel <- rpcResponse
					}
				}(node)
			}
			findNodeList.mutex.Unlock()
			break
		}
	}
	fmt.Println("Lookup done?")

	return findNodeList.responded
		
		//todo send next round

	//  <- rpcChannel



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
	// return nil
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	// TODO
	return []byte("TODO")
}

func (kademlia *Kademlia) Store(data []byte) []byte {
	// TODO
	return []byte("TODO")
}
