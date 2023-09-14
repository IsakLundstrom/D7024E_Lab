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
	b := make([]byte, 128)
	connection.Read(b)
	args := strings.SplitN(string(b), " ", 2)

	switch args[0] {
	case "ping":
		// TODO: Implement response
		connection.Write([]byte("pong"))
	case "put":
		fmt.Println("PUTTING")
		res := kademlia.Store([]byte(args[1]))
		//TODO check if res is bytes
		connection.Write(res)
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

	defer connection.Close()

}
