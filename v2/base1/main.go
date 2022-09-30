package main

import (
	"Gweb/v2/base1/Gweb"
	"net/http"
)

func main() {
	engine := Gweb.New()
	engine.GET("/", func(context *Gweb.Context) {
		context.String(http.StatusOK, "Hello World")
	})
	err := engine.Run(":8001")
	if err != nil {
		panic(err)
	}
}
