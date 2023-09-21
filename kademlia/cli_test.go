package kademlia

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

// TODO I don't understand how to test the cli it will not cooperate
func TestCli(t *testing.T) {

	mockServer, mockClient := net.Pipe()

	kademliaId := NewKademliaID("0000000000000000000000000000000000000000")
	address := "localhost:8000"
	contact := NewContact(kademliaId, address)
	network := CreateNetwork(&contact)
	node := CreateKademlia(&network)
	// socketAddress := "/tmp/echo.sock"

	// listener, err := net.Listen("unix", socketAddress)
	// if err != nil {
	// 	t.Errorf("recieved err [%s] on %s", err, socketAddress)
	// }
	// defer listener.Close()

	testWords := [2]string{"put test", "get 01234567890123456789"}

	for _, command := range testWords {
		args := strings.SplitN(command, " ", 2)
		fmt.Println("loop?", command)
		// connection, err := listener.Accept()
		// if err != nil {
		// 	t.Errorf("recieved err [%s] on listener.Accept()", err)
		// }
		// fmt.Println("connection accepted")

		// conn, err := net.Dial("unix", socketAddress)
		// if err != nil {
		// 	t.Errorf("[%s]", err)
		// 	return
		// }
		// conn.Write([]byte(command))
		// response := make([]byte, 128)
		// conn.Read(response)

		go CliHandler(mockServer, &node)

		fmt.Printf("command till mock: [%s]\n", command)
		mockClient.Write([]byte(command))
		fmt.Println("response?")

		response := make([]byte, 512)
		mockClient.Read(response)

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
		}

		if expectedAnswer != "not nil" && string(response) != expectedAnswer {
			t.Errorf("Cli test fail, expected [%s] but got [%s]", expectedAnswer, string(response))
		}
		if expectedAnswer == "not nil" && string(response) == "" {
			t.Errorf("Cli test fail, expected [%s] but got [%s]", expectedAnswer, string(response))
		}
		defer mockClient.Close()
		defer mockServer.Close()
	}

}
