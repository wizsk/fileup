package fileup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/websocket"
)

const CHUNK_SIZE int64 = 1024 * 1024 * 2 // 2MB

var ErrCheckSumDontMatch = errors.New("checum don't match")

type Upper struct {
	conn     *websocket.Conn
	root     string
	fileInfo FileInfo
	buff     []byte
}

func NewUpper(rootDir string) *Upper {
	if rootDir == "" {
		log.Fatal("rootdir cant't be emty")
	}
	return &Upper{
		buff: make([]byte, CHUNK_SIZE),
		root: rootDir,
	}

}

func (u *Upper) createFile() (*os.File, error) {
	var read int
	var err error

	for read <= 0 {
		read, err = u.conn.Read(u.buff)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	// read file info
	u.fileInfo = FileInfo{}
	if err := json.Unmarshal(u.buff[:read], &u.fileInfo); err != nil {
		return nil, err
	}

	_, err = u.conn.Write([]byte(`{"status": "file creatiion successfull success"}`))
	if err != nil {
		return nil, err
	}

	u.fileInfo.path = filepath.Join(u.root, u.fileInfo.Name)
	return createCniqueFile(u.fileInfo.path)
}

func createCniqueFile(path string) (*os.File, error) {
	_, err := os.Stat(path)
	// err nill file existss
	if err == nil {
		path += fmt.Sprintf("-%d", time.Now().UnixNano())
	}

	return os.Create(path)
}

func (u *Upper) saveToFile() error {
	for {
		file, err := u.createFile()
		if err != nil {
			return err
		}
		defer file.Close()

		err = save(u.conn, file, u.fileInfo.Size)
		if err != nil {
			return err
		}
		file.Close()
		err = u.checkFile()
		if err != nil {
			// ignoring the err for now
			_ = writeJson(u.conn, StatusMsg{Error: true, Type: "error", Body: err.Error()})
			return err
		}
	}
}

func (u *Upper) checkFile() error {
	var read int
	var err error

	for read <= 0 {
		read, err = u.conn.Read(u.buff)
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

	var acualSum Sha256
	if err = json.Unmarshal(u.buff[:read], &acualSum); err != nil {
		return err
	}

	matched := acualSum.Sum != sum
	if matched {
		return ErrCheckSumDontMatch
	}

	return nil
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

		// read only the file part
		if totalRead+CHUNK_SIZE >= size {
			// read only until the file not more than that
			buf = buf[:size-totalRead]
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

func writeJson(w io.Writer, i interface{}) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
