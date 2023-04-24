package nestgin

import (
	"fmt"
	"github.com/JackWSK/go-nest/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

const errorKey = "__err__"
const bodyKey = "__body__"
const contextKey = "__contextKey__"

type context struct {
	c *gin.Context
}

func (th *context) Next() {
	th.c.Next()
}

func (th *context) JSON(code int, obj any) {
	th.c.JSON(code, obj)
}

func (th *context) GetBody() any {
	if body, ok := th.c.Get(bodyKey); ok {
		return body
	}
	return nil
}

func (th *context) SetBody(body any) {
	th.c.Set(bodyKey, body)
}

func (th *context) SetError(err error) {
	th.c.Set(errorKey, err)
}

func (th *context) GetError() error {
	if err, ok := th.c.Get(errorKey); ok {
		return err.(error)
	}
	return nil
}

type GinServer struct {
	engine            *gin.Engine
	exceptionHandlers []server.ExceptionHandler
}

func New() *GinServer {
	engine := gin.Default()

	engine.Use(func(ctx *gin.Context) {
		c := &context{c: ctx}
		defer func() {
			c.c = nil
		}()
		ctx.Set(contextKey, c)
		ctx.Next()
		fmt.Println("xxx")
	})

	// middleware

	server := &GinServer{
		engine: engine,
	}
	server.handleException()

	return server
}

func (th *GinServer) handleWrite() {
	th.engine.Use(func(ctx *gin.Context) {
		c := ctx.MustGet(contextKey).(*context)
		ctx.Next()
		if err := c.GetError(); err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		} else {
			//ctx.Data(http.StatusOK)
		}
	})
}

func (th *GinServer) handleException() {
	th.engine.Use(func(ctx *gin.Context) {
		c := ctx.MustGet(contextKey).(*context)
		defer func() {
			err := recover()
			if err != nil {
				if len(th.exceptionHandlers) > 0 {
					for _, handler := range th.exceptionHandlers {
						body := handler.HandlePanic(c, err)
						c.SetError(nil)
						c.SetBody(body)
					}
				} else {
					panic(err)
				}
			}
		}()
		ctx.Next()
		if err := c.GetError(); err != nil {
			if len(th.exceptionHandlers) > 0 {
				c.SetError(nil)
				for _, handler := range th.exceptionHandlers {
					body := handler.HandleError(c, err)
					c.SetBody(body)
				}
			}
			return
		}
	})
}

func (th *GinServer) Handle(httpMethod, relativePath string, handler server.HandlerFunc) {
	th.engine.Handle(httpMethod, relativePath, func(ctx *gin.Context) {
		c := ctx.MustGet(contextKey).(*context)
		body, err := handler(c)
		c.SetError(err)
		c.SetBody(body)
	})
}

func (th *GinServer) UseInterceptors(interceptors ...server.Interceptor) {
	for _, interceptor := range interceptors {
		th.engine.Use(func(ctx *gin.Context) {
			c := ctx.MustGet(contextKey).(*context)
			err := interceptor.Handle(c)
			if err != nil {
				c.SetError(err)
			}
		})
	}
}

func (th *GinServer) Run(addr string) error {
	return th.engine.Run(addr)
}
