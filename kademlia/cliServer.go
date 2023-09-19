package kademlia

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func CLIServer(kademlia *Kademlia) {

	listener, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Println(err)
		return
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		CliHandler(connection, kademlia)
	}
}

func CliHandler(connection net.Conn, kademlia *Kademlia) {
	defer connection.Close()

	b := make([]byte, 128)
	connection.Read(b)
	args := strings.SplitN(string(b), " ", 2)

	switch args[0] {
	case "ping":
		c := NewContact(NewKademliaIDString(BOOTSTRAP_ID), BOOTSTRAP_IP)
		rpc := kademlia.network.SendPingMessage(kademlia, &c)
		if rpc.Type == PONG {
			connection.Write([]byte("PONG"))
		} else {
			connection.Write([]byte("Connection timedout..."))
		}
	case "put":
		fmt.Println("PUTTING")
		res, status := kademlia.Store([]byte(args[1]))
		if status == "FAIL" {
			connection.Write([]byte(status))
			return
		}

		connection.Write([]byte(res))
	case "get":
		fmt.Println("GETTING")
		res := kademlia.LookupData(args[1])
		//TODO check what type res will be
		connection.Write(res)
	case "exit":
		fmt.Println("EXITING")
		connection.Write([]byte("exiting"))
		connection.Close()
		os.Exit(0)
	default:
		fmt.Println("Invalid input!")
	}

}
