package main

import (
	"log"
	"net/http"
	"net/rpc"
)

type Cal struct{}

func (c *Cal) Square(num int, result *int) error {
	*result = num * num
	return nil
}

func main() {
	rpc.Register(new(Cal))
	rpc.HandleHTTP()

	log.Println("starting rpc service")

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("start rpc error")
	}
}
