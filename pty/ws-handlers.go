package pty

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var terminalManager = NewTerminalManager()

func HandleDisconnect(conn *websocket.Conn) {
	log.Println("Client disconnected")
	conn.Close()
}

func HandleFetchDir(conn *websocket.Conn, replId string) {
	log.Println("Fetching directory")
	dir, _ := os.Getwd()
	dirPath := filepath.Join(dir, "pty", "workspace")
	fmt.Println(dirPath)
	files, err := FetchDir(dirPath)

	if err != nil {
		log.Println("FetchDir error: ", err)
		return
	}
	conn.WriteJSON(files)
}

func HandleFetchContent(conn *websocket.Conn, path string) {
	data := FetchFileContent(path)

	conn.WriteJSON(string(data))
}

func HandleUpdateContent(conn *websocket.Conn, path string, content string, replId string) {
	log.Println("File path: ", path)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Println("Error writing file: ", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("File updated successfully"))
}

func HandleRequestTerminal(conn *websocket.Conn, replId string) {
	terminalID := uuid.New().String()

	_, err := terminalManager.CreatePty(terminalID, replId, func(data string, terminalID string) {
		conn.WriteJSON(map[string]interface{}{
			"event":      "terminal",
			"terminalID": terminalID,
			"data":       data,
		})
	})

	if err != nil {
		log.Printf("Failed to create PTY: %v", err)
		conn.WriteJSON(map[string]interface{}{
			"event": "terminalError",
			"error": "Failed to create terminal",
		})
		return
	}

	conn.WriteJSON(map[string]interface{}{
		"event":      "terminalCreated",
		"terminalID": terminalID,
	})
}

func HandleTerminalData(conn *websocket.Conn, data string, terminalID string) {
	if err := terminalManager.Write(terminalID, data); err != nil {
		log.Printf("Error writing to terminal %s: %v", terminalID, err)
		conn.WriteJSON(map[string]interface{}{
			"event": "terminalError",
			"error": err.Error(),
		})
	}
}
