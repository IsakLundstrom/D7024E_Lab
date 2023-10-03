package kademlia

import "sync"

type RefreshMap struct {
	rMap 	map[string](chan int)
	mutex	sync.Mutex
}