package pty

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func HandleDisconnect(conn *websocket.Conn) {
	log.Println("Client disconnected")
	conn.Close()
}

func HandleFetchDir(conn *websocket.Conn, replId string) {
	log.Println("Fetching directory")
	dirPath := filepath.Join("os.Getwd()", "pty", "workspace")
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
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory: ", err)
		return
	}
	filePath := filepath.Join(dir, "pty", "workspace")
	log.Println("File path: ", filePath)

	if err := SaveFile(filePath, content); err != nil {
		log.Println("Error updating the file content", err)
	}
}

func HandleRequestTerminal(conn *websocket.Conn, replId string) {
	id := uuid.New()
	NewTerminalManager().CreatePty(id.string, "sharvesh", func(data string, id int) {
		conn.WriteJSON(map[string]interface{}{
			"event": "terminal",
			"data":  data,
		})
	})
}

func HandleTerminalData(conn *websocket.Conn, data string) {
	TerminalManager.Write(conn, data)
}
