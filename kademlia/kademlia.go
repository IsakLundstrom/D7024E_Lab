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

func CreateKademlia(network *Network) Kademlia {
	return Kademlia{NewRoutingTable(*network.myContact), network, map[KademliaID][]byte{}}

}

func (kademlia *Kademlia) JoinNetwork() {
	fmt.Println("Joining started...")
	// add bootstrap node to routing table
	kademlia.table.AddContact(NewContact(NewKademliaIDString(BOOTSTRAP_ID), BOOTSTRAP_IP))

	// lookup on itself
	res := kademlia.LookupContact(&kademlia.table.me)

	fmt.Println("Join lookup result:", res)
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	// Response channel and response storage
	rpcChannel := make(chan RPC) //TODO maybe add channel limit
	findNodeList := NewFindNodeList()

	// Init closest node as furthest away as possible
	closestNode := NewContact(kademlia.table.me.ID.InverseBitwise(), "")
	closestNode.CalcDistance(kademlia.network.myContact.ID)

	// Get alpha closet nodes in my routing table and set these as first candidates
	alphaClosestNodes := kademlia.table.FindClosestContacts(target.ID, alpha)
	findNodeList.candidates.contacts = alphaClosestNodes

	// Round variables
	roundNr := 1
	roundTimeout := 300 * time.Millisecond
	foundCloserNode := true
	finalRound := false

	fmt.Println("Start rounds / iterative process...")
	mainLoop: for {
		// New round
		fmt.Println("New round started:", roundNr)
		fmt.Println("closestnode:", closestNode.ID.String(), " distance:", closestNode.distance)
		
		roundEndTime := time.Now().Add(roundTimeout)

		// Send requests
		kademlia.newRequestRound(&findNodeList, target, &rpcChannel, foundCloserNode)
		if !foundCloserNode {
			fmt.Println("No closer node found in previous round, now just wait for last responses.")
			finalRound = true
		}

		foundCloserNode = false

		// Listen for round responses
		round: for {
			select {
			case <-time.After(time.Until(roundEndTime)):
				if finalRound {
					fmt.Println("Final round over")
					break mainLoop
				}
				fmt.Println("Round over")
				break round
			case rpcResponse := <-rpcChannel:
				findNodeList.mutex.Lock()
				findNodeList.responded = append(findNodeList.responded, rpcResponse.Sender)
				// Check if >= k have responded already
				if len(findNodeList.responded) >= 20 {
					fmt.Println("k contacts have responded already -> done")
					findNodeList.mutex.Unlock()
					break mainLoop
				}
				fmt.Println("Find node response from", rpcResponse.Sender.String())
				findNodeList.updateCandidates(&kademlia.table.me, target, &rpcResponse.Nodes)
				// Update closestNode
				if rpcResponse.Sender.Less(&closestNode) {
					closestNode = rpcResponse.Sender
					foundCloserNode = true
				}
				findNodeList.mutex.Unlock()
			}

		}
		roundNr++
	}
	fmt.Println("Lookup done?")
	fmt.Println(kademlia.table.String())

	return findNodeList.responded
}

func (kademlia *Kademlia) newRequestRound(findNodeList *FindNodeList, target *Contact, rpcChannel *chan RPC, round bool) {
	findNodeList.mutex.Lock()
	var maxRequestNodes int
	if round {
		maxRequestNodes = alpha
	} else {
		maxRequestNodes = k - len(findNodeList.queried)
	}
	numCandidates := findNodeList.candidates.Len()
	if maxRequestNodes > numCandidates {
		maxRequestNodes = numCandidates
	}
	fmt.Println("Requesting to", maxRequestNodes, "more nodes")
	nodes := findNodeList.candidates.GetContacts(maxRequestNodes)
	for _, node := range nodes {
		findNodeList.queried = append(findNodeList.queried, node)
		go func(node Contact) {
			rpcResponse := kademlia.network.SendFindContactReqMessage(kademlia, node, target.ID)
			if rpcResponse.Type == FIND_NODE_RSP {
				*rpcChannel <- rpcResponse
			}
		}(node)
	}
	findNodeList.candidates.contacts = findNodeList.candidates.contacts[numCandidates:] // remove called nodes
	findNodeList.mutex.Unlock()
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	// TODO
	return []byte("TODO")
}

func (kademlia *Kademlia) Store(data []byte) []byte {
	// TODO
	return []byte("TODO")
}
