package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	server, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		panic(err)
	}
	log.Printf("Listen to 127.0.0.1:8081...")

	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("Accept failed %v", err)
			continue
		}
		go process(client)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Read error %v", err)
			break
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Printf("Write error %v", err)
			break
		}
	}
}
