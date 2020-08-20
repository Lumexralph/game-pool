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

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func serveWs(pool *game.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
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

func SetupRoutes() {
	pool := game.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}
