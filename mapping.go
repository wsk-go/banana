package banana

import (
	"github.com/JackWSK/banana/defines"
	"net/http"
)

type Mapping interface {
	GetMethod() string
	GetPath() string
	GetHandler() defines.Handler
	GetRequiredQuery() []string
}

type RequestMapping struct {
	Method        string
	Path          string
	Handler       defines.Handler
	RequiredQuery []string
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

func (th RequestMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

type GetMapping struct {
	Path          string
	Handler       defines.Handler
	RequiredQuery []string
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

func (th GetMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

type PostMapping struct {
	Path          string
	Handler       defines.Handler
	RequiredQuery []string
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

func (th PostMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

type PutMapping struct {
	Path          string
	Handler       defines.Handler
	RequiredQuery []string
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

func (th PutMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

type DeleteMapping struct {
	Path          string
	Handler       defines.Handler
	RequiredQuery []string
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

func (th DeleteMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}
