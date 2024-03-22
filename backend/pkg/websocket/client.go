package websocket

import (
	"fmt"
	"log"
	"sync"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	mu   sync.Mutex
	lng  string
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		c.lng = GetResponse("The following is text in a specific language. Please find this language, and output only the language name itself. The text to analyze is: " + string(p))
		fmt.Printf("Language: %+v\n", c.lng)
		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
		}
	}
