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
	bytesRead, err := connection.Read(buffer)
	if err != nil {
        log.Println("Error reading data:", err.Error())
        return
    }
	args := strings.SplitN(string(buffer[:bytesRead]), " ", 2)
	
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
		res, err := kademlia.Store([]byte(args[1]))
		if err != nil {
			connection.Write([]byte(err.Error()))
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
	case "forget":
		hash := args[1]
		kademlia.refreshMap.mutex.Lock()
		ch, exist := kademlia.refreshMap.rMap[hash]
		if exist {
			ch <- 0
			delete(kademlia.refreshMap.rMap, hash)
			connection.Write([]byte("Deleted the data with hash " + hash))
		} else {			
			connection.Write([]byte("Data with hash " + hash + " doesn't exist"))
		}
		kademlia.refreshMap.mutex.Unlock()
	default:
		fmt.Println("Invalid input!")
	}
}
