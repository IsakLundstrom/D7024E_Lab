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
func SendMsg() {

	connection, err := net.Dial("udp", "172.19.0.2:4000")
	if err != nil {
		fmt.Println(err)
		return
	}
	connection.Write([]byte("Hello"))
	connection.Close()

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
