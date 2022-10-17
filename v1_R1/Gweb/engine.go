package Gweb

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

// 定义Engine结构体
type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 1.根据req中内容生成 方法—路径
	key := req.Method + "-" + req.URL.Path
	fmt.Println(key)
	// 2.从router映射表中查找对应的处理函数
	handlerFunc, ok := e.router[key]
	if ok {
		handlerFunc(w, req)
	} else { // 3.404 未匹配到的情况
		fmt.Fprintf(w, "404 页面未找到")
	}
}

func (e *Engine) addRoute(method string, path string, handler HandlerFunc) {
	showTheRoute(method, path)
	key := fmt.Sprintf("%s-%s", method, path)
	e.router[key] = handler
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute("GET", path, handler)
}

func (e *Engine) printAllRoute() {
	for key, value := range e.router {
		fmt.Printf("URL:  %s    func:  %v", key, value)
	}
}

func (e *Engine) Run(port int) error {
	sPort := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(sPort, e)
	return err
}

func showTheRoute(method string, path string) {
	fmt.Printf("\tmethod:%s \t\t path:%s \t\t \n", method, path)
}
