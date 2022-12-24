package Gweb

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(ctx *context)

// 定义Engine结构体
type Engine struct {
	router *Router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := newContext(w, req)
	e.router.handle(context)
}

func (e *Engine) addRoute(method string, path string, handler HandlerFunc) {
	showTheRoute(method, path)
	e.router.addRoute(method, path, handler)
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute("GET", path, handler)
}

func (e *Engine) Run(port int) error {
	sPort := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(sPort, e)
	return err
}

func showTheRoute(method string, path string) {
	fmt.Printf("\tmethod:%s \t\t path:%s \t\t \n", method, path)
}
