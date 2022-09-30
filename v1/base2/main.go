package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.paht=%q \n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q]=%q \n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND", req.URL)
	}
}

func main() {
	engine := new(Engine)
	err := http.ListenAndServe(":8001", engine)
	log.Fatalln(err.Error())
}
