package nestgin

import (
	"fmt"
	"github.com/JackWSK/go-nest/server"
	"github.com/gin-gonic/gin"
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
	engine        *gin.Engine
	errorHandlers []server.ExceptionHandler
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

func defaultHandler(ctx *gin.Context) {
	ctx.Next()
	if err, ok := ctx.Get(errorKey); ok && err != nil {
		ctx.String(500, fmt.Sprintf("%+v", err))
		return
	} else if body, ok := ctx.Get(bodyKey); ok {
		ctx.JSON(200, body)
		return
	}
	_ = ctx.AbortWithError(500, fmt.Errorf("something went wrong"))
}

func (th *GinServer) handleException() {
	th.engine.Use(func(ctx *gin.Context) {
		c := ctx.MustGet(contextKey).(*context)
		defer func() {
			err := recover()
			if err != nil {
				if len(th.errorHandlers) > 0 {
					for _, handler := range th.errorHandlers {
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
			if len(th.errorHandlers) > 0 {
				c.SetError(nil)
				for _, handler := range th.errorHandlers {
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
