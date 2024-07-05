package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	e := New()
	e.GET("/sb/ydx", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
	})

	e.POST("/pn/ydx", func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
		}
	})

	log.Fatal(e.Run(":9901"))
}
