package fileup

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	RootDir string
	err     error  // for testing
	origin  string // for testing
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(s.HandleConn).ServeHTTP(w, r)
}

func (s *Server) HandleConn(conn *websocket.Conn) {
	u := NewUpper(s.RootDir)
	u.conn = conn
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		for {
			_, err := conn.Write([]byte("ping"))
			if err != nil {
				conn.Close()
				log.Printf("conn form %q has been cloesd\n", conn.Request().RemoteAddr)
				break
			}
			time.Sleep(3 * time.Second)
		}
		done <- struct{}{}
	}()

	s.err = u.saveToFile()
	if s.err != nil {
		log.Println(s.err)
	}
	<-done
}

// origin will be selected randomly
func (s *Server) Serve(route string) {
	s.origin = randAddr()
	http.Handle(route, websocket.Handler(s.HandleConn))

	log.Printf("litening at %q\n", s.origin)
	if err := http.ListenAndServe(s.origin, nil); err != nil {
		log.Fatal(err)
	}
}

func randAddr() string {
	// Define a range for the random port (e.g., 1024 to 65535)
	minPort := 1024
	maxPort := 65535

	// Generate a random port within the specified range
	randomPort := rand.Intn(maxPort-minPort+1) + minPort

	return "localhost:" + strconv.Itoa(randomPort)
}
