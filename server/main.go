package main

import (
	"fmt"
	"net/http"
	"server/handler"
)

// // client play
// type ClientPlay struct {
// 	client *Client
// 	plays  [][2]uint8
// }

func main() {
	fmt.Println("Distributed Chat App v0.01")
	handler.SetupRoutes()
	http.ListenAndServe(":8080", nil)
}
