package main

import "learning/practice/function-hook/handler"

func main() {
	h := handler.LoginHandler{}
	h.Handler = &h
	h.Execute()
}
