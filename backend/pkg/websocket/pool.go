package websocket

import (
	"fmt"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

type Pool struct {
	Register  	chan *Client
	Unregister	chan *Client
	Clients   	map[*Client]bool
	Broadcast 	chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GetResponse(question string) string {
	apiKey := "API keys must be kept secret, please replace with your own"
	client := resty.New()
	response, err := client.R().
        SetAuthToken(apiKey).
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "model":      "gpt-3.5-turbo",
            "messages":   []interface{}{map[string]interface{}{"role": "system", "content": question}},
            "max_tokens": 50,
        }).
        Post(apiEndpoint)
	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}
	body := response.Body()
	var output ChatCompletionResponse
	err = json.Unmarshal(body, &output)
	if err != nil {
		return "error"
	}
	content := output.Choices[0].Message.Content
	return content
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			client.lng = "English";
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				message.Body = GetResponse("This is a message to someone who speaks: " + client.lng + ". Output this message translated into that language if necessary, or the original message if it does not need to be translated. The message is: " + message.Body)
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
