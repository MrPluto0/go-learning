package tcp

import (
	"bufio"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	closeChan := make(chan bool)
	if err != nil {
		t.Error(err)
		return
	}
	addr := listener.Addr().String()
	go ListenAndServe(listener, NewEchoHandler(), closeChan)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 10; i++ {
		content := strconv.Itoa(rand.Int())
		_, err := conn.Write([]byte(content + "\n"))
		if err != nil {
			t.Error(err)
			return
		}
		bufReader := bufio.NewReader(conn)
		line, _, err := bufReader.ReadLine()
		if err != nil {
			t.Error(err)
			return
		}
		if string(line) != content {
			t.Error("unequal request and response: ", content, string(line))
			return
		}
	}
	_ = conn.Close()

	// test perfect close
	for i := 0; i < 5; i++ {
		_, _ = net.Dial("tcp", addr)
	}
	closeChan <- true
	time.Sleep(time.Second)
}
