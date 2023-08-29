package fileup

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
)

const (
	BUFF_SIZE = 2 * 1024 * 1024 // 2MB
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  BUFF_SIZE,
	WriteBufferSize: BUFF_SIZE,
}

type Upper struct {
	RootDir  string
	BuffSize uint
}

type Message struct {
	IsError bool   `json:"is_error"`
	Body    string `json:"body"`
}

func NewUpper(root string) Upper {
	return Upper{
		RootDir:  root,
		BuffSize: BUFF_SIZE,
	}
}

func (up Upper) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	msg, name, err := conn.ReadMessage()
	if err != nil || msg != websocket.TextMessage {
		return
	}

	f, err := os.Create(filepath.Join(up.RootDir, string(name)))
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	buff := make([]byte, BUFF_SIZE)

	for {
		msg, buff, err = conn.ReadMessage()
		if err != nil {
			return
		}

		// if msg type no textMessage then it's binnary object;;
		// but i do think if text file is sent then it will be a pain :)
		// i guess noo
		if msg != websocket.BinaryMessage {
			break
		}
		_, err = f.Write(buff)
		if err != nil {
			log.Println(err)
			return
		}
	}

	fileMsg, err := checkFile(f.Name(), buff)
	if err != nil {
		// handle err
	}

	conn.WriteJSON(fileMsg)
	if err != nil {
		// handle err
		return
	}
}
