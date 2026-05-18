package server

import (
	"fmt"
	"log"
	"time"
)

type RegisterClientArgs struct {
	Id   *int
	Port int
}

type RegisterClientResponse struct {
	Id int
}

func (s *Server) RegisterClient(req RegisterClientArgs, res *RegisterClientResponse) error {
	var clientId int
	if req.Id == nil {
		s.clientMaxId++
		clientId = s.clientMaxId
	} else {
		clientId = *req.Id
	}

	s.clients[clientId] = Client{
		id:   clientId,
		port: req.Port,
	}

	log.Printf("Registered client %d on port %d", clientId, req.Port)
	res.Id = clientId

	go func() {
		s.broadcastCh <- Message{
			id:        0,
			text:      fmt.Sprintf("Client %d has joined the chat", clientId),
			timestamp: time.Now(),
			authorId:  0,
		}
	}()

	return nil
}
