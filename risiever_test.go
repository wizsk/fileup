package fileup

import (
	"io"
	"testing"
)

type mockConn struct {
}

var _ ConnReader = &mockConn{}

func (mc *mockConn) ReadMessage() (int, []byte, error) {
	return 0, nil, nil
}
func (mc *mockConn) WriteJSON(interface{}) error {
	return nil
}

type mockFile struct {
}

var _ io.ReadWriteCloser = &mockFile{}

func (mf *mockFile) Read(p []byte) (n int, err error) {
	return 0, nil
}
func (mf *mockFile) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (mf *mockFile) Close() error {
	return nil
}

func TestXX(t *testing.T) {
}
