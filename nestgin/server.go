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

func (th *context) GetBody() (value any, exists bool) {
	return th.c.Get(bodyKey)
}

func (th *context) SetBody(body any) {
	th.c.Set(bodyKey, body)
}

func (th *context) setError(err error) {
	th.c.Set(errorKey, err)
}

type GinServer struct {
	engine        *gin.Engine
	errorHandlers []server.ErrorHandler
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
	server.handleResponseAndError()

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

func (th *GinServer) handleResponseAndError() {
	th.engine.Use(func(ctx *gin.Context) {
		ctx.Next()
		c := ctx.MustGet(contextKey).(*context)
		if err, ok := ctx.MustGet(errorKey).(error); ok && err != nil {
			for _, handler := range th.errorHandlers {
				body, err := handler.Handle(c, err)
				c.setError(err)
				c.SetBody(body)
			}
			return
		} else if body, ok := ctx.Get(bodyKey); ok {
			ctx.JSON(200, body)
			return
		}
		_ = ctx.AbortWithError(500, fmt.Errorf("something went wrong"))
	})
}

func (th *GinServer) Handle(httpMethod, relativePath string, handler server.HandlerFunc) {
	th.engine.Handle(httpMethod, relativePath, func(ctx *gin.Context) {
		c := ctx.MustGet(contextKey).(*context)
		body, err := handler(c)
		c.setError(err)
		c.SetBody(body)
	})
}

func (th *GinServer) UseInterceptors(interceptors ...server.Interceptor) {
	for _, interceptor := range interceptors {
		th.engine.Use(func(ctx *gin.Context) {
			err := interceptor.Intercept(ctx.MustGet(contextKey).(*context))
			if err != nil {
				ctx.Set(errorKey, err)
				ctx.Abort()
			}
		})
	}
}

func (th *GinServer) Run(addr string) error {
	return th.engine.Run(addr)
}
