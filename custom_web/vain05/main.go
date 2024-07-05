package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	e := New()
	e.GET("/sb/ydx", func(ctx *Context) {
		fmt.Println("\nname: ", ctx.GetQuery("name"))
		ctx.String(http.StatusOK, "sbydx : %v, %v", ctx.Path, ctx.GetQuery("name"))
	})

	e.POST("/pn/ydx", func(ctx *Context) {
		h := H{}
		for k, v := range ctx.Request.Header {
			h[k] = v
		}
		ctx.JSON(http.StatusOK, h)
	})

	log.Fatal(e.Run(":9901"))
}
