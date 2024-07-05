package main

import (
	"fmt"
	"net/http"
)

type DefaultRouter struct {
	m map[string]http.HandlerFunc
}

func newDefaultRouter() DefaultRouter {
	return DefaultRouter{
		m: make(map[string]http.HandlerFunc),
	}
}

func (d *DefaultRouter) addRoute(method string, pattern string, handlerFunc http.HandlerFunc) {
	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handlerFunc == nil {
		panic("http: nil handler")
	}
	key := d.makeRouteKey(method, pattern)
	if _, exist := d.m[key]; exist {
		panic("http: multiple registrations for " + pattern)
	}
	fmt.Printf("%#v", d)
	d.m[key] = handlerFunc
}

func (d *DefaultRouter) makeRouteKey(method string, pattern string) string {
	return method + "-" + pattern
}

func (d *DefaultRouter) GET(pattern string, handlerFunc http.HandlerFunc) {
	d.addRoute(http.MethodGet, pattern, handlerFunc)
}

func (d *DefaultRouter) POST(pattern string, handlerFunc http.HandlerFunc) {
	d.addRoute(http.MethodPost, pattern, handlerFunc)
}
