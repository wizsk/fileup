package fileup

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	path string
}

type Sha256 struct {
	Sum string `json:"sum"`
}


const (
	msgTypeSha256 string = "sha256"
	// msgType string = ""

)

type StatusMsg struct {
	Type  string `json:"type"` // 
	Body  string `json:"body"`
	Error bool   `json:"error"`
}
