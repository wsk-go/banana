package server

type Interceptor interface {
	Intercept(ctx Context) error
}

type ErrorHandler interface {
	Handle(ctx Context, err error) (any, error)
}

type ResponseHandler interface {
	Handle(ctx Context, body any) any
}