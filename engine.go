package banana

import (
	"time"
)

// Cookie data for c.Cookie
type Cookie struct {
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Path        string    `json:"path"`
	Domain      string    `json:"domain"`
	MaxAge      int       `json:"max_age"`
	Expires     time.Time `json:"expires"`
	Secure      bool      `json:"secure"`
	HTTPOnly    bool      `json:"http_only"`
	SameSite    string    `json:"same_site"`
	SessionOnly bool      `json:"session_only"`
}

type Handler func(Context) error

type EngineMiddlewareFunc func(ctx Context) error

type Engine interface {
	// Add allows you to specify a HTTP method to register a route
	Add(method, path string, handlers Handler)
	// Listen serves HTTP requests from the given addr.
	Listen(addr string) error
	// Use middleware function
	Use(middlewares ...EngineMiddlewareFunc)
}
