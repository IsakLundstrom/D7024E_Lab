package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	args := os.Args

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

		fmt.Println("pinging bootstrap...")
		CliSend("ping ")

	case "put":
		if len(args) != 3 {
			fmt.Println("nothing to put")
			return
		}
		fmt.Println("put call...")
		CliSend("put " + args[2])

	case "get":
		if len(args) != 3 {
			fmt.Println("need a hash to know what to get")
			return
		}
		fmt.Println("get call...")
		CliSend("get " + args[2])

	case "exit":
		fmt.Println("exit call...")
		CliSend("exit ")

	default:
		fmt.Println("Invalid input!")
	}
}

func CliSend(data string) {
	fmt.Println("Want to send:", data)
	connection, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Want to write")
	connection.Write([]byte(data))
	fmt.Println("waiting for response")
	response := make([]byte, 512)
	fmt.Println("response", response)
	connection.Read(response)
	fmt.Println(string(response))
	defer connection.Close()
}
