package banana

import (
	"github.com/JackWSK/banana/defines"
	"net/http"
)

type Mapping interface {
	GetMethod() string
	GetPath() string
	GetHandler() defines.Handler
}

type RequestMapping struct {
	Method  string
	Path    string
	Handler defines.Handler
}

func (th RequestMapping) GetMethod() string {
	return th.Method
}

func (th RequestMapping) GetPath() string {
	return th.Path
}

func (th RequestMapping) GetHandler() defines.Handler {
	return th.Handler
}

type GetMapping struct {
	Path    string
	Handler defines.Handler
}

func (th GetMapping) GetMethod() string {
	return http.MethodGet
}

func (th GetMapping) GetPath() string {
	return th.Path
}

func (th GetMapping) GetHandler() defines.Handler {
	return th.Handler
}

type PostMapping struct {
	Path    string
	Handler defines.Handler
}

func (th PostMapping) GetMethod() string {
	return http.MethodPost
}

func (th PostMapping) GetPath() string {
	return th.Path
}

func (th PostMapping) GetHandler() defines.Handler {
	return th.Handler
}

type PutMapping struct {
	Path    string
	Handler defines.Handler
}

func (th PutMapping) GetMethod() string {
	return http.MethodPut
}

func (th PutMapping) GetPath() string {
	return th.Path
}

func (th PutMapping) GetHandler() defines.Handler {
	return th.Handler
}

type DeleteMapping struct {
	Path    string
	Handler defines.Handler
}

func (th DeleteMapping) GetMethod() string {
	return http.MethodDelete
}

func (th DeleteMapping) GetPath() string {
	return th.Path
}

func (th DeleteMapping) GetHandler() defines.Handler {
	return th.Handler
}
