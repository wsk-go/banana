package server

type Interceptor interface {
	Intercept(ctx Context) error
}
