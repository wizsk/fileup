package fileup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/websocket"
)

const origin = "localhost:8001"

const CHUNK_SIZE int64 = 1024 * 1024 * 2 // 2MB

var SaveFileDir = "tmp"

var ErrCheckSumDontMatch = errors.New("checum don't match")

type Upper struct {
	conn     *websocket.Conn
	fileInfo FileInfo
	err      error
}

func (u *Upper) handleConn(ws *websocket.Conn) {
	u.conn = ws
	defer ws.Close()

	err := u.saveToFile()
	u.err = err
	if err != nil {
		log.Println(err)
	}
}

func (u *Upper) Serve() {
	http.Handle("/", websocket.Handler(u.handleConn))

	if err := http.ListenAndServe(origin, nil); err != nil {
		log.Fatal(err)
	}
}

func (u *Upper) createFile() (*os.File, error) {
	buf := make([]byte, CHUNK_SIZE)
	var read int
	var err error

	fmt.Println("?", read, buf)
	for read <= 0 {
		read, err = u.conn.Read(buf)
		fmt.Println("?", read, buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	fmt.Println("?", read, buf)
	// read file info
	u.fileInfo = FileInfo{}
	if err := json.Unmarshal(buf[:read], &u.fileInfo); err != nil {
		return nil, err
	}

	_, err = u.conn.Write([]byte(`{"status": "file creatiion successfull success"}`))
	if err != nil {
		return nil, err
	}

	u.fileInfo.path = filepath.Join(SaveFileDir, u.fileInfo.Name)
	return os.Create(u.fileInfo.path)
}

func (u *Upper) checkFile() error {
	buf := make([]byte, CHUNK_SIZE)
	var read int
	var err error

	for read <= 0 {
		read, err = u.conn.Read(buf)
		if err != nil {
			return err
		}
	}

	file, err := os.Open(u.fileInfo.path)
	if err != nil {
		return err
	}
	defer file.Close()

	sum, err := calculateSHA256Checksum(file)
	if err != nil {
		return err
	}

	if string(buf[:read]) != sum {
		return ErrCheckSumDontMatch
	}

	return nil
}

func (u *Upper) saveToFile() error {
	fmt.Println("0")
	for {

		fmt.Println("0+11")
		file, err := u.createFile()
		if err != nil {
			fmt.Println("e1")
			return err
		}
	panic("f")
		fmt.Println("1")

		err = save(u.conn, file, u.fileInfo.Size)
		if err != nil {
			return err
		}
		fmt.Println("2")
		file.Close()
		err = u.checkFile()
		if err != nil {
			return err
		}
		fmt.Println("3")
	}
}

// save reads from the conn,  writes it to the w
// and retuns error
func save(conn io.Reader, w io.Writer, size int64) error {
	var totalRead int64
	buf := make([]byte, CHUNK_SIZE)
	for {
		read, err := conn.Read(buf)
		// eof err should not occur here
		if err != nil {
			return nil
		}
		_, err = w.Write(buf[:read])
		if err != nil {
			return nil
		}

		totalRead += int64(read)

		if totalRead >= size {
			return nil
		}

		// hacky math
		if totalRead+CHUNK_SIZE >= size {
			buf = buf[:size-totalRead+CHUNK_SIZE]
		}
	}

}

func calculateSHA256Checksum(data io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", err
	}
	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum), nil
}
