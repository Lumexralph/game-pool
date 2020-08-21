package main

import (
	"fmt"
	"net/http"
	"server/handler"
)

func main() {
	fmt.Println("Distributed Chat App v0.01")
	handler.SetupRoutes()
	http.ListenAndServe(":8080", nil)
}
