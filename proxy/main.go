package main

import (
	"log"
	"net/http"
	"proxy/https"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello world\n"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello world\n"))
	})
	log.Println("listening on :8443")
	log.Println(https.ListenAndServer(":8443", "test.crt", "test.key", mux))
}
