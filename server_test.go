package fileup

import (
	"sync"
	"testing"
	"time"

	// "github.com/shirou/gopsutil/cpu"
	"golang.org/x/net/websocket"
)

func TestServer(t *testing.T) {
	// go cpuUsages(t)
	s := NewServer("tmp")
	go s.Serve("/ws") // starting real server :)
	time.Sleep(50 * time.Millisecond)
	t.Log("TestServer server started at", s.origin)

	w := "ws://" + s.origin + "/ws"

	var wg sync.WaitGroup
	for i := 0; i <= 20; i++ {
		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()

			conn, err := websocket.Dial(w, "", s.origin)
			Nil(t, err)
			time.Sleep(1 * time.Second)
			_ = conn
		}(t)
	}
	wg.Wait()

	time.Sleep(20 * time.Second)
}

// func cpuUsages(t *testing.T) {
// 	for {
// 		// Get CPU usage percentages
// 		percentages, err := cpu.Percent(time.Second, false)
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			return
// 		}

// 		// Print CPU usage for each core
// 		for i, pct := range percentages {
// 			t.Logf("CPU %d Usage: %.2f%%\n", i, pct)
// 		}

// 		// Sleep for a while before getting CPU usage again
// 		time.Sleep(time.Second)
// 	}
// }
