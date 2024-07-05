package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

/*
Context
1.管理请求和响应
请求信息：context 通常包含有关请求的信息，如请求路径、HTTP方法、请求头、查询参数和请求体等。
响应信息：context 也用于设置和管理响应信息，如响应状态码、响应头和响应体。
*/
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	//req info
	Method string
	Path   string

	//resp info
	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	c := new(Context)
	c.Writer = w
	c.Request = r
	c.Method = r.Method
	c.Path = r.URL.Path

	return c
}

// head
func (c *Context) requestHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) GetHeader(key string) string {
	return c.requestHeader(key)
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) ContentType() string {
	return c.requestHeader("Content-Type")
}

// req param

func (c *Context) GetQuery(key string) string {
	if c.Request == nil {
		return ""
	}

	return c.Request.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

var (
	jsonContentType = []string{"application/json; charset=utf-8"}

	plainContentType = []string{"text/plain; charset=utf-8"}
)

func (c *Context) JSON(code int, obj any) {
	c.writeContentType(jsonContentType)
	c.Status(code)

	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(code int, format string, values ...any) {
	c.writeContentType(plainContentType)
	c.Status(code)

	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) writeContentType(value []string) {
	header := c.Writer.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
