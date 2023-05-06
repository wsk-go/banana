package fiberengine

import (
	"github.com/JackWSK/banana/defines"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"time"
)

type FiberEngine struct {
	app *fiber.App
}

type Context struct {
	ctx *fiber.Ctx
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Context().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Context().Done()
}

func (c *Context) Err() error {
	return c.ctx.Context().Err()
}

func (c *Context) Value(key any) any {
	return c.ctx.Context().Value(key)
}

// SetValue is used to store a new key/value pair exclusively for this context.
func (c *Context) SetValue(key any, value any) {
	c.ctx.Context().SetUserValue(key, value)
}

// RemoveValue removes the given key and the value under it in ctx.
func (c *Context) RemoveValue(key any) {
	c.ctx.Context().RemoveUserValue(key)
}

// Header It writes a header in the response.
func (c *Context) Header(key any, value any) {
	c.ctx.Context().SetUserValue(key, value)
}

func (c *Context) BodyParser(out any) error {
	return c.ctx.BodyParser(out)
}

func (c *Context) Append(field string, values ...string) {
	c.ctx.Append(field, values...)
}

func (c *Context) Attachment(filename ...string) {
	c.ctx.Attachment(filename...)
}

func (c *Context) BaseURL() string {
	return c.ctx.BaseURL()
}

func (c *Context) JSON(data any) error {
	return c.ctx.JSON(data)
}

func (c *Context) JSONP(data any, callback ...string) error {
	return c.ctx.JSONP(data, callback...)
}

func (c *Context) XML(data interface{}) error {
	return c.ctx.XML(data)
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.ctx.MultipartForm()
}

func (c *Context) Next() error {
	return c.ctx.Next()
}

func (c *Context) Body() []byte {
	return c.ctx.Body()
}

func New(config ...fiber.Config) *FiberEngine {
	app := fiber.New(config...)
	app.Use(func(ctx *fiber.Ctx) error {
		//ctx.Context().SetUserValue("__context__", &Context{
		//	ctx: ctx,
		//})
		return ctx.Next()
	})
	return &FiberEngine{
		app: fiber.New(config...),
	}
}

func (f *FiberEngine) App() *fiber.App {
	return f.app
}

func (f *FiberEngine) Run(addr string) error {
	return f.app.Listen(addr)
}

func (f *FiberEngine) Add(method, path string, handler defines.Handler) {
	f.app.Add(method, path, func(ctx *fiber.Ctx) error {
		//c := ctx.Context().Value("__context__").(*Context)
		return handler(&Context{
			ctx: ctx,
		})
	})
}
