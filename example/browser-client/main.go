package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wizsk/fileup"
)

func main() {
	port := ":8001"
	updir := "uploads"
	upRoute := "/api/upload/"
	os.Mkdir("uploads", os.ModePerm)

	s := fileup.NewSaver(upRoute, updir, nil)

	http.Handle("/", http.FileServer(http.Dir("ui")))
	http.HandleFunc(upRoute, s.Handeler)

	log.Println("litenting at http://localhost" + port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
