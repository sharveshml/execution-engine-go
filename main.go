package main

import (
	"log"
	"net/http"

	"index/pty"
)

func main() {
	http.HandleFunc("/ws", pty.InitWebSocket)
	log.Println("Server starting on port: 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
