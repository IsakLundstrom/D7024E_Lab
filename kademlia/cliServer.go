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
		log.Println("Socket error:", err)
		return
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Connection error: ", err)
			continue
		}
		CliHandler(connection, kademlia)
	}
}

func CliHandler(connection net.Conn, kademlia *Kademlia) {
	defer connection.Close()

	buffer := make([]byte, 512)
	connection.Read(buffer)
	args := strings.SplitN(string(buffer), " ", 2)
	
	fmt.Println("Recived a", args[0], "command, now parsing it...")
	switch args[0] {
	case "ping":
		c := NewContact(NewKademliaID(BOOTSTRAP_ID), BOOTSTRAP_IP)
		rpc := kademlia.network.SendPingMessage(kademlia, &c)
		if rpc.Type == PONG {
			connection.Write([]byte("Pong from " + kademlia.table.me.String()))
		} else {
			connection.Write([]byte("Connection timedout..."))
		}
	case "put":
		res, status := kademlia.Store([]byte(args[1]))
		if status == "FAIL" {
			connection.Write([]byte(status))
			return
		}
		connection.Write([]byte(res))
	case "get":
		contacts, data := kademlia.LookupData(args[1])
		res := "Could not find the data"
		if data != "" {
			res = fmt.Sprintf("The node: %s\nThe data: %s", contacts[0].String(), data)
		}
		connection.Write([]byte(res))
	case "exit":
		connection.Write([]byte("exiting"))
		connection.Close()
		os.Exit(0)
	default:
		fmt.Println("Invalid input!")
	}
}
