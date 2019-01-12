package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/ttacon/glog"
)

// DefaultPort is the default port for the server
const DefaultPort = 2000

var port = flag.Int("port", DefaultPort, "port to run server")
var primaryURL = flag.String("primaryURL", "", "primary URL to proxy (its response will be returned to the client)")
var secondaryURL = flag.String("secondaryURL", "", "secondary URL to proxy (its response will be silenced)")

func main() {
	flag.Parse()

	router := mux.NewRouter()
	router.NotFoundHandler = proxy()
	fmt.Printf("Listening on :%d\n", *port)

	if *primaryURL == "" {
		glog.Fatal("No primary URL specified")
	} else {
		fmt.Printf("Primary: %s\n", *primaryURL)
	}

	if *secondaryURL == "" {
		glog.Fatal("No secondary URL specified")
	} else {
		fmt.Printf("Secondary: %s\n", *secondaryURL)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", *port), router)

}

func proxy() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request: %s\n", r.URL.Path)

		fakeRa := &http.Request{}
		copier.Copy(&fakeRa, &r)

		fakeRb := &http.Request{}
		copier.Copy(&fakeRb, &r)

		buf := new(bytes.Buffer)
		if r.Body != nil {
			buf.ReadFrom(r.Body)
		}

		fakeRa.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		fakeRb.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

		go (func() {
			fakeW := httptest.NewRecorder()
			reverseProxy(*secondaryURL, fakeW, fakeRa)
		})()

		reverseProxy(*primaryURL, w, fakeRb)
	})
}

func reverseProxy(target string, w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(target)
	glog.FatalIf(err)

	proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	proxy.ServeHTTP(w, r)
}
