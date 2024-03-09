package websocket

import (
	"fmt"
	"log"
	"sync"
	"github.com/gorilla/websocket"
)

import (
	"context"
	"os"
	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
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

func GetResponse(c *Client, client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		c.Pool.Broadcast <- Message{Type: 1, Body: resp.Choices[0].Text}
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
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
		message := Message{Type: messageType, Body: string(p)}
		viper.SetConfigFile(".env")
		viper.ReadInConfig()
		//apiKey := viper.GetString("API_KEY")
		apiKey := os.Getenv("API_KEY")
		fmt.Println("API Key:", apiKey)
		if apiKey == "" {
			panic("Missing API KEY")
		}
		ctx := context.Background()
		client := gpt3.NewClient(apiKey)
		GetResponse(c, client, ctx, langauge_translator + message.Body)
		fmt.Printf("Message Received: %+v\n", message)
	}
}
