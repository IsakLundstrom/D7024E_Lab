package kademlia

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func CLIServer() {

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
		CliHandler(connection)

	}
}

func CliHandler(connection net.Conn) {
	b := make([]byte, 128)
	connection.Read(b)
	args := strings.SplitN(string(b), " ", 2)
	fmt.Println("args", args)

	switch args[0] {
	case "ping":
		// TODO: Implement response
		connection.Write([]byte("pong"))
	}

	defer connection.Close()

}
