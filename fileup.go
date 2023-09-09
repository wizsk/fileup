package fileup

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Saver struct {
	UpDir           string
	UpRoute         string
	IncomePleateExt string // default is ".part"
	Err             ErrHandler
}

type ErrHandler func(w http.ResponseWriter, r *http.Request, httpMsg string, httpCode int, err error)

const (
	permFile = 0664
	permDir  = 0755
)

// upRoute is the "example/route/path"
func NewSaver(upRoute, upDir string, errHandeler ErrHandler) *Saver {
	if errHandeler == nil {
		errHandeler = func(w http.ResponseWriter, r *http.Request, httpMsg string, httpCode int, err error) {
			if httpMsg == "" {
				httpMsg = "sorry someting went wrong"
			}

			http.Error(w, httpMsg, httpCode)
			log.Printf("[Error] Request from %q: err: %v\n", r.RemoteAddr, err)
		}
	}

	return &Saver{
		UpRoute:         upRoute,
		UpDir:           upDir,
		IncomePleateExt: ".part",
		Err:             errHandeler,
	}
}

func (s *Saver) Handeler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.postHandler(w, r)

	case http.MethodPatch:
		s.patchHandler(w, r)

	case http.MethodHead:
		// do stuff

	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}

func (s *Saver) patchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/offset+octet-stream" {
		s.Err(w, r,
			"bad request", http.StatusBadRequest,
			fmt.Errorf(`wanted "application/offset+octet-stream" got %q`, r.Header.Get("Content-Type")))
		return
	}

	uploadOffset, err := uploadOffset(r)
	if err != nil {
		s.Err(w, r, "", http.StatusBadRequest, err)
		return
	}

	file, err := s.openFile(r)
	if err != nil {
		s.Err(w, r, "", http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		if fileStat.Size() != uploadOffset {
			s.Err(w, r, "", http.StatusBadRequest, err)
			return
		}
	}

	defer r.Body.Close()
	bytesWritten, err := io.Copy(file, r.Body)
	if err != nil {
		s.Err(w, r, "", http.StatusBadRequest, err)
		return
	}

	err = s.renameFile(file, r)
	if err != nil {
		s.Err(w, r, "", http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Upload-Offset", strconv.FormatInt(uploadOffset+bytesWritten, 10))
}

// renameFile cheks if the file is fully written if true
// then it will try to rename the file without uuuid;
// if a file already exists with that given name then it will keep the uuid like "foo.mp4-uuid"
func (s *Saver) renameFile(file *os.File, r *http.Request) error {
	upSize, err := uploadSize(r)
	if err != nil {
		return err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	if upSize != fileStat.Size() {
		return nil
	}

	fileName := strings.TrimPrefix(r.URL.Path, s.UpRoute)
	fileName = filepath.Join(s.UpDir, fileName)

	if _, err := os.Stat(fileName); err == nil {
		fileName = strings.TrimRight(file.Name(), s.IncomePleateExt)
	}

	return os.Rename(file.Name(), fileName)
}

// this cretes the file
func (s *Saver) postHandler(w http.ResponseWriter, r *http.Request) {
	fileName, err := s.getFilePath(r)
	if err != nil {
		s.Err(w, r, "", http.StatusBadRequest, err)
		return
	}

	_, err = os.Create(fileName)
	if err != nil {
		s.Err(w, r, "", http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// chekFile cheks the file exists or not
// and retuns a FileInfo || an err
func (s *Saver) openFile(r *http.Request) (*os.File, error) {
	filePath, err := s.getFilePath(r)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, permFile)
}

func (s *Saver) getFilePath(r *http.Request) (string, error) {
	fileName := strings.TrimPrefix(r.URL.Path, s.UpRoute)
	uuid := r.Header.Get("Upload-UUID")
	if uuid == "" {
		return "", errors.New("tempFileName: uuid is empty")
	}

	// name-uuid.incompleteExt
	fileName = fmt.Sprintf("%s-%s%s", fileName, uuid, s.IncomePleateExt)

	return filepath.Join(s.UpDir, fileName), nil
}
