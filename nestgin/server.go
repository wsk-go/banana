package nestgin

import (
	"fmt"
	"github.com/JackWSK/go-nest/server"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	engine *gin.Engine
}

func New() *GinServer {
	engine := gin.Default()

	// handler
	engine.Use(func(ctx *gin.Context) {
		ctx.Next()
		if err, ok := ctx.Get("__err__"); ok {
			ctx.AbortWithError(500, fmt.Errorf("%+v", err))
		} else {
			if body, ok := ctx.Get("__body__"); ok {
				ctx.JSON(200, body)
			}
		}
	})
	return &GinServer{
		engine: gin.Default(),
	}
}

func (th *GinServer) Handle(httpMethod, relativePath string, handler server.HandlerFunc) {
	th.engine.Handle(httpMethod, relativePath, func(ctx *gin.Context) {
		body, err := handler(ctx)
		ctx.Set("__err__", err)
		ctx.Set("__body__", body)
	})
}

func (th *GinServer) Run(addr string) error {
	return th.engine.Run(addr)
}
