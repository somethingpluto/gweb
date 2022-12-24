package Gweb

import "fmt"

type Router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

func (r *Router) addRoute(method string, path string, handler HandlerFunc) {
	key := method + "-" + path
	r.handlers[key] = handler
}

func (r Router) handle(ctx *context) {
	key := ctx.Method + "-" + ctx.Path
	handler, ok := r.handlers[key]
	if ok {
		handler(ctx)
	} else {
		fmt.Fprintf(ctx.Writer, "404页面未找打")
	}
}
