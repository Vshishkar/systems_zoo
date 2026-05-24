package server

import (
	"log"
	"time"
)

type SendMessageArgs struct {
	ClientId int
	Text     string
}

type SendMessageResponse struct {
	Success bool
}

func (s *Server) SendMessage(req SendMessageArgs, res *SendMessageResponse) error {
	log.Printf("Received message from client %d: %s", req.ClientId, req.Text)
	s.RecordMessage(req)

	res.Success = true
	return nil
}

func (s *Server) RecordMessage(req SendMessageArgs) {
	s.lock.Lock()
	defer s.lock.Unlock()

	message := Message{
		id:        len(s.messages) + 1,
		text:      req.Text,
		timestamp: time.Now(),
		authorId:  req.ClientId,
	}
	s.messages = append(s.messages, message)
	s.broadcastCh <- message
}
