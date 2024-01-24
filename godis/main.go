package main

import (
	"example.com/redis/lib/logger"
	"example.com/redis/tcp"
)

func main() {
	handler := tcp.NewEchoHandler()

	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address:    ":6666",
		MaxConnect: 10,
		Timeout:    1000,
	}, handler)

	if err != nil {
		logger.Fatal(err)
	}
}
