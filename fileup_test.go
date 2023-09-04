package fileup

import (
	"encoding/json"
	"testing"

	"golang.org/x/net/websocket"
)

func TestHandeler(t *testing.T) {
	u := &Upper{}
	// starting real server :)
	go u.Serve()

	w := "ws://" + origin
	conn, err := websocket.Dial(w, "", origin)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer conn.Close()

	f := FileInfo{
		Name: "foo",
		Size: 12312313,
	}
	fj, err := json.Marshal(f)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = conn.Write(fj)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = conn.Write([]byte("skldflksdfjlksdfjlkdjfjkldsjflkdsjfkldsjflkdsjfkldsjfkldsjflkjskdfjklsdfjklsdfjklsdjklfdsklfdjskl"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if u.err != nil {
		t.Error(err)
		t.FailNow()
	}

	// for {
	// 	b := make([]byte, 10)
	// 	i, err := conn.Read(b)
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			t.Error(err)
	// 			t.FailNow()
	// 		} else {
	// 			t.Log(err)
	// 		}
	// 		return
	// 	}
	//
	// 	t.Logf("received %d bytes; data %q", i, string(b[:i]))
	// }
}
