package websocket

import (
	"fmt"
	"log"
	"github.com/cjr03/MultilingualChat-AI/backend/pkg/ai"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			client.lng = "English"
			fmt.Println("Size of Connection Pool:", len(pool.Clients))
			for c := range pool.Clients {
				c.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool:", len(pool.Clients))
			for c := range pool.Clients {
				c.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}

		case message := <-pool.Broadcast:
			fmt.Println("Broadcasting message to clients")
			for c := range pool.Clients {
				message.Body = ai.GetResponse(
					"This is a message to someone who speaks: " + c.lng + ". Output this message translated into that language if necessary, or original if not. The message is: " + message.Body,
				)
				if err := c.Conn.WriteJSON(message); err != nil {
					fmt.Println("Error sending message:", err)
				}
			}
		}
	}
}
