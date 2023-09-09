package fileup

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func (s *Saver) checkUploadOffSet(w http.ResponseWriter, r *http.Request, f *os.File) bool {
	uploadOffset, err := uploadOffset(r)
	if err != nil {
		s.Err(w, r, "", http.StatusBadRequest, err)
		return false
	}

	if fileStat, err := f.Stat(); err != nil {
		s.Err(w, r, "", http.StatusInternalServerError, err)
		return false
	} else if fileStat.Size() != uploadOffset {
		s.Err(w, r, "file offset size dosen't match", http.StatusBadRequest, err)
		return false
	}

	return true
}

func uploadOffset(r *http.Request) (int64, error) {
	uploadOffset, err := strconv.ParseInt(r.Header.Get("Upload-Offset"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid upload offset: %v", err)
	}
	return uploadOffset, nil
}

func uploadSize(r *http.Request) (int64, error) {
	uploadOffset, err := strconv.ParseInt(r.Header.Get("Upload-Size"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid upload offset: %v", err)
	}
	return uploadOffset, nil
}
