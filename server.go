package fileup

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type Server struct {
	RootDir string
	err     error // for testing
}

// root is the directory where the file will be saved.
// Func will exit with `log.Fatal` if the root direcoty name is not provided properly
func NewServer(root string) *Server {
	stat, err := os.Stat(root)
	if err != nil {
		log.Fatalf("NewServer: while opening %q: %s\n", root, err)
	}
	if !stat.IsDir() {
		log.Fatalf("NewServer: %q is a not a directory\n", root)
	}

	return &Server{
		RootDir: root,
	}
}

func (s *Server) HandleConn(conn *websocket.Conn) {
	u := NewUpper(s.RootDir)
	u.conn = conn
	defer conn.Close()

	s.err = u.saveToFile()
	if s.err != nil {
		log.Println(s.err)
	}
}

// Serve serves the websocket
// route is the path where the http server will use
// origin is the local addr.  example : `localhost:8001`
func (s *Server) Serve(route, origin string) {
	http.Handle(route, websocket.Handler(s.HandleConn))

	log.Printf("litening at %q\n", origin)
	if err := http.ListenAndServe(origin, nil); err != nil {
		log.Fatal(err)
	}
}
