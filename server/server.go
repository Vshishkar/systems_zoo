package server

import (
	"fmt"
	"net/rpc"
	"time"
)

type Server struct {
	messages    []Message
	clients     map[int]Client
	connections map[int]*rpc.Client
	clientMaxId int
	broadcastCh chan Message
}

type Client struct {
	id   int
	port int
}

type Message struct {
	id        int
	text      string
	timestamp time.Time
	authorId  int
}

func (s *Server) EstablishConnection(clientId int) {
	client, err := rpc.Dial("tcp", fmt.Sprintf(":%d", s.clients[clientId].port))
	if err != nil {
		fmt.Printf("Failed to connect to client %d: %v\n", clientId, err)
		return
	}
	s.connections[clientId] = client
}

func MakeServer() *Server {
	return &Server{
		messages:    []Message{},
		clients:     make(map[int]Client),
		connections: make(map[int]*rpc.Client),
		clientMaxId: 0,
		broadcastCh: make(chan Message),
	}
}

func (s *Server) Start() {
	for {
		select {
		case message := <-s.broadcastCh:
			s.broadcastMessage(message)
		}
	}
}
