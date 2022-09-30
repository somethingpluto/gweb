package main

import (
	"Gweb/v1/base3/Gweb"
	"fmt"
	"net/http"
)

func main() {
	engine := Gweb.New()
	engine.GET("/", indexHandler)
	engine.GET("/hello", helloHandler)
	err := engine.Run(":8001")
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
