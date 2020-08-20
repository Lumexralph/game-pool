package pool

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID         string          `json:"id"`
	Conn       *websocket.Conn `json:"-"`
	Pool       *Pool           `json:"-"`
	mu         sync.Mutex
	Name       string `json:"name"`
	player     bool
	TotalScore int8 `json:"totalScore"`
	lowerBound uint8
	upperBound uint8
}

// to help viewing the client in I/O stream
// implement the Stringer interface
func (c *Client) String() string {
	return fmt.Sprintf("{ name: %q, id: %q, TotalScore: %d }", c.Name, c.ID, c.TotalScore)
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}

		// check if the user wants to play
		c.Pool.Ping <- msg
		fmt.Printf("Message Received: %+v\n", msg)
	}
}

func GenerateClientID() string {
	uuid := strings.Replace(uuid.New().String(), "-", "", -1)
	return uuid
}
