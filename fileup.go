package fileup

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Saver struct {
	UpDir           string
	UpRoute         string
	IncomePleateExt string // default is ".part"
	Err             ErrHandler
}

const (
	permFile = 0664
	permDir  = 0755

	defaultExt = ".part"
)

func (s *Saver) Handeler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// post for creaate new file
		s.postHandler(w, r)

	case http.MethodPut:
		// patch for appending to the file
		s.putHandler(w, r)

	case http.MethodPatch:
		if err := s.renameFile(w, r); err != nil {
			s.Err(w, r, "file corropted", http.StatusBadRequest, err)
		}

	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
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

// patchHandeler receives data and appends it to the file
// it gets the file name form url and uuid
func (s *Saver) putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/offset+octet-stream" {
		s.Err(w, r,
			"bad request", http.StatusBadRequest,
			fmt.Errorf(`patchHanderler: expected "application/offset+octet-stream" got %q`, r.Header.Get("Content-Type")))
		return
	}

	file, err := s.openFile(r)
	if err != nil {
		s.Err(w, r, "", http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	if !s.checkUploadOffSet(w, r, file) {
		return
	}

	defer r.Body.Close()
	if _, err = io.Copy(file, r.Body); err != nil {
		s.Err(w, r, "", http.StatusInternalServerError, err)
		return
	}
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

// getFilePath generates filePaht
// eg. "upload/file.mp4-uuid.part"
func (s *Saver) getFilePath(r *http.Request) (string, error) {
	uuid := r.Header.Get("Upload-UUID")
	if uuid == "" {
		return "", errors.New("tempFileName: uuid is empty")
	}

	// name.ext-uuid.part
	fileName := fmt.Sprintf("%s-%s%s", strings.TrimPrefix(r.URL.Path, s.UpRoute), uuid, s.IncomePleateExt)

	return filepath.Join(s.UpDir, fileName), nil
}

// renameFile cheks if the file is fully written if true
// then it will try to rename the file without uuuid;
// if a file already exists with that given name then it will keep the uuid like "foo.mp4-uuid"
func (s *Saver) renameFile(w http.ResponseWriter, r *http.Request) error {
	// aka with ".part"
	partFileName, err := s.getFilePath(r)
	if err != nil {
		return err
	}
	upSize, err := uploadSize(r)
	if err != nil {
		return err
	}

	fileStat, err := os.Stat(partFileName)
	if err != nil {
		return err
	} else if fileStat.Size() != upSize {
		os.Remove(partFileName)
		return errors.New("renameFile: file size don't match")
	}

	sha256 := r.Header.Get("Sha256")
	if sha256 != "" {
		// todo add sha
	}

	uuid := r.Header.Get("Upload-UUID")
	// trying to remove uuid
	fileName := filepath.Join(s.UpDir, strings.TrimPrefix(r.URL.Path, s.UpRoute))
	// no err means file alreay exists
	if _, err = os.Stat(fileName); err == nil {
		fileName = fmt.Sprintf("%s-%s", fileName, uuid)
	}

	return os.Rename(partFileName, fileName)
}
