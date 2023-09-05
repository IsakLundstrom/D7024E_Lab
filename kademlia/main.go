package main

import (
	"fmt"
	"log"
)

func main() {
	var myRoutingTable *RoutingTable
	var myId *KademliaID
	var myContact Contact

	fmt.Println("hello world")
	if IsBootstrap() {
		fmt.Println("IM BOOTERS")
		myId = NewKademliaID(BOOTSTRAP_ID)

	} else {
		fmt.Println("NOT BOOTSTRAp")
		myId = NewRandomKademliaID()
	}

	myIp, err := GetMyIp()
	if err != nil {
		log.Fatal(err)
	}
	myContact = NewContact(myId, myIp)
	myRoutingTable = NewRoutingTable(myContact)
	fmt.Println(myRoutingTable)
	Server()

}
