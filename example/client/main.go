package main

import (
	"encoding/json"
	"log"

	"github.com/wizsk/fileup"
	"golang.org/x/net/websocket"
)

const origin = "localhost:8001"

func main() {
	u := &fileup.Upper{}
	// starting real server :)
	go u.Serve()

	w := "ws://" + origin
	conn, err := websocket.Dial(w, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	f := fileup.FileInfo{
		Name: "foo",
		Size: 12312313,
	}
	fj, err := json.Marshal(f)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(fj)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("skldflksdfjlksdfjlkdjfjkldsjflkdsjfkldsjflkdsjfkldsjfkldsjflkjskdfjklsdfjklsdfjklsdjklfdsklfdjskl"))
	if err != nil {
		log.Fatal(err)
	}
}
