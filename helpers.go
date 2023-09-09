package fileup

import (
	"fmt"
	"net/http"
	"strconv"
)

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
