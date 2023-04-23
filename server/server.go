package server

type Server interface {
	Handle(httpMethod, relativePath string, handlers HandlerFunc)

	Run(addr string) error
}
