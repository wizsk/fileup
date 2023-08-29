package main

import (
	"net/http"

	"github.com/wizsk/fileup"
)

func main() {
	up := fileup.NewUpper("tmp")
	d := http.FileServer(http.Dir("ui"))
	http.Handle("/", d)
	http.HandleFunc("/ws", up.handleWebSocket)
	http.ListenAndServe(":8002", nil)
}
