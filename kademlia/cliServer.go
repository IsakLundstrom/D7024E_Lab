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

	b := make([]byte, 512)
	connection.Read(b)
	args := strings.SplitN(string(b), " ", 2)
	
	fmt.Println("args[0]:", args[0])
	switch args[0] {
	case "ping":
		c := NewContact(NewKademliaID(BOOTSTRAP_ID), BOOTSTRAP_IP)
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
		contacts, data := kademlia.LookupData(args[1])
		fmt.Println("contact:", contacts, "data:", data, "data == ", data == "")
		res := "Could not find the data"
		if data != "" {
			res = fmt.Sprintf("The node: %s\nThe data: %s", contacts[0].String(), data)
		}
		connection.Write([]byte(res))
	case "exit":
		fmt.Println("EXITING")
		connection.Write([]byte("exiting"))
		connection.Close()
		os.Exit(0)
	default:
		fmt.Println("Invalid input!")
	}
}
