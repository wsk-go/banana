package fiberengine

import (
	"github.com/gofiber/fiber/v2"
	_recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wsk-go/banana"
	"mime/multipart"
	"time"
)

type FiberEngine struct {
	app *fiber.App
}

const contextKey = "__FiberEngineContextKey__"

func New(config ...fiber.Config) *FiberEngine {
	app := fiber.New(config...)
	app.Use(_recover.New())
	app.Use(func(ctx *fiber.Ctx) error {
		defer func() {
			ctx.Context().RemoveUserValue(contextKey)
		}()
		ctx.Context().SetUserValue(contextKey, &Context{
			ctx: ctx,
		})
		return ctx.Next()
	})
	return &FiberEngine{
		app: app,
	}
}

func (f *FiberEngine) App() *fiber.App {
	return f.app
}

func (f *FiberEngine) Listen(addr string) error {
	return f.app.Listen(addr)
}

func (f *FiberEngine) Use(middlewares ...banana.EngineMiddlewareFunc) {
	for _, middleware := range middlewares {
		f.app.Use(func(ctx *fiber.Ctx) error {
			c := ctx.Context().Value(contextKey).(*Context)
			return middleware(c)
		})
	}
}

func (f *FiberEngine) Add(method, path string, handler banana.Handler) {
	f.app.Add(method, path, func(ctx *fiber.Ctx) error {
		c := ctx.Context().Value(contextKey).(*Context)
		return handler(c)
	})
}

type Context struct {
	ctx *fiber.Ctx
}

func (c *Context) Query(key string, defaultValue ...string) string {
	return c.ctx.Query(key, defaultValue...)
}

func (c *Context) QueryInt(key string, defaultValue ...int) int {
	return c.ctx.QueryInt(key, defaultValue...)
}

func (c *Context) QueryBool(key string, defaultValue ...bool) bool {
	return c.ctx.QueryBool(key, defaultValue...)
}

func (c *Context) QueryFloat(key string, defaultValue ...float64) float64 {
	return c.ctx.QueryFloat(key, defaultValue...)
}

func (c *Context) QueryParser(out any) error {
	return c.ctx.QueryParser(out)
}

func (c *Context) Params(key string, defaultValue ...string) string {
	return c.ctx.Params(key, defaultValue...)
}

func (c *Context) AllParams() map[string]string {
	return c.ctx.AllParams()
}

func (c *Context) ParamsParser(out interface{}) error {
	return c.ctx.ParamsParser(out)
}

func (c *Context) ParamsInt(key string, defaultValue ...int) (int, error) {
	return c.ctx.ParamsInt(key, defaultValue...)
}

func (c *Context) Protocol() string {
	return c.ctx.Protocol()
}

func (c *Context) Path() string {
	return c.ctx.Path()
}

func (c *Context) OriginalURL() string {
	return c.ctx.OriginalURL()
}

func (c *Context) Get(key string, defaultValue ...string) string {
	return c.ctx.Get(key, defaultValue...)
}

func (c *Context) GetRespHeader(key string, defaultValue ...string) string {
	return c.ctx.GetRespHeader(key, defaultValue...)
}

func (c *Context) GetRespHeaders() map[string]string {
	return c.ctx.GetRespHeaders()
}

func (c *Context) GetReqHeaders() map[string]string {
	return c.ctx.GetReqHeaders()
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(key)
}

func (c *Context) FormValue(key string, defaultValue ...string) string {
	return c.ctx.FormValue(key, defaultValue...)
}

func (c *Context) Download(file string, filename ...string) error {
	return c.ctx.Download(file, filename...)
}

func (c *Context) Cookies(key string, defaultValue ...string) string {
	return c.ctx.Cookies(key, defaultValue...)
}

func (c *Context) Cookie(cookie *banana.Cookie) {
	c.ctx.Cookie(&fiber.Cookie{
		Name:        cookie.Name,
		Value:       cookie.Value,
		Path:        cookie.Path,
		Domain:      cookie.Domain,
		MaxAge:      cookie.MaxAge,
		Expires:     cookie.Expires,
		Secure:      cookie.Secure,
		HTTPOnly:    cookie.HTTPOnly,
		SameSite:    cookie.SameSite,
		SessionOnly: cookie.SessionOnly,
	})
}

func (c *Context) Status(status int) {
	c.ctx.Status(status)
}

// Write appends p into response body.
func (c *Context) Write(p []byte) (int, error) {
	return c.ctx.Write(p)
}

// Writef appends f & a into response body writer.
func (c *Context) Writef(f string, a ...interface{}) (int, error) {
	return c.ctx.Writef(f, a...)
}

// WriteString appends s to response body.
func (c *Context) WriteString(s string) (int, error) {
	return c.ctx.WriteString(s)
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

func (c *Context) RawContext() any {
	return c.ctx
}
