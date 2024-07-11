package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	e := New()
	//e.Use(func(c *Context) {
	//	// Start timer
	//	t := time.Now()
	//	log.Println("start: ", t.String())
	//	// Process request
	//	c.Next()
	//	// Calculate resolution time
	//	log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
	//})

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

	ydx := e.Group("/ydx")
	{
		ydx.Use(func(c *Context) {
			log.Printf("\n/ydx\n")
		})

		ydx.GET("/sb", func(c *Context) {
			c.JSON(http.StatusOK, H{
				"sb": "ydx",
			})
		})
	}

	log.Fatal(e.Run(":9901"))
}
