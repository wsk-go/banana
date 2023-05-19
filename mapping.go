package banana

import (
	"net/http"
)

type Mapping interface {
	GetMethod() string
	GetPath() string
	GetHandler() Handler
	GetRequiredQuery() []string
	GetUserInfo(key string) (any, bool)
}

type RequestMapping struct {
	Method        string
	Path          string
	Handler       Handler
	RequiredQuery []string
	UserInfo      map[string]any
}

func (th RequestMapping) GetMethod() string {
	return th.Method
}

func (th RequestMapping) GetPath() string {
	return th.Path
}

func (th RequestMapping) GetHandler() Handler {
	return th.Handler
}

func (th RequestMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

func (th RequestMapping) GetUserInfo(key string) (any, bool) {
	if th.UserInfo == nil {
		return nil, false
	}
	v, ok := th.UserInfo[key]
	return v, ok
}

type GetMapping struct {
	Path          string
	Handler       Handler
	RequiredQuery []string
	UserInfo      map[string]any
}

func (th GetMapping) GetMethod() string {
	return http.MethodGet
}

func (th GetMapping) GetPath() string {
	return th.Path
}

func (th GetMapping) GetHandler() Handler {
	return th.Handler
}

func (th GetMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

func (th GetMapping) GetUserInfo(key string) (any, bool) {
	if th.UserInfo == nil {
		return nil, false
	}
	v, ok := th.UserInfo[key]
	return v, ok
}

type PostMapping struct {
	Path          string
	Handler       Handler
	RequiredQuery []string
	UserInfo      map[string]any
}

func (th PostMapping) GetMethod() string {
	return http.MethodPost
}

func (th PostMapping) GetPath() string {
	return th.Path
}

func (th PostMapping) GetHandler() Handler {
	return th.Handler
}

func (th PostMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

func (th PostMapping) GetUserInfo(key string) (any, bool) {
	if th.UserInfo == nil {
		return nil, false
	}
	v, ok := th.UserInfo[key]
	return v, ok
}

type PutMapping struct {
	Path          string
	Handler       Handler
	RequiredQuery []string
	UserInfo      map[string]any
}

func (th PutMapping) GetMethod() string {
	return http.MethodPut
}

func (th PutMapping) GetPath() string {
	return th.Path
}

func (th PutMapping) GetHandler() Handler {
	return th.Handler
}

func (th PutMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

func (th PutMapping) GetUserInfo(key string) (any, bool) {
	if th.UserInfo == nil {
		return nil, false
	}
	v, ok := th.UserInfo[key]
	return v, ok
}

type DeleteMapping struct {
	Path          string
	Handler       Handler
	RequiredQuery []string
	UserInfo      map[string]any
}

func (th DeleteMapping) GetMethod() string {
	return http.MethodDelete
}

func (th DeleteMapping) GetPath() string {
	return th.Path
}

func (th DeleteMapping) GetHandler() Handler {
	return th.Handler
}

func (th DeleteMapping) GetRequiredQuery() []string {
	return th.RequiredQuery
}

func (th DeleteMapping) GetUserInfo(key string) (any, bool) {
	if th.UserInfo == nil {
		return nil, false
	}
	v, ok := th.UserInfo[key]
	return v, ok
}
