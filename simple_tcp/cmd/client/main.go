package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"

	"github.com/vshishkar/simple-tcp/client"
	"github.com/vshishkar/simple-tcp/server"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <message>")
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Invalid port number")
	}

	rpc_client, err := rpc.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}

	request := server.RegisterClientArgs{
		Port: port,
	}

	response := server.RegisterClientResponse{}
	err = rpc_client.Call("Server.RegisterClient", request, &response)
	if err != nil {
		log.Fatal(err)
	}

	main_client := &client.Client{
		Id:   response.Id,
		Port: port,
	}

	err = rpc.Register(main_client)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RPC server listening on", port)
	go acceptRPCConnections(listener)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter messages to send (type 'exit' to quit):")

	for _, message := range response.Messages {
		fmt.Printf("%d: %s\n", message.AuthorId, message.Text)
	}

	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}

		err = rpc_client.Call("Server.SendMessage", server.SendMessageArgs{
			ClientId: main_client.Id,
			Text:     text,
		}, &server.SendMessageResponse{})
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading from stdin:", err)
	}

}

func acceptRPCConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}
