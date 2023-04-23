package server

import (
	"net/http"
)

type HandlerFunc func(ctx Context) (body any, err error)

type Mapping interface {
	GetMethod() string
	GetPath() string
	GetHandler() HandlerFunc
}

type RequestMapping struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func (th RequestMapping) GetMethod() string {
	return th.Method
}

func (th RequestMapping) GetPath() string {
	return th.Path
}

func (th RequestMapping) GetHandler() HandlerFunc {
	return th.Handler
}

type GetMapping struct {
	Path    string
	Handler HandlerFunc
}

func (th GetMapping) GetMethod() string {
	return http.MethodGet
}

func (th GetMapping) GetPath() string {
	return th.Path
}

func (th GetMapping) GetHandler() HandlerFunc {
	return th.Handler
}

type PostMapping struct {
	Path    string
	Handler HandlerFunc
}

func (th PostMapping) GetMethod() string {
	return http.MethodPost
}

func (th PostMapping) GetPath() string {
	return th.Path
}

func (th PostMapping) GetHandler() HandlerFunc {
	return th.Handler
}

type PutMapping struct {
	Path    string
	Handler HandlerFunc
}

func (th PutMapping) GetMethod() string {
	return http.MethodPut
}

func (th PutMapping) GetPath() string {
	return th.Path
}

func (th PutMapping) GetHandler() HandlerFunc {
	return th.Handler
}

type DeleteMapping struct {
	Path    string
	Handler HandlerFunc
}

func (th DeleteMapping) GetMethod() string {
	return http.MethodDelete
}

func (th DeleteMapping) GetPath() string {
	return th.Path
}

func (th DeleteMapping) GetHandler() HandlerFunc {
	return th.Handler
}
