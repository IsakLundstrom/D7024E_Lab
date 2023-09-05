package main

import (
	"fmt"
	"kadlab/kademlia"
)

func main() {
	fmt.Println("hello world")
	if kademlia.IsBootstrap() {
		fmt.Println("IM BOOTERS")
	} else {
		fmt.Println("NOT BOOTSTRAp")
	}
	kademlia.Server()

}
