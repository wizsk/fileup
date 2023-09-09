package fileup

import (
	"log"
	"net/http"
)

type ErrHandler func(w http.ResponseWriter, r *http.Request, httpMsg string, httpCode int, err error)

var defaultErrHandeer ErrHandler = func(w http.ResponseWriter, r *http.Request, httpMsg string, httpCode int, err error) {
	if httpMsg == "" {
		httpMsg = "sorry someting went wrong"
	}

	http.Error(w, httpMsg, httpCode)
	log.Printf("[Error] Request from %q: err: %v\n", r.RemoteAddr, err)
}

// upRoute is the "example/route/path"
func NewSaver(upRoute, upDir string, errHandeler ErrHandler) *Saver {
	if errHandeler == nil {
		errHandeler = defaultErrHandeer
	}

	return &Saver{
		UpRoute:         upRoute,
		UpDir:           upDir,
		IncomePleateExt: defaultExt,
		Err:             errHandeler,
	}
}
