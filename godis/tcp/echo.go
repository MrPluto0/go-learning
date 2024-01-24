package tcp

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	wait "example.com/redis/lib/sync"
)

type EchoClient struct {
	conn    net.Conn
	waiting wait.Wait
}

func (c *EchoClient) Close() error {
	c.waiting.WaitWithTimeout(10 * time.Second)
	c.conn.Close()
	return nil
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Bool
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Load() {
		conn.Close()
		return
	}

	client := &EchoClient{
		conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("connection info")
				h.activeConn.Delete(client)
			} else {
				log.Println(err)
			}
			return
		}
		// use lock to confirm the msg written to connection
		client.waiting.Add(1)
		conn.Write([]byte(msg))
		client.waiting.Done()
	}
}

func (h *EchoHandler) Close() error {
	log.Println("handler shutting down...")
	h.closing.Store(true)
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		client.Close()
		return true
	})
	return nil
}
