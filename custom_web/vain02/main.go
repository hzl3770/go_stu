package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {
}

/*
 http.DefaultServeMux 里面实现了对路由-》handle的封装

*/

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/sb/ydx":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	case "/pn/ydx":
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}

}

func main() {
	e := new(Engine)
	log.Fatal(http.ListenAndServe(":9901", e))
}
