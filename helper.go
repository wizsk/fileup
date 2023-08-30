package fileup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var (
	ErrorSHA256DontMatch = errors.New("sha256 dosen't match")
)

// createFile takes the *Upper and reads the data
// figures out the file name and creates the file
// and saves it to the *Upper
//
// if error happends it retusn it
func createFile(up *Upper) error {
	var filename struct {
		FileName string `json:"file_name"`
	}
	err := json.Unmarshal(up.Data, &filename)
	if err != nil {
		return err
	}

	name := filename.FileName
	name = filepath.Join(up.RootDir, name)

	_, err = os.Stat(name)
	// err == nil mans the file exists
	if err == nil {
		// for randomness
		// TODO:
		// 		- Find a better solution
		name += "-" + strconv.Itoa(int(time.Now().UnixNano()))
	}

	up.CurrentFile, err = os.Create(name)
	up.CurrentFileName = name
	return err
}

func (up *Upper) checkFile() (Message, error) {
	var expectedSum struct {
		ShaSum string `json:"checksum"`
	}

	err := json.Unmarshal(up.Data, &expectedSum)
	if err != nil {
		return Message{IsError: true, Body: "something unexpected happened"}, err
	}

	var file io.ReadCloser
	file, err = up.openFile(up.CurrentFileName)
	if err != nil {
		msg := Message{
			IsError: true,
			Body:    "could't save file",
		}
		return msg, err
	}

	sum, err := calculateSHA256Checksum(file)
	if err != nil {
		msg := Message{
			IsError: true,
			Body:    "could't save file",
		}
		return msg, err
	}

	if sum != string(expectedSum.ShaSum) {
		msg := Message{
			IsError: true,
			Body:    ErrorSHA256DontMatch.Error(),
		}
		return msg, err
	}

	msg := Message{
		IsError: false,
		Body:    "file upload successfull",
	}

	return msg, nil
}

func calculateSHA256Checksum(data io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", err
	}
	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum), nil
}

// for debuggin perposes
func PrintMemUsage(r string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v%s", m.NumGC, r)
}
