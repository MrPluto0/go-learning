package main

import (
	"net/http"
)

func cal(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query().Get("num")
	w.Write([]byte(val + "xx"))
}

func main() {
	http.HandleFunc("/", cal)
	http.ListenAndServe(":8080", nil)
}
