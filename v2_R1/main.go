package main

import (
	"Gweb/v1_R1/Gweb"
	"fmt"
	"net/http"
)

func main() {
	engine := Gweb.New()
	engine.GET("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello World")
	})
	err := engine.Run(8080)
	if err != nil {
		panic(err)
	}
}
