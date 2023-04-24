package nestgin

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handle httprouter.Handle

type Engine struct {
	router *httprouter.Router
}

func NewEngine() *Engine {
	return &Engine{
		router: httprouter.New(),
	}
}

func (th *Engine) Get(path string, handle Handle) {
	th.Handle(http.MethodGet, path, handle)
}

// GET is a shortcut for router.Handle(http.MethodGet, path, handle)
func (th *Engine) GET(path string, handle Handle) {
	th.Handle(http.MethodGet, path, handle)
}

// HEAD is a shortcut for router.Handle(http.MethodHead, path, handle)
func (th *Engine) HEAD(path string, handle Handle) {
	th.Handle(http.MethodHead, path, handle)
}

// OPTIONS is a shortcut for router.Handle(http.MethodOptions, path, handle)
func (th *Engine) OPTIONS(path string, handle Handle) {
	th.Handle(http.MethodOptions, path, handle)
}

// POST is a shortcut for router.Handle(http.MethodPost, path, handle)
func (th *Engine) POST(path string, handle Handle) {
	th.Handle(http.MethodPost, path, handle)
}

// PUT is a shortcut for router.Handle(http.MethodPut, path, handle)
func (th *Engine) PUT(path string, handle Handle) {
	th.Handle(http.MethodPut, path, handle)
}

// PATCH is a shortcut for router.Handle(http.MethodPatch, path, handle)
func (th *Engine) PATCH(path string, handle Handle) {
	th.Handle(http.MethodPatch, path, handle)
}

// DELETE is a shortcut for router.Handle(http.MethodDelete, path, handle)
func (th *Engine) DELETE(path string, handle Handle) {
	th.Handle(http.MethodDelete, path, handle)
}

func (th *Engine) Handle(method, path string, handle Handle) {
	th.router.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handle(writer, request, params)
	})
}

func (th *Engine) Run() error {
	return http.ListenAndServe(":12345", th)
}

// Implement the ServeHTTP method on our new type
func (th *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.router.ServeHTTP(w, r)
}
