package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wizsk/fileup"
	"golang.org/x/net/websocket"
)

func main() {
	d := http.FileServer(http.Dir("ui"))
	http.Handle("/", d)

	rootDir := "tmp"
	// igroing the err
	_ = os.Mkdir(rootDir, os.ModePerm)
	s := fileup.NewServer(rootDir, "")
	http.Handle("/ws", websocket.Handler(s.HandleConn))
	// starting server

	log.Println("serving at http://localhost:8002")
	if err := http.ListenAndServe(":8002", nil); err != nil {
		log.Fatal(err)
	}

}
