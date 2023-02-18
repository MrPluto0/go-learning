package main

import (
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var num = 12
	var res int
	if err := client.Call("Cal.Square", 12, &res); err != nil {
		log.Fatal("Failed to call Cal.Square: ", err)
	}
	log.Printf("%d^2 = %d", num, res)
}
