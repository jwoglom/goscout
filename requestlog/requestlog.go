package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// DefaultPort is the default port for the server
const DefaultPort = 4000

var port = flag.Int("port", DefaultPort, "port to run server")

func main() {
	flag.Parse()

	router := mux.NewRouter()
	router.NotFoundHandler = log()
	fmt.Printf("Listening on :%d\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), router)

}

func log() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request: %s\n", r.URL)
		fmt.Printf("type: %s\n", r.Method)
		for k, v := range r.Form {
			fmt.Printf("form: %s = %s\n", k, v)
		}
		for k, v := range r.Header {
			fmt.Printf("header: %s = %s\n", k, v)
		}
		fmt.Printf("\n")
		if r.Body != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			fmt.Printf("body: %s\n", buf)
		}
		fmt.Printf("\n\n")

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok: port %d\n", *port)
	})
}
