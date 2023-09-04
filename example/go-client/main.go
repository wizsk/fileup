package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wizsk/fileup"
	"golang.org/x/net/websocket"
)

const origin = "localhost:8001"

func main() {
	// ignoring erross
	_ = os.Mkdir("tmp", os.ModePerm)

	// starting real server :)
	s := fileup.NewServer("tmp")
	// go s.Serve("/", origin) // you can use this
	// or
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// you can do your suff here and
		websocket.Handler(s.HandleConn).ServeHTTP(w, r)
	})
	go func() {
		if err := http.ListenAndServe(origin, nil); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := websocket.Dial("ws://"+origin, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// write to conn
	handle(conn)

	// waiting for the server
	time.Sleep(2 * time.Second)

	fmt.Println("file was written to './tmp'")
}

func handle(conn *websocket.Conn) {
	// this can be with out loop too!
	for i := 0; i <= 10; i++ {
		// garbadge data
		gd := make([]byte, 1024*1024*10) // 10mb data
		_, err := rand.Read(gd)
		if err != nil {
			log.Fatal(err)
		}

		sendFileInfo(conn, "random-file-name", int64(len(gd)))

		// stream the data
		if _, err = io.Copy(conn, bytes.NewReader(gd)); err != nil {
			log.Fatal(err)
		}
		sendSha256(conn, bytes.NewReader(gd))
		readShaStatus(conn)
	}
}

func sendFileInfo(conn *websocket.Conn, name string, lenth int64) {
	_, err := conn.Write(jsonify(
		fileup.FileInfo{
			Name: name,
			Size: lenth,
		}))

	if err != nil {
		log.Fatal(err)
	}

}

func sendSha256(conn *websocket.Conn, r io.Reader) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		log.Fatal(err)
	}
	jsn := jsonify(fileup.Sha256{Sum: fmt.Sprintf("%x", h.Sum(nil))})
	if _, err := conn.Write(jsn); err != nil {
		log.Fatal(err)
	}
}

func readShaStatus(conn *websocket.Conn) {
	buf := make([]byte, fileup.CHUNK_SIZE)

	// read is blocking
	read, err := conn.Read(buf)

	stat := fileup.StatusMsg{}
	// err := websocket.Message.Receive(conn, &stat)
	if err = json.Unmarshal(buf[:read], &stat); err != nil {
		log.Fatal(err)
	}

	if stat.Type != "sha256" || stat.Error {
		log.Fatalf("%+v\n", stat)
	}
}

func jsonify(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
