package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	args := os.Args

	fmt.Println("KWODIKAD")
	fmt.Println(args, len(args))

	if len(args) < 2 {
		fmt.Println("No arguments")
		return
	}
	if len(args) > 3 {
		fmt.Println("To many arguments")
		return
	}
	switch args[1] {
	case "ping":
		fmt.Println("PINGING")
		CliSend("ping " + args[2])

	}

}

func CliSend(data string) {
	connection, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		log.Println(err)
		return
	}
	connection.Write([]byte(data))
	response := make([]byte, 128)
	connection.Read(response)
	fmt.Println("Node response", string(response))
	defer connection.Close()
}
