package kademlia

import (
	"fmt"
	"log"
	"time"
)

var k int = 20
var alpha int = 3

type Kademlia struct {
	table   *RoutingTable
	network *Network
	store   DataStorage
}

func CreateKademlia(network *Network) Kademlia {
	return Kademlia{
		table:   NewRoutingTable(*network.myContact),
		network: network,
		store:   NewDataStorage()}
}

func (kademlia *Kademlia) JoinNetwork() {
	fmt.Println("Joining started...")
	// add bootstrap node to routing table
	kademlia.table.AddContact(NewContact(NewKademliaID(BOOTSTRAP_ID), BOOTSTRAP_IP))

	// lookup on itself
	res := kademlia.LookupContact(kademlia.table.me.ID)

	fmt.Println("Join lookup result:")
	for _, c := range res {
		fmt.Println("  ", c.String())
	}
}

func (kademlia *Kademlia) LookupContact(targetID *KademliaID) []Contact {
	contacts, _ := kademlia.iterativeFind(targetID, FIND_NODE_REQ)
	return contacts
}

func (kademlia *Kademlia) LookupData(hash string) ([]Contact, string) {
	return kademlia.iterativeFind(NewKademliaID(hash), FIND_VALUE_REQ)
}

func (kademlia *Kademlia) Store(data []byte) (string, string) {
	rpcChannel := make(chan RPC)

	hash := NewKademliaID(GetHash(data))
	contacts := kademlia.LookupContact(hash)

	kademlia.table.me.CalcDistance(hash)
	meCloser := len(contacts) == 0 || kademlia.table.me.Less(&contacts[len(contacts)-1])

	if meCloser {
		if len(contacts) == k {
			contacts = contacts[:len(contacts)-1]
		}
		contacts = append(contacts, kademlia.table.me)
	}

	okCounter := 0
	hasCounter := 0
	failCounter := 0

	for _, c := range contacts {

		go func(c Contact) {
			rpcResponse := kademlia.network.SendStoreReqMessage(kademlia, &c, *hash, data)
			if rpcResponse.Type == STORE_RSP {
				rpcChannel <- rpcResponse
			} else {
				log.Println("Recived response type:", rpcResponse.Type, "but expected:", STORE_RSP)
			}
		}(c)
	}

	timeOut := 300 * time.Millisecond
	endTime := time.Now().Add(timeOut)

listen:
	for {
		select {
		case <-time.After(time.Until(endTime)):
			fmt.Println("Store wait for responses timedout!")
			break listen
		case rpcResponse := <-rpcChannel:
			switch string(rpcResponse.Data) {
			case "ok":
				okCounter++
			case "has":
				hasCounter++
			default:
				failCounter++
			}
		}
	}

	fmt.Println("sentTo", len(contacts), "okCounter:", okCounter, "hasCounter:", hasCounter, "failCounter:", failCounter)

	if len(contacts) == okCounter+hasCounter {
		return hash.String(), "OK"
	}
	return hash.String(), "FAIL"
}

func (kademlia *Kademlia) iterativeFind(targetID *KademliaID, findType RPCType) ([]Contact, string) {
	// Response channel and response storage
	rpcChannel := make(chan RPC) //TODO maybe add channel limit
	findNodeList := NewFindNodeList()

	// Init closest node as furthest away as possible
	closestNode := NewContact(kademlia.table.me.ID.InverseBitwise(), "")
	closestNode.CalcDistance(kademlia.network.myContact.ID)

	// Get alpha closet nodes in my routing table and set these as first candidates
	alphaClosestNodes := kademlia.table.FindClosestContacts(targetID, alpha)
	findNodeList.candidates.contacts = alphaClosestNodes

	// Round variables
	roundNr := 1 // only used for prints
	roundTimeout := 300 * time.Millisecond
	foundCloserNode := true
	finalRound := false

	fmt.Println("Start rounds / iterative process...")
iterativeProcess:
	for {
		// New round
		fmt.Println("New round started:", roundNr)
		fmt.Println("closestnode:", closestNode.ID.String(), " distance:", closestNode.distance)

		roundEndTime := time.Now().Add(roundTimeout)

		// Send requests
		kademlia.newRequestRound(&findNodeList, targetID, findType, &rpcChannel, foundCloserNode)
		if !foundCloserNode {
			fmt.Println("No closer node found in previous round, now just wait for last responses.")
			finalRound = true
		}

		foundCloserNode = false

		// Listen for round responses
	listen:
		for {
			select {
			case <-time.After(time.Until(roundEndTime)):
				if finalRound {
					fmt.Println("Final round over")
					findNodeList.mutex.Lock()
					break iterativeProcess
				}
				break listen
			case rpcResponse := <-rpcChannel:
				switch rpcResponse.Type {
				case FIND_VALUE_RSP:
					if rpcResponse.Data != nil {
						if len(findNodeList.responded) > 0 {
							c := ContactCandidates{findNodeList.responded}
							c.Sort()
							kademlia.network.SendStoreReqMessage(kademlia, &c.contacts[0], rpcResponse.TargetID, rpcResponse.Data)
						}
						return []Contact{rpcResponse.Sender}, string(rpcResponse.Data)
					}
					fallthrough
				case FIND_NODE_RSP:
					findNodeList.mutex.Lock()
					findNodeList.responded = append(findNodeList.responded, rpcResponse.Sender)

					fmt.Println("Find node response from", rpcResponse.Sender.String())
					findNodeList.updateCandidates(&kademlia.table.me, targetID, &rpcResponse.Nodes)
					// Update closestNode
					if rpcResponse.Sender.Less(&closestNode) {
						closestNode = rpcResponse.Sender
						foundCloserNode = true
					}
					findNodeList.mutex.Unlock()
				default:
					log.Println("Recived response type:", rpcResponse.Type, "but expected:", FIND_NODE_RSP, "or", FIND_VALUE_RSP)
				}
			}
		}
		fmt.Println("Round over")

		findNodeList.mutex.Lock()
		done := findNodeList.checkKClosest(k)

		if done {
			break iterativeProcess
		}
		findNodeList.mutex.Unlock()
		roundNr++
	}

	fmt.Println("Lookup done")
	fmt.Println(kademlia.table.String())

	// Final part to extract the closest nodes which repsonded
	kClosest := ContactCandidates{findNodeList.responded}
	kClosest.Sort()
	numNodes := kClosest.Len()
	if numNodes > k {
		numNodes = k
	}
	defer findNodeList.mutex.Unlock()

	return kClosest.GetContacts(numNodes), ""
}

func (kademlia *Kademlia) newRequestRound(findNodeList *FindNodeList, targetID *KademliaID, findType RPCType, rpcChannel *chan RPC, foundCloserNode bool) {
	findNodeList.mutex.Lock()
	// Selects another alpha nodes until a round doesn't find a closer node -> requests to each of the k closest nodes that it has not already queried.
	numRequests := alpha
	if !foundCloserNode {
		if findNodeList.candidates.Len() > k {
			numRequests = k
		} else {
			numRequests = findNodeList.candidates.Len()
		}
	}
	// Bound numRequests to numCandidates
	numCandidates := findNodeList.candidates.Len()
	if numRequests > numCandidates {
		numRequests = numCandidates
	}
	fmt.Println("Requesting to", numRequests, "more nodes")
	nodes := findNodeList.candidates.GetContacts(numRequests)
	for _, node := range nodes {
		findNodeList.queried = append(findNodeList.queried, node)
		go func(node Contact) {
			rpcResponse := kademlia.network.SendFindReqMessage(kademlia, node, targetID, findType)
			*rpcChannel <- rpcResponse
		}(node)
	}
	findNodeList.candidates.contacts = findNodeList.candidates.contacts[numCandidates:] // remove called nodes
	findNodeList.mutex.Unlock()
}
