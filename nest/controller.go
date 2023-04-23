package nest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerFunc func(ctx *gin.Context) (body any, err error)

type Mapping interface {
	getMethod() string
	getPath() string
	getHandler() string
}

type RequestMapping struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func (th RequestMapping) getMethod() string {
	return th.Method
}

func (th RequestMapping) getPath() string {
	return th.Path
}

func (th RequestMapping) getHandler() HandlerFunc {
	return th.Handler
}

type GetMapping struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func (th GetMapping) getMethod() string {
	return http.MethodGet
}

func (th GetMapping) getPath() string {
	return th.Path
}

func (th GetMapping) getHandler() HandlerFunc {
	return th.Handler
}

type PostMapping struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func (th PostMapping) getMethod() string {
	return http.MethodPost
}

func (th PostMapping) getPath() string {
	return th.Path
}

func (th PostMapping) getHandler() HandlerFunc {
	return th.Handler
}

type PutMapping struct {
	Method  string
	Path    string
	Handler func(ctx *gin.Context) (body any, err error)
}

func (th PutMapping) getMethod() string {
	return http.MethodPut
}

func (th PutMapping) getPath() string {
	return th.Path
}

func (th PutMapping) getHandler() HandlerFunc {
	return th.Handler
}

type DeleteMapping struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func (th DeleteMapping) getMethod() string {
	return http.MethodDelete
}

func (th DeleteMapping) getPath() string {
	return th.Path
}

func (th DeleteMapping) getHandler() HandlerFunc {
	return th.Handler
}
