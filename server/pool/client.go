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
	Player     bool `json:"player"`
	TotalScore int8 `json:"totalScore"`
	lowerBound uint8
	upperBound uint8
}

// String method to help viewing the client in I/O stream,
// and implement the Stringer interface for formatting client.
func (c *Client) String() string {
	return fmt.Sprintf("{ name: %q, id: %q, TotalScore: %d }", c.Name, c.ID, c.TotalScore)
}

// Read method will continually read from the connection stream to the pool Ping channel and
// handle when a client connection stream is closed.
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		if err := c.Conn.Close(); err != nil {
			log.Println("client: error closing connection")
		}
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

// GenerateClientID is an utility function to generate identifiers for a client connection
func GenerateClientID() string {
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	return id
}
