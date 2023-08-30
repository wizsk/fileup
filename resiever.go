package fileup

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// for ease of testing
type ConnReader interface {
	ReadMessage() (int, []byte, error)
	WriteJSON(interface{}) error
}

const (
	BUFF_SIZE = 2 * 1024 * 1024 // 2MB
)

var (
	ErrWrongMSG = errors.New("wrong message")
)

type Upper struct {
	RootDir  string
	BuffSize uint
	// Data is for storing msg data
	// problem is it will alloate and alocate a lot!
	// TODO:
	//		- Fix this || !
	//				   or not
	Data            []byte
	CurrentFile     io.WriteCloser
	CurrentFileName string
	wsUp            websocket.Upgrader
	createFile      func(*Upper) error
	openFile        func(name string) (io.ReadCloser, error)
}

type Message struct {
	IsError bool   `json:"is_error"`
	Body    string `json:"body"`
}

func NewUpper(root string) Upper {
	return Upper{
		RootDir:  root,
		BuffSize: BUFF_SIZE,
		wsUp: websocket.Upgrader{
			ReadBufferSize:  BUFF_SIZE,
			WriteBufferSize: BUFF_SIZE,
		},
		createFile: createFile,
		openFile: func(name string) (io.ReadCloser, error) {
			return os.Open(name)
		},
	}
}

func (up *Upper) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := up.wsUp.Upgrade(w, r, nil)
	if err != nil {
		// handle er
		return
	}
	defer conn.Close()

	up.getData(conn)
}

func (up *Upper) getData(conn ConnReader) error {
	// just incase of there is a pannic
	// or early return
	defer func() {
		if up.CurrentFile != nil {
			_ = up.CurrentFile.Close()
		}
	}()

	var msg int
	var err error

	for {
		// read the name
		msg, up.Data, err = conn.ReadMessage()
		if err != nil {
			// return fmt.Errorf("while reading conn: %v", err)
			return err
		}
		if msg != websocket.TextMessage {
			return ErrWrongMSG
		}

		err = up.createFile(up)
		if err != nil {
			// return fmt.Errorf("while creaing file: %v", err)
			return err
		}

		// while the message type is Binnary read the up.Data and save it to the
		// file else break and do post processing
		for {
			msg, up.Data, err = conn.ReadMessage()
			if err != nil {
				// return fmt.Errorf("while reading conn: %v", err)
				return err
			}

			// if msg type no textMessage then it's binnary object
			if msg != websocket.BinaryMessage {
				break
			}

			_, err = up.CurrentFile.Write(up.Data)
			if err != nil {
				// return fmt.Errorf("while writing up.Data to file: %v", err)
				return err
			}

		}

		// checing file intrigrity
		// fmt.Println("checksum", string(up.Buff))
		fileMsg, err := up.checkFile()
		if err != nil {
			return err
		}
		conn.WriteJSON(fileMsg)
		if err != nil {
			// TODO:
			//		- delete the file
			return err
		}

		// cleanin up
		up.CurrentFile.Close()
		up.CurrentFile = nil
		up.CurrentFileName = ""
		if err != nil {
			return err
		}
	}
}
