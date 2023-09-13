package fileup

import (
	"fmt"
	"net/http"
	"os"
)

type ErrHandler func(w http.ResponseWriter, r *http.Request, httpMsg string, httpCode int, err error)

var defaultErrHandeer ErrHandler = func(w http.ResponseWriter, r *http.Request, httpMsg string, httpCode int, err error) {
	if httpMsg == "" {
		httpMsg = "sorry someting went wrong"
	}

	http.Error(w, httpMsg, httpCode)
	fmt.Printf("[Error] Request from %q: err: %v\n", r.RemoteAddr, err)
}

// upRoute is the "example/route/path"
func NewSaver(upRoute, upDir string) *Saver {
	return &Saver{
		UpRoute:         upRoute,
		UpDir:           upDir,
		IncomePleateExt: defaultExt,
		Err:             defaultErrHandeer,
	}
}

// NewSaverMkdir creates the dir and Retuns a Server
func NewSaverMkdir(upRoute, upDir string) (*Saver, error) {
	 err := os.Mkdir(upDir, os.ModePerm)
	if !os.IsExist(err) {
		return nil, err
	}

	return &Saver{
		UpRoute:         upRoute,
		UpDir:           upDir,
		IncomePleateExt: defaultExt,
		Err:             defaultErrHandeer,
	}, nil

}
