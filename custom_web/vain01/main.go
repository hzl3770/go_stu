package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/sb/ydx", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
	})

	http.HandleFunc("/pn/ydx", func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
		}
	})

	log.Fatal(http.ListenAndServe(":9901", nil))
}
