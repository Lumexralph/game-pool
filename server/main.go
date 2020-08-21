package main

import (
	"fmt"
	"net/http"
	"server/handler"
)

func main() {
	fmt.Println("Distributed Network Game v0.01")
	handler.SetupRoutes()
	http.ListenAndServe(":8080", nil)
}
