package test

import (
	"net/http"
	"testing"
)

func BenchmarkHttp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http.Get("http://localhost:8080?num=123")
	}
}
