package server

type Context interface {
	Next()

	GetBody() any

	SetBody(body any)
}
