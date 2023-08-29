package fileup

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	BUFF_SIZE = 2 * 1024 * 1024 // 2MB
)

type Upper struct {
	RootDir  string
	BuffSize uint
	wsUp     websocket.Upgrader
}

type Message struct {
	IsError bool   `json:"is_error"`
	Body    string `json:"body"`
}

func NewUpper(root string) Upper {
	return Upper{
		RootDir:  root,
		BuffSize: BUFF_SIZE,
		wsUp: websocket.Upgrader{
			ReadBufferSize:  BUFF_SIZE,
			WriteBufferSize: BUFF_SIZE,
		},
	}
}

func (up Upper) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := up.wsUp.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	msg, name, err := conn.ReadMessage()
	if err != nil || msg != websocket.TextMessage {
		log.Println(err)
		return
	}

	f, err := up.createFile(string(name))
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	buff := make([]byte, up.BuffSize)

	for {
		msg, buff, err = conn.ReadMessage()
		if err != nil {
			log.Println(err)
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
		log.Println(err)
		// handle err
	}

	conn.WriteJSON(fileMsg)
	if err != nil {
		log.Println(err)
		// handle err
		return
	}
}
