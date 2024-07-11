package main

import (
	"net/http"
)

var (
	default404Body = []byte("404 page not found")
)

type HandlerFunc func(*Context)

type HandleChain []HandlerFunc

type router struct {
	handleMap map[string]HandleChain
}

func newDefaultRouter() router {
	return router{
		handleMap: make(map[string]HandleChain),
	}
}

func (r *router) addRoute(method string, pattern string, handlerFunc ...HandlerFunc) {
	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handlerFunc == nil {
		panic("http: nil handler")
	}
	key := r.makeRouteKey(method, pattern)
	if _, exist := r.handleMap[key]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	r.handleMap[key] = handlerFunc
}

func (r *router) makeRouteKey(method string, pattern string) string {
	return method + "-" + pattern
}

func (r *router) route(c *Context) {
	k := r.makeRouteKey(c.Method, c.Path)
	//log.Printf("\nk %v\n", k)

	if handlers, ok := r.handleMap[k]; ok {
		c.handles = handlers
	} else {
		c.handles = append(c.handles, func(ctx *Context) {
			ctx.Data(http.StatusNotFound, default404Body)
		})
	}
	//log.Printf("\nctx: %v", c)
	c.Next()
}
