package banana

import (
	"context"
	"mime/multipart"
)

type Context interface {
	context.Context
	// BodyParser binds the request body to a struct.
	// It supports decoding the following content types based on the Content-Type header:
	// application/json, application/xml, application/x-www-form-urlencoded, multipart/form-data
	// If none of the content types above are matched, it will return a ErrUnprocessableEntity error
	BodyParser(out any) error

	// Header It writes a header in the response.
	Header(key any, value any)

	// SetValue is used to store a new key/value pair exclusively for this context.
	SetValue(key any, value any)

	// RemoveValue removes the given key and the value under it in ctx.
	RemoveValue(key any)

	// Append the specified value to the HTTP response header field.
	// If the header is not already set, it creates the header with the specified value.
	Append(field string, values ...string)

	// Attachment sets the HTTP response Content-Disposition header field to attachment.
	Attachment(filename ...string)

	// BaseURL returns (protocol + host + base path).
	BaseURL() string

	// Query returns the query string parameter in the url.
	// Defaults to empty string "" if the query doesn't exist.
	// If a default value is given, it will return that value if the query doesn't exist.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	Query(key string, defaultValue ...string) string

	// QueryInt returns integer value of key string parameter in the url.
	// Default to empty or invalid key is 0.
	//
	//	GET /?name=alex&wanna_cake=2&id=
	//	QueryInt("wanna_cake", 1) == 2
	//	QueryInt("name", 1) == 1
	//	QueryInt("id", 1) == 1
	//	QueryInt("id") == 0
	QueryInt(key string, defaultValue ...int) int

	// QueryBool returns bool value of key string parameter in the url.
	// Default to empty or invalid key is true.
	//
	//	Get /?name=alex&want_pizza=false&id=
	//	QueryBool("want_pizza") == false
	//	QueryBool("want_pizza", true) == false
	//	QueryBool("name") == false
	//	QueryBool("name", true) == true
	//	QueryBool("id") == false
	//	QueryBool("id", true) == true
	QueryBool(key string, defaultValue ...bool) bool

	// QueryFloat returns float64 value of key string parameter in the url.
	// Default to empty or invalid key is 0.
	//
	//	GET /?name=alex&amount=32.23&id=
	//	QueryFloat("amount") = 32.23
	//	QueryFloat("amount", 3) = 32.23
	//	QueryFloat("name", 1) = 1
	//	QueryFloat("name") = 0
	//	QueryFloat("id", 3) = 3
	QueryFloat(key string, defaultValue ...float64) float64

	// QueryParser binds the query string to a struct.
	QueryParser(out any) error

	// Params is used to get the route parameters.
	// Defaults to empty string "" if the param doesn't exist.
	// If a default value is given, it will return that value if the param doesn't exist.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	Params(key string, defaultValue ...string) string

	// AllParams Params is used to get all route parameters.
	// Using Params method to get params.
	AllParams() map[string]string

	// ParamsParser binds the param string to a struct.
	ParamsParser(out interface{}) error

	// ParamsInt is used to get an integer from the route parameters
	// it defaults to zero if the parameter is not found or if the
	// parameter cannot be converted to an integer
	// If a default value is given, it will return that value in case the param
	// doesn't exist or cannot be converted to an integer
	ParamsInt(key string, defaultValue ...int) (int, error)

	// JSON converts any interface or string to JSON.
	// Array and slice values encode as JSON arrays,
	// except that []byte encodes as a base64-encoded string,
	// and a nil slice encodes as the null JSON value.
	// This method also sets the content header to application/json.
	JSON(data any) error

	// JSONP sends a JSON response with JSONP support.
	// This method is identical to JSON, except that it opts-in to JSONP callback support.
	// By default, the callback name is simply callback.
	JSONP(data any, callback ...string) error

	// XML converts any interface or string to XML.
	// This method also sets the content header to application/xml.
	XML(data interface{}) error

	// MultipartForm returns request's multipart form.
	//
	// Returns ErrNoMultipartForm if request's content-type
	// isn't 'multipart/form-data'.
	//
	// All uploaded temporary files are automatically deleted after
	// returning from RequestHandler. Either move or copy uploaded files
	// into new place if you want retaining them.
	//
	// Use SaveMultipartFile function for permanently saving uploaded file.
	//
	// The returned form is valid until your request handler returns.
	//
	// See also FormFile and FormValue.
	MultipartForm() (*multipart.Form, error)

	// Next executes the next method in the stack that matches the current route.
	Next() error

	// Body contains the raw body submitted in a POST request.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	Body() []byte

	// Status sets the HTTP status for the response.
	// This method is chainable.
	Status(status int)

	// Write appends p into response body.
	Write(p []byte) (int, error)

	// Writef appends f & a into response body writer.
	Writef(f string, a ...interface{}) (int, error)

	// WriteString appends s to response body.
	WriteString(s string) (int, error)

	// RawContext return the context your framework use
	RawContext() any

	// Protocol contains the request protocol string: http or https for TLS requests.
	// Please use Config.EnableTrustedProxyCheck to prevent header spoofing, in case when your app is behind the proxy.
	Protocol() string

	// Path returns the path part of the request URL.
	// Optionally, you could override the path.
	Path() string

	// OriginalURL contains the original request URL.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	OriginalURL() string

	// Get returns the HTTP request header specified by field.
	// Field names are case-insensitive
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	Get(key string, defaultValue ...string) string

	// GetRespHeader returns the HTTP response header specified by field.
	// Field names are case-insensitive
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	GetRespHeader(key string, defaultValue ...string) string

	// GetRespHeaders returns the HTTP response headers.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	GetRespHeaders() map[string]string

	// GetReqHeaders returns the HTTP request headers.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	GetReqHeaders() map[string]string

	// FormFile returns the first file by key from a MultipartForm.
	FormFile(key string) (*multipart.FileHeader, error)

	// FormValue returns the first value by key from a MultipartForm.
	// Search is performed in QueryArgs, PostArgs, MultipartForm and FormFile in this particular order.
	// Defaults to the empty string "" if the form value doesn't exist.
	// If a default value is given, it will return that value if the form value does not exist.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	FormValue(key string, defaultValue ...string) string

	// Download transfers the file from path as an attachment.
	// Typically, browsers will prompt the user for download.
	// By default, the Content-Disposition header filename= parameter is the filepath (this typically appears in the browser dialog).
	// Override this default with the filename parameter.
	Download(file string, filename ...string) error

	// Cookies are used for getting a cookie value by key.
	// Defaults to the empty string "" if the cookie doesn't exist.
	// If a default value is given, it will return that value if the cookie doesn't exist.
	// The returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	Cookies(key string, defaultValue ...string) string

	// Cookie sets a cookie by passing a cookie struct.
	Cookie(cookie *Cookie)

	// RoutePath Path you write on route
	RoutePath() string
	// RouteMethod Method you write on route
	RouteMethod() string
	// RouteParams Params you write on route
	RouteParams() []string
}

type ValidateFunc func(obj any) error

func BodyParser[T any](ctx Context, fs ...ValidateFunc) (*T, error) {
	var t T
	err := ctx.BodyParser(&t)
	if err != nil {
		return nil, err
	}

	err = validate(&t, fs)
	return &t, err
}

func QueryParser[T any](ctx Context, fs ...ValidateFunc) (*T, error) {
	var t T
	err := ctx.QueryParser(&t)
	if err != nil {
		return nil, err
	}

	err = validate(&t, fs)
	return &t, err
}

func ParamParser[T any](ctx Context, fs ...ValidateFunc) (*T, error) {
	var t T
	err := ctx.ParamsParser(&t)

	if err != nil {
		return nil, err
	}

	err = validate(&t, fs)
	return &t, err
}

func validate[T any](t *T, fs []ValidateFunc) (err error) {
	if len(fs) > 0 {
		for _, f := range fs {
			if err = f(t); err != nil {
				return err
			}
		}
	}

	return nil
}
