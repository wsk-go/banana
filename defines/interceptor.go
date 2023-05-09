package defines

type Interceptor interface {
	Pre(ctx Context) error
	After(ctx Context) error
	Completion(ctx Context)
}
