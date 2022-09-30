package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "url.paht = %q \n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q]=%q \n", k, v)
	}
}
