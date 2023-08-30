package main

import (
	"net/http"

	"github.com/wizsk/fileup"
)

func main() {
	// go func() {
	// 	for range time.Tick(time.Second) {
	// 		fileup.PrintMemUsage("\r")
	// 	}
	//
	// }()
	up := fileup.NewUpper("tmp")
	d := http.FileServer(http.Dir("ui"))
	http.Handle("/", d)
	http.HandleFunc("/ws", up.HandleWebSocket)
	http.ListenAndServe(":8002", nil)

}
