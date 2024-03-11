package websocket

import (
	"fmt"
	"log"
	"sync"
	"github.com/gorilla/websocket"
)

import (
    "encoding/json"
    "github.com/go-resty/resty/v2"
)

const (
    apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

var langauge_translator string = "The following message should be translated into " + language + ". If it is already in " + language + " please output only the message. If it must be translated, please output only the translated message as your answer. The message is: "
var language string = "English"

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	mu   sync.Mutex
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func GetResponse(question String){
	apiKey := ""
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
    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        fmt.Println("Error while decoding JSON response:", err)
        return
    }
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content
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
		message := Message{Type: messageType, Body: GetResponse(language_translator + string(p))}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
