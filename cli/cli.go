package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	args := os.Args

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

		fmt.Println("Pinging bootstrap...")
		CliSend("ping ")

	case "put":
		if len(args) != 3 {
			fmt.Println("Nothing to put")
			return
		}
		fmt.Println("put call...")
		CliSend("put " + args[2])

	case "get":
		if len(args) != 3 {
			fmt.Println("Need a hash to know what to get")
			return
		}
		_, err := hex.DecodeString(args[2])
		if err != nil || len(args[2]) != 40 {
			fmt.Println("Invalid hash (need 20 hex values)")
			return
		}
		fmt.Println("get call...")

		CliSend("get " + args[2])

	case "exit":
		fmt.Println("exit call...")
		CliSend("exit ")

	case "forget":
		if len(args) != 3 {
			fmt.Println("Need a hash to know what to forget")
			return
		}
		_, err := hex.DecodeString(args[2])
		if err != nil || len(args[2]) != 40 {
			fmt.Println("Invalid hash (need 20 hex values)")
			return
		}
		fmt.Println("forget call...")

		CliSend("forget " + args[2])

	default:
		fmt.Println("Invalid input!")
	}
}

func CliSend(data string) {
	connection, err := net.Dial("unix", "/tmp/echo.sock")
	defer connection.Close()
	if err != nil {
		log.Println(err)
		return
	}
	connection.Write([]byte(data))
	response := make([]byte, 1024)
	bytesRead, err := connection.Read(response)
	if err != nil {
        log.Println("Error reading data:", err.Error())
        return
    }
	fmt.Println(string(response[:bytesRead]))
}
