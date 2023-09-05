package kademlia

import (
	"fmt"
	"net"
)

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {

}
func SendMsg(address string, msg string) {

	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	connection.Write([]byte(msg))
	connection.Close()

}

func Server() {
	port, err := net.Listen("tcp", ":4000")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		connection, err := port.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(connection.RemoteAddr(), connection)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
