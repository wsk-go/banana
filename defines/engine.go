package defines

import (
	"context"
	"mime/multipart"
)

type Handler = func(Context) error

type Engine interface {
	Add(method, path string, handlers Handler)

	Listen(addr string) error
}

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
}
