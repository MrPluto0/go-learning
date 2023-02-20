package client_test

import (
	"log"
	"net/rpc"
	"testing"
)

func BenchmarkRPC(b *testing.B) {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var num = 12
	var res int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Call("Cal.Square", num, &res)
	}
}

func BenchmarkCal(b *testing.B) {
	var res int
	var num = 12

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = num * num
	}

	b.Log(res)
}
