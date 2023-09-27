package kademlia

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func setupNode() *Kademlia {
	kademliaId := NewKademliaID("0000000000000000000000000000000000000000")
	address := "localhost:8000"
	contact := NewContact(kademliaId, address)
	network := CreateNetwork(&contact)
	node := CreateKademlia(&network)
	return &node

}

func TestCliPut(t *testing.T) {

	mockServer, mockClient := net.Pipe()

	node := setupNode()

	var savedHash string
	testWords := [5]string{"put test", "put ", "get 0123456789ABCDEF0123456789ABCDEF01234567", "get hash", "ex "}

	for _, command := range testWords {
		args := strings.SplitN(command, " ", 2)

		go CliHandler(mockServer, node)

		if command == "get hash" {
			command = "get " + savedHash
		}

		mockClient.Write([]byte(command))

		response := make([]byte, 512)
		mockClient.Read(response)

		time.Sleep(500 * time.Millisecond)

		var expectedAnswer string
		switch args[0] {
		case "ping":
			expectedAnswer = "PONG"
		case "put":
			expectedAnswer = "not nil"
		case "get":
			expectedAnswer = "not nil"
		case "exit":
			expectedAnswer = "exiting"
		default:
			expectedAnswer = "not nil"
		}

		if command == "put test" {
			savedHash = string(response)
		}

		fmt.Printf("---------------command [%s] räspåns [%s] [%v]\n", command, string(response), response == nil)

		if expectedAnswer != "not nil" && string(response) != expectedAnswer {
			t.Errorf("Cli test fail for command [%s], expected [%s] but got [%s]", command, expectedAnswer, string(response))
		}
		if expectedAnswer == "not nil" && string(response) == "" {
			t.Errorf("Cli test fail, expected [%s] but got [%s]", expectedAnswer, string(response))
		}
		defer mockClient.Close()
		defer mockServer.Close()
	}

}
