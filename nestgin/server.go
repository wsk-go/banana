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

type GinServer struct {
	engine *gin.Engine
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
	engine.Use(func(ctx *gin.Context) {
		ctx.Next()
		if err, ok := ctx.Get(errorKey); ok && err != nil {
			ctx.String(500, fmt.Sprintf("%+v", err))
			return
		} else if body, ok := ctx.Get(bodyKey); ok {
			ctx.JSON(200, body)
			return
		}
		_ = ctx.AbortWithError(500, fmt.Errorf("something went wrong"))
	})
	return &GinServer{
		engine: engine,
	}
}

func (th *GinServer) Handle(httpMethod, relativePath string, handler server.HandlerFunc) {
	th.engine.Handle(httpMethod, relativePath, func(ctx *gin.Context) {
		body, err := handler(ctx.MustGet(contextKey).(*context))
		ctx.Set(errorKey, err)
		ctx.Set(bodyKey, body)
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
