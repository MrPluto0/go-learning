package main

import "function-hook/handler"

func main() {
	h := handler.LoginHandler{}
	h.Handler = &h
	h.Execute()
}
