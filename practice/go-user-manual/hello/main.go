package main

import (
	"fmt"
	"learning/practice/go-user-manual/greetings"
)

func main() {
	name := "Gypsophlia"
	message, err := greetings.Hello(name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(message)

	names := []string{
		"Gypsophlia",
		"Notseefire",
	}
	messages, err := greetings.Hellos(names)
	if err != nil {
		fmt.Println(err)
	}
	for _, message := range messages {
		fmt.Println(message)
	}
}
