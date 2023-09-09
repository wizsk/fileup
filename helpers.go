package fileup

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func logErr(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Sorry somwting went wrong", http.StatusBadRequest)
	log.Printf("req form %q: err: %v\n", r.RemoteAddr, err)
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
