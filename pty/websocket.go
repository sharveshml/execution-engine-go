package pty

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not upgrade connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	replId := strings.Split(r.Host, ":")[0]

	if replId == "" {
		http.Error(w, "replId not found", http.StatusBadRequest)
		return
	}

	wd, _ := os.Getwd()
	fmt.Println(wd)
	rootContent, err := FetchDir(wd + "/pty" + "/workspace")

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Cannot get the files in in workspace", http.StatusBadRequest)
		return
	}

	HandleWebSocket(conn, replId, rootContent)
}

func HandleWebSocket(conn *websocket.Conn, replId string, file *[]File) {

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Read Error: ", err)
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Unmarshall error: ", err)
			continue
		}

		event, ok := msg["event"].(string)
		if !ok {
			log.Println("Invalid event type")
			continue
		}

		log.Println("Received event: ", event)

		switch event {
		case "disconnect":
			handleDisconnect(conn)
		case "fetchDir":
			handleFetchDir(conn)
		case "fetchContent":
			path, _ := msg["path"].(string)
			handleFetchContent(conn, path)
		case "updateContent":
			path, _ := msg["path"].(string)
			content, _ := msg["content"].(string)
			handleUpdateContent(conn, path, content, replId)
		case "requestTerminal":
			handleRequestTerminal(conn, replId)
		case "terminalData":
			data, _ := msg["data"].(string)
			handleTerminalData(conn, data)
		default:
			log.Println("Unknown event: ", event)
		}
	}
}
