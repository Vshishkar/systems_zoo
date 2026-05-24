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

type MessageDto struct {
	Id        int
	Text      string
	Timestamp time.Time
	AuthorId  int
}

type RegisterClientResponse struct {
	Messages []MessageDto
	Id       int
}

func (s *Server) RegisterClient(req RegisterClientArgs, res *RegisterClientResponse) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

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
	res.Messages = make([]MessageDto, len(s.messages))
	for i, msg := range s.messages {
		res.Messages[i] = MessageDto{
			Id:        msg.id,
			Text:      msg.text,
			Timestamp: msg.timestamp,
			AuthorId:  msg.authorId,
		}
	}

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
