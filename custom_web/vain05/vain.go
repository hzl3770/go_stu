package main

import (
	"fmt"
	"net/http"
)

/*
Engine掌管整个生命周期
Run启动
*/
type Engine struct {
	*RouterGroup

	router router

	groups []*RouterGroup
}

func New() *Engine {
	e := &Engine{
		router: newDefaultRouter(),
	}

	e.RouterGroup = &RouterGroup{
		engine:   e,
		basePath: "/",
	}

	return e
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r)
	ctx.Writer = w
	ctx.Request = r
	fmt.Printf("url: %#v", r.URL)
	fmt.Printf("engine: %+v", e)
	e.router.route(ctx)
}
