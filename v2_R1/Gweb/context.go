package Gweb

import "net/http"

type context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	statusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *context {
	return &context{Writer: w, Req: req, Path: req.URL.Path, Method: req.Method}
}
