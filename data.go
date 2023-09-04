package fileup

import "io"

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	path string
	w    io.WriteCloser
}
