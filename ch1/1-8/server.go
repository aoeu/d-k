package main

// This code was snarfed from an example in `go doc http.ListenAndServe`.

import (
	"io"
	"log"
	"net/http"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}