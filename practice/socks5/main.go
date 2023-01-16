package main

import (
	"bufio"
	"fmt"
	"io"
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
	err := auth(reader, conn)
	if err != nil {
		log.Printf("client %v auth failed:%v\n", conn.RemoteAddr(), err)
	}
	log.Panicln("auth success")
}

func auth(reader *bufio.Reader, conn net.Conn) error {
	const sock5Ver = 0x05

	ver, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read version failed: %w", err)
	}
	if ver != sock5Ver {
		return fmt.Errorf("not supported version: %v", ver)
	}

	methodSize, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read method size failed: %w", err)
	}
	method := make([]byte, methodSize)
	_, err = io.ReadFull(reader, method)
	if err != nil {
		return fmt.Errorf("read method body failed: %w", err)
	}

	log.Println("version", ver, "method", method)

	_, err = conn.Write([]byte{sock5Ver, 0x00})
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}
