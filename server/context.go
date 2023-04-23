package server

type Context interface {
	Next()

	GetBody() (value any, exists bool)
}
