package fileup

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type Server struct {
	RootDir string
	Origin  string // used for testing mostly
	err     error  // for testing
}

// root is the directory where the file will be saved.
// Func will exit with `log.Fatal` if the root direcoty name is not provided properly
//
// origin is the local addr.
// example : `localhost:8001`
func NewServer(root string, origin string) *Server {
	stat, err := os.Stat(root)
	if err != nil {
		log.Fatalf("NewServer: while opening %q: %s\n", root, err)
	}
	if !stat.IsDir() {
		log.Fatalf("NewServer: %q is a not a directory\n", root)
	}

	return &Server{
		RootDir: root,
		Origin:  origin,
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

func (s *Server) Serve() {
	http.Handle("/", websocket.Handler(s.HandleConn))

	if err := http.ListenAndServe(s.Origin, nil); err != nil {
		log.Fatal(err)
	}
}
