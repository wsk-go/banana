package server

type Interceptor interface {
	Handle(ctx Context) error
}

type ExceptionHandler interface {
	// HandleError Handle convert err to body, return nil means give up handle
	HandleError(ctx Context, err error) (body any)

	// HandlePanic Handle convert err from panic to body, return nil means give up handle
	HandlePanic(ctx Context, err any) (body any)
}

type ResponseHandler interface {
	Handle(ctx Context, body any)
}
