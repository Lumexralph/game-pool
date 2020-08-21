package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	game "server/pool"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Upgrade upgrades the HTTP request to a websocket connection.
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

// serveWs handler listens on the /ws endpoint for websocket connection.
func serveWs(pool *game.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := Upgrade(w, r)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%+v\n", err)
	}

	// create new client
	client := &game.Client{
		ID:   game.GenerateClientID(),
		Conn: conn,
		Pool: pool,
	}

	// when a client newly joins
	pool.Register <- client
	client.Read()
}

// SetupRoutes handles the creation of a new game pool and the router handlers.
func SetupRoutes() {
	pool := game.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}
