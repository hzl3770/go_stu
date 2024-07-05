package main

import (
	"fmt"
	"log"
)

func main() {
	e := New()
	e.GET("/sb/ydx", func(ctx *Context) {
		fmt.Fprintf(ctx.Writer, "URL.Path = %q\n", ctx.Path)

	})

	e.POST("/pn/ydx", func(ctx *Context) {
		for k, v := range ctx.Request.Header {
			fmt.Fprintf(ctx.Writer, "Header[%q] = %q\n", k, v)
		}
	})

	log.Fatal(e.Run(":9901"))
}
