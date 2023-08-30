package fileup

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGetData(t *testing.T) {
	mf := &mockFile{}
	up := mockUpper(mf)
	mc := newMockConn()
	err := up.getData(mc)
	if err != nil && !errors.Is(err, io.EOF) {
		t.Errorf("expected EOF err; got %v", err)
		t.FailNow()
	}
	t.Log(mf.buff.String())
}

type mockConn struct {
	reqCount, toalRes int
	results           []mockConnRes
}

type mockConnRes struct {
	msg  int
	data []byte
	err  error
}

var _ ConnReader = &mockConn{}

func (mc *mockConn) ReadMessage() (int, []byte, error) {
	if mc.reqCount < mc.toalRes {
		r := mc.results[mc.reqCount]
		mc.reqCount++ // next request
		return r.msg, r.data, r.err
	}

	return 0, nil, io.EOF
}
func (mc *mockConn) WriteJSON(v interface{}) error {
	fmt.Println(v)
	return nil
}

func (mc *mockConn) current() mockConnRes {
	if mc.reqCount < mc.toalRes {
		r := mc.results[mc.reqCount]
		return r
	}
	return mockConnRes{}
}

type mockFile struct {
	buff bytes.Buffer
}

var _ io.ReadWriteCloser = &mockFile{}

func (mf *mockFile) Read(p []byte) (n int, err error) {
	return mf.buff.Read(p)
}
func (mf *mockFile) Write(p []byte) (n int, err error) {
	return mf.buff.Write(p)
}
func (mf *mockFile) Close() error {
	return nil
}

func mockUpper(mf *mockFile) Upper {
	return Upper{
		RootDir:  "", // this wont be needed for testing
		BuffSize: BUFF_SIZE,
		wsUp: websocket.Upgrader{
			ReadBufferSize:  BUFF_SIZE,
			WriteBufferSize: BUFF_SIZE,
		},
		createFile: func(u *Upper) error {
			u.CurrentFile = mf
			u.CurrentFileName = "test"
			return nil
		},
		openFile: func(name string) (io.ReadCloser, error) {
			return mf, nil
		},
	}
}

func newMockConn() *mockConn {
	data := "some random data here it's just text file"
	sha := `{"checksum": "76cf9c0e4eacef622ddc73c583be62db9472500c65bf31d890cbff35a675b69a"}` // sha of the text
	mcon := []mockConnRes{
		{
			msg:  websocket.TextMessage,
			data: []byte("fine name"),
			err:  nil,
		},
		{
			msg:  websocket.BinaryMessage,
			data: []byte(data),
			err:  nil,
		},
		{
			msg:  websocket.TextMessage,
			data: []byte(sha),
			err:  nil,
		},
	}

	return &mockConn{
		reqCount: 0,
		toalRes:  len(mcon),
		results:  mcon,
	}
}

// func newMockFile() *mockFile {
// 	return &mockFile{}
// }
