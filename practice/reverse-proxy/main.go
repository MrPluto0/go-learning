// Reverse Proxy：From 127.0.0.1:8082 To https://www.baidu.com
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	director := func(req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = targetHost
		req.Host = targetHost
	}
	proxy := &httputil.ReverseProxy{Director: director}

	return proxy, nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewProxy("www.baidu.com")
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8082", nil))
}
