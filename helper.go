package fileup

import (
	"crypto/sha256"
	"encoding/hex"
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

func (up Upper) createFile(name string) (*os.File, error) {
	name = filepath.Join(up.RootDir, name)

	_, err := os.Stat(name)
	// err == nil mans the file exists
	if err == nil {
		// for randomness
		name += "-" + strconv.Itoa(int(time.Now().UnixNano()))
	}

	return os.Create(name)
}

func checkFile(fileName string, expectedSum []byte) (Message, error) {
	file, err := os.Open(fileName)
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

	if sum != string(expectedSum) {
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