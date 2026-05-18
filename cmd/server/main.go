package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/vshishkar/simple-tcp/server"
)

func main() {
	server := server.MakeServer()
	err := rpc.Register(server)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RPC server listening on 9000")

	go server.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go rpc.ServeConn(conn)
	}

}
