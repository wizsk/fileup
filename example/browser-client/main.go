package main

import (
	"log"
	"net/http"

	"github.com/wizsk/fileup"
)

func main() {
	port := ":8001"
	updir := "uploads"
	upRoute := "/api/upload/"

	// you can do this
	// os.Mkdir("uploads", os.ModePerm)
	// s := fileup.NewSaver(upRoute, updir)

	// or this
	s, err := fileup.NewSaverMkdir(upRoute, updir)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("ui")))
	http.HandleFunc(upRoute, s.Handeler)

	log.Println("litenting at http://localhost" + port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
