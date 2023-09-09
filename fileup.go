package fileup

import (
	"errors"
	"fmt"
	"io"
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
}

const (
	uploadRoute    = "/api/upload/"
	uploadDir      = "uploads"
	incompleateExt = ".part"

	PermFile = 0664
	PermDir  = 0755
)

func NewSaver(upRoute, upDir string) *Saver {
	return &Saver{
		UpRoute:         upRoute,
		UpDir:           upDir,
		IncomePleateExt: ".part",
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("ui")))
	s := NewSaver(uploadRoute, uploadDir)

	http.HandleFunc(uploadRoute, s.Handeler)
	http.ListenAndServe(":8080", nil)
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
		logErr(w, r, fmt.Errorf(`wanted "application/offset+octet-stream" got %q`, r.Header.Get("Content-Type")))
		return
	}

	uploadOffset, err := uploadOffset(r)
	if err != nil {
		logErr(w, r, err)
		return
	}

	file, err := s.openFile(uploadRoute, r)
	if err != nil {
		logErr(w, r, err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		if fileStat.Size() != uploadOffset {
			logErr(w, r, fmt.Errorf("file offset dont match"))
			return
		}
	}

	defer r.Body.Close()
	bytesWritten, err := io.Copy(file, r.Body)
	if err != nil {
		logErr(w, r, err)
		return
	}

	err = s.renameFile(file, r)
	if err != nil {
		logErr(w, r, err)
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

	fileName := strings.TrimPrefix(r.URL.Path, uploadRoute)
	fileName = filepath.Join(uploadDir, fileName)

	if _, err := os.Stat(fileName); err == nil {
		fileName = strings.TrimRight(file.Name(), incompleateExt)
	}

	return os.Rename(file.Name(), fileName)
}

// this cretes the file
func (s *Saver) postHandler(w http.ResponseWriter, r *http.Request) {
	fileName, err := s.getFilePath(r)
	if err != nil {
		logErr(w, r, err)
		return
	}

	_, err = os.Create(fileName)
	if err != nil {
		logErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// chekFile cheks the file exists or not
// and retuns a FileInfo || an err
func (s *Saver) openFile(baseUrl string, r *http.Request) (*os.File, error) {
	filePath, err := s.getFilePath(r)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, PermFile)
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
