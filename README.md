<a name="sS9Oc"></a>
# 1.http包的使用
在go语言中，我们可以利用官方包http很快的完成一个Web后台的搭建。
```go
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

```
	这是一个非常基础的http包的使用，我们利用http.HandleFunc()将路径与对应的处理函数关联了起来，然后使用http.ListenAndServe()函数在 8001,端口启动了我们的服务。<br /> 	其中最关键的也就是最后的那一个启动函数，http.ListenAndServe()。
```go
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```
	参数：addr string （格式需要为 :端口的形式）<br />    handler Handler (Handler为一个接口类型)
```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
	需要实现ServerHTTP方法<br />参数：ResponseWriter,Request，这两个类型都是http定义好的可直接使用。

也就是，如果我们希望使用http.ListenAndServe()来监听端口启动服务，我们只需要传入一个实现了这个方法的结构体即可。
<a name="tDBtU"></a>
# 2.实现ServeHTTP
```go
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

```
	1. 定义了一个Engine结构体<br />2.Engine实现了ServeHTTP方法<br />3.在ServeHTTP中 通过switch 人为的指定了不同路劲的处理方法<br />4.main方法中 通过http.ListenAndServe 传入engine结构体 启动服务
<a name="ihHwX"></a>
# 3.Gweb Demo版本
通过上面两个例子，我们可以有一个大概的感觉。http帮助我们处理网络请求的流程<br />![image.png](https://cdn.nlark.com/yuque/0/2022/png/23003149/1664371924426-db91203a-795b-4fc8-b8c8-11f14eb9afdd.png#clientId=ubf0c0c61-16b9-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=531&id=u0c9a84fd&margin=%5Bobject%20Object%5D&name=image.png&originHeight=796&originWidth=1758&originalType=binary&ratio=1&rotation=0&showTitle=false&size=46700&status=done&style=none&taskId=u41529e96-4c1c-4f22-9d7b-c65744d11ea&title=&width=1172)<br />所有的请求都会被http.ListenAndServe拦截，然后根据不同的路径去匹配不同的处理函数，那么如何存储 这个路径与HandleFunc的关系表呢？<br />自然是使用我们Go语言中自带的map，来进行映射。<br />HTTP请求，在请求时的方式是不一样的。所以我们在存储路径时也需要保存其方法即 "方法-路径"的形式。这里我们就需要一个根据方法和路径 生成规定格式的路由映射表中的Key。<br />有了生成固定格式的key，我们在使用时肯定不希望需要手动的传入方法，而是通过直接调用对应的函数传入路径和处理函数（eg:Gin中的GET(path,handler)）。<br />最后启动http服务，模仿Gin。对http.ListenAndServe进行一层包裹。
```go
package Gweb

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
	engine.addRoute("PUT", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	handler, ok := engine.router[key]
	if ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND %q", req.URL.Path)
	}
}

func (engine *Engine) Run(port string) error {
	err := http.ListenAndServe(port, engine)
	return err
}

```
```go
type Engine struct {
	router map[string]HandlerFunc
}
```
	定义了Engine结构体，map类型的成员变量router，map的key为string类型即 路径，value为 HandlerFunc 类型 即处理函数的类型。
```go
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}
```
	New方法，返回Engine结构体指针，同时初始化了map
```go
func (engine Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}
```
	addRoute函数功能为：向router中添加key value。参数为 method 请求方法，pattern 请求路径，handler  处理函数。
```go
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}
```
	这里就是为了方便使用进一步分装了 addRoute方法，直接写死请求方法，用户使用时根据不同的请求方法调用对应函数 传入路径和处理函数即可。
```go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	handler, ok := engine.router[key]
	if ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND %q", req.URL.Path)
	}
}
```
实现ServeHTTP方法。（这里可以对所有请求进行一些通用的处理方法，后续的中间件即可放在这里）<br />根据我们规定好的格式 通过 method和path 凭借 router的key。然后去engine中的router寻找，如果找到即调用对应的处理函数。失败 报404
```go
func (engine *Engine) Run(port string) error {
	err := http.ListenAndServe(port, engine)
	return err
}
```
	Run方法，传入port 然后在函数当中调用 http.ListenAndServe()方法启动Web服务。



<a name="YWy0O"></a>

1.为什么需要context

先再次祭出http包处理网络请求的过程<br /><br />http.ListenAndServe()，第一个参数是端口 第二个参数是实现了ServeHTTP方法的结构体。<br />而在ServeHTTP中，参数要求为 http.ResponseWriter 和 *http.Request。<br />可以理解为，go将每个请求的response和request给我们传递了进来，我们后面的所有工作都是围绕response 和 request莱特完成。但是每次总是向下传递两个参数还是有些麻烦，而且其中一些信息我们希望先提取出来，而不是拿着request一个一个的拎出来，这时我们可以用context来封装一下，每个网络请求的基本信息，也就是上下文context。
<a name="QHF0u"></a>

2.context结构体中的内容

    type Context struct {
    	Writer     http.ResponseWriter
    	Req        *http.Request
    	Path       string
    	Method     string
    	StatusCode int
    }
    

    我们定义了一个Context结构体<br />成员为：<br />Writer：http.ResponseWriter类型，即网络请求的response<br />Req：*http.Requst类型 即网球请求<br />Path：string类型 网络请求的路径 初始化时可以从Req中获取<br />Method：string类型 网络请求的方法 初始化时可以从Req总获取<br />StatusCode：状态响应码<br />可以看出来，contxt结构体还是围绕http.ResponseWriter 和 http.Writer。使用contxt时为了提取里面的属性 从而更加方便我们后面的操作。

<a name="t27ht"></a>

3.功能方法

有时候不同功能之间可能会出现统一端代码，例如，每个请求响应都需要设置响应转台码。我们可以定义一个设置响应状态码的函数，这样在不同函数中就无需重复代码。
<a name="YIGGU"></a>

3.1 定义Header

    func (c *Context) SetHeader(key string, value string) {
    	c.Writer.Header().Set(key, value)
    }

    接收一对key value ，调用context中Writer的Head Set方法，设置请求头的Key value。

<a name="YYsXl"></a>

3. 2 定义状态码

   func (c *Context) Status(code int) {
   c.StatusCode = code
   c.Writer.WriteHeader(code)
   }


    Context结构体中StatusCode字段进行赋值。然后调用contxt中Writer的WriteHeader方法设置响应状态码。

<a name="wkhwt"></a>

3.3 获取form中key的value

    func (c *Context) PostForm(key string) string {
    	return c.Req.FormValue(key)
    }

    通过context中的Req中的FormValue(key) 获取form中对应key 的value值。

<a name="DgVek"></a>

3.4 获取query中key的value

    func (c *Context) Query(key string) string {
    	return c.Req.URL.Query().Get(key)
    }

    通过 context中的Req中的URL.Query.Get方法获取对应key的value

<a name="XZQJy"></a>

4. 二次封装

有了上面的基础方法我们就可以更好的封装我们所要进行的操作。作为一个Web后端框架，返回的数据Type应该是多种多样的，应该包含所有类型。这些类型都需要我们通过设置Header中的key value来完成配置。
<a name="msDbr"></a>

4.1 返回JSON结构体

    func (c *Context) JSON(code int, obj interface{}) {
    	c.SetHeader("Content-Type", "application/json")
    	c.Status(code)
    	encoder := json.NewEncoder(c.Writer)
    	err := encoder.Encode(obj)
    	if err != nil {
    		http.Error(c.Writer, err.Error(), 500)
    	}
    }

    接收两个参数：int类型的状态码 code，接口类型的 obj（也就是所有类型）。<br />首先我们需要通过SetHeader方法，在Header的设置响应数据格式。application/json。然后设置状态码。最后是对传入的obj进行json类型转换。转换失败则 通过http.Error方法报错。

<a name="spqJc"></a>

4.2 返回HTML

    func (c *Context) HTML(code int, html string) {
    	c.SetHeader("Content-Type", "text/html")
    	c.Status(code)
    	c.Writer.Write([]byte(html))
    }


<a name="q0cM6"></a>

4.3 返回String

    func (c *Context) String(code int, format string, values ...interface{}) {
    	c.SetHeader("Content-Type", "text/plain")
    	c.Status(code)
    	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
    }

<a name="uToUn"></a>

4.4 返回字节数据

    func (c *Context) Data(code int, data []byte) {
    	c.Status(code)
    	c.Writer.Write(data)
    }


<a name="GsFhz"></a>

5.Router

通过context结构体 我们实现了对每一个请求的细化管理。因此我们也不再使用之前的map来管理各个Route。我们模仿context的思想，创建一个router结构体，管理各个router。
<a name="UAwpG"></a>

5.1 struct router

    type router struct {
    	handlers map[string]HandlerFunc
    }

    定义一个router结构体 结构体成员为handlers map类型，key为string value为HandlerFunc类型的函数。

<a name="XTgPz"></a>

5.2 addRouter方法

    func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
    	log.Printf("Routes %4s - %s", method, pattern)
    	key := method + "-" + pattern
    	r.handlers[key] = handler
    }

<a name="q1O2E"></a>

5.3 handle方法

当路由的key被命中时调用的方法

    func (r *router) handle(c *Context) {
    	key := c.Method + "-" + c.Path
    	handler, ok := r.handlers[key]
    	if ok {
    		handler(c)
    	} else {
    		c.String(http.StatusNotFound, "404 not found %s \n", c.Path)
    	}
    }

    参数为我们定义的context

<a name="EvGKe"></a>

6.Gweb修改

<a name="mL4hf"></a>

6.1 Engine结构体

此处我们使用了router结构体来管理我们的router，所以Engine里面的router需要改变类型。从map类型改为*router

    type Engine struct {
    	router *router
    }

<a name="fdXnJ"></a>

6.2 addRoute

添加route前：

    func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    	key := method + "-" + pattern
    	engine.router[key] = handler
    }
    

    路由的添加 key 的组装在Gweb中完成<br />添加route后：

    func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    	engine.router.addRoute(method, pattern, handler)
    }

    交给了engine中的router完成，代码上更加简洁，同时各个成员负责的功能也更加明显。

<a name="I49w7"></a>

源代码

当前框架的源代码为三个go文件。

1. context.go 定义了每个请求的context上下文
2. Gweb.go 对外暴露提供方法的主文件
3. router.go 定义了router的处理方法

   package Gweb

   import (
   "encoding/json"
   "fmt"
   "net/http"
   )

   type H map[string]interface{}

   type Context struct {
   Writer     http.ResponseWriter
   Req        *http.Request
   Path       string
   Method     string
   StatusCode int
   }

   func newContext(w http.ResponseWriter, req *http.Request) *Context {
   return &Context{
   Writer: w,
   Req:    req,
   Path:   req.URL.Path,
   Method: req.Method,
   }
   }

   // PostForm
   // @Description: 获取form表单中对应key的value
   // @receiver c
   // @param key
   // @return string
   //
   func (c *Context) PostForm(key string) string {
   return c.Req.FormValue(key)
   }

   // Query
   // @Description: 获取URL中的query字段中key的value
   // @receiver c
   // @param key
   // @return string
   //
   func (c *Context) Query(key string) string {
   return c.Req.URL.Query().Get(key)
   }

   // Status
   // @Description: 设置响应状态码
   // @receiver c
   // @param code
   //
   func (c *Context) Status(code int) {
   c.StatusCode = code
   c.Writer.WriteHeader(code)
   }

   // SetHeader
   // @Description: 设置Header中的KeyValue
   // @receiver c
   // @param key
   // @param value
   //
   func (c *Context) SetHeader(key string, value string) {
   c.Writer.Header().Set(key, value)
   }

   func (c *Context) String(code int, format string, values ...interface{}) {
   c.SetHeader("Content-Type", "text/plain")
   c.Status(code)
   c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
   }

   func (c *Context) JSON(code int, obj interface{}) {
   c.SetHeader("Content-Type", "application/json")
   c.Status(code)
   encoder := json.NewEncoder(c.Writer)
   err := encoder.Encode(obj)
   if err != nil {
   http.Error(c.Writer, err.Error(), 500)
   }
   }

   func (c *Context) Data(code int, data []byte) {
   c.Status(code)
   c.Writer.Write(data)
   }

   func (c *Context) HTML(code int, html string) {
   c.SetHeader("Content-Type", "text/html")
   c.Status(code)
   c.Writer.Write([]byte(html))
   }


    package Gweb
    
    import (
    	"net/http"
    )
    
    type HandlerFunc func(*Context)
    
    type Engine struct {
    	router *router
    }
    
    func New() *Engine {
    	return &Engine{router: newRouter()}
    }
    
    func (engine Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    	engine.router.addRoute(method, pattern, handler)
    }
    
    func (engine *Engine) GET(pattern string, handler HandlerFunc) {
    	engine.addRoute("GET", pattern, handler)
    }
    
    func (engine *Engine) POST(pattern string, handler HandlerFunc) {
    	engine.addRoute("POST", pattern, handler)
    }
    
    func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
    	engine.addRoute("PUT", pattern, handler)
    }
    
    func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    	context := newContext(w, req)
    	engine.router.handle(context)
    }
    
    func (engine *Engine) Run(port string) error {
    	err := http.ListenAndServe(port, engine)
    	return err
    }
    

    package Gweb
    
    import (
    	"log"
    	"net/http"
    )
    
    type router struct {
    	handlers map[string]HandlerFunc
    }
    
    func newRouter() *router {
    	return &router{handlers: make(map[string]HandlerFunc)}
    }
    
    func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
    	log.Printf("Routes %4s - %s", method, pattern)
    	key := method + "-" + pattern
    	r.handlers[key] = handler
    }
    
    func (r *router) handle(c *Context) {
    	key := c.Method + "-" + c.Path
    	handler, ok := r.handlers[key]
    	if ok {
    		handler(c)
    	} else {
    		c.String(http.StatusNotFound, "404 not found %s \n", c.Path)
    	}
    }
    
