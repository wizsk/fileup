package fileup

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

const origin = "localhost:8001"
const saveFileDir = "tmp"
const readFileDir = "_tmp"

func TestHandeler(t *testing.T) {
	checkForTestFiels(t)
	s := NewServer(saveFileDir, origin)
	go s.Serve() // starting real server :)

	w := "ws://" + origin
	conn, err := websocket.Dial(w, "", origin)
	Nil(t, err)
	defer conn.Close()

	for _, fileName := range readDirecotory(t, readFileDir) {
		// fileName :=
		f, err := getFileInfo(fileName)
		Nil(t, err)
		fj, err := json.Marshal(f)
		Nil(t, err)
		_, err = conn.Write(fj)
		Nil(t, err)

		file, err := os.Open(fileName)
		Nil(t, err)
		_, err = io.Copy(conn, file)
		file.Close()
		Nil(t, err)

		file, err = os.Open(fileName)
		Nil(t, err)
		sha, err := calculateSHA256Checksum(file)
		file.Close()
		Nil(t, err)
		fj, err = json.Marshal(Sha256{Sum: sha})
		Nil(t, err)
		_, err = conn.Write(fj)
		Nil(t, err)
		t.Log(sha)
	}

	// wait for the conn and so on
	time.Sleep(5 * time.Second)
	// reading the err form server
	Nil(t, s.err)
}

func getFileInfo(path string) (FileInfo, error) {
	s, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Name: s.Name(),
		Size: s.Size(),
	}, err
}

func Nil(t *testing.T, v any) {
	t.Helper()
	if v != nil {
		t.Error(v)
		t.FailNow()
	}
}

func readDirecotory(t *testing.T, dir string) []string {
	t.Helper()
	dirE, err := os.Open(dir)
	Nil(t, err)

	dirlis, err := dirE.ReadDir(0)
	Nil(t, err)

	var fileNames []string
	for _, f := range dirlis {
		if !f.IsDir() {
			fileNames = append(fileNames, filepath.Join(dir, f.Name()))
		}
	}

	return fileNames
}

func checkForTestFiels(t *testing.T) {
	err := os.Mkdir(saveFileDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		Nil(t, err)
	}

	err = os.Mkdir(readFileDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		Nil(t, err)
	}

	if len(readDirecotory(t, readFileDir)) == 0 {
		t.Errorf("please put some files in the %q directory which will be used for testing", readFileDir)
		t.FailNow()
	}
}
