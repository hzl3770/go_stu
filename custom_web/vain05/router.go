package main

import (
	"fmt"
	"net/http"
)

var (
	default404Body = []byte("404 page not found")
)

type HandlerFunc func(*Context)

type router struct {
	middleware []HandlerFunc
	handles    map[string]HandlerFunc
}

func newDefaultRouter() router {
	return router{
		handles: make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method string, pattern string, handlerFunc HandlerFunc) {
	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handlerFunc == nil {
		panic("http: nil handler")
	}
	key := r.makeRouteKey(method, pattern)
	if _, exist := r.handles[key]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	r.handles[key] = handlerFunc
}

func (r *router) makeRouteKey(method string, pattern string) string {
	return method + "-" + pattern
}

func (r *router) route(c *Context) {
	k := r.makeRouteKey(c.Method, c.Path)
	fmt.Printf("k %v", k)
	if handler, ok := r.handles[k]; ok {
		handler(c)
		return
	}

	c.Data(http.StatusNotFound, default404Body)
}

func (r *router) Use(middleware ...HandlerFunc) {
	r.middleware = append(r.middleware, middleware...)
}
