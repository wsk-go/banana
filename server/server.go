package server

type Server interface {
	UseInterceptors(interceptors ...Interceptor)

	Handle(httpMethod, relativePath string, handlers HandlerFunc)

	Run(addr string) error
}
