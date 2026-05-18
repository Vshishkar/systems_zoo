package server

import (
	"log"

	"github.com/vshishkar/simple-tcp/client"
)

func (s *Server) broadcastMessage(message Message) {
	for clientId := range s.clients {
		connection, ok := s.connections[clientId]
		if !ok {
			s.EstablishConnection(clientId)
			connection = s.connections[clientId]
		}

		if connection == nil {
			log.Printf("No connection available for client %d, skipping message broadcast\n", clientId)
			continue
		}

		log.Printf("Broadcasting message to client %d: %s\n", clientId, message.text)
		err := connection.Call("Client.ReceiveMessage", client.ReceiveMessageArgs{
			ClientId: message.authorId,
			Text:     message.text,
		}, &client.ReceiveMessageResponse{})
		if err != nil {
			log.Printf("Failed to send message to client %d: %v\n", clientId, err)
		}
	}
}
