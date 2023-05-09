package defines

type Middleware interface {
	Handle(ctx Context) error
}

type MiddlewareFunc func(ctx Context) error
