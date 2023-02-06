package main

import (
	"encoding/json"
	"fmt"
)

type studentSon struct {
	SonName string `json:"sonName"`
	SonAge  int    `json:"sonAge"`
}

type student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	studentSon
}

func main() {
	var student = student{
		Name:       "test",
		studentSon: studentSon{SonName: "123"},
	}
	res, _ := json.Marshal(student)

	// check out the difference
	fmt.Println(string(res))
	fmt.Printf("%+v", student)
}
