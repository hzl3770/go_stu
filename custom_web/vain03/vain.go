package main

import (
	"fmt"
	"net/http"
)

/*
Engine掌管整个生命周期
Run启动
*/
var (
	default404Body = []byte("404 page not found")
)

type Engine struct {
	DefaultRouter
}

func New() *Engine {
	return &Engine{
		DefaultRouter: newDefaultRouter(),
	}
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("url: %#v", r.URL)
	k := e.makeRouteKey(r.Method, r.URL.Path)
	if handler, ok := e.m[k]; ok {
		handler(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write(default404Body)
	if err != nil {
		fmt.Errorf("failed to write to client, err: %v", err)
	}
}
