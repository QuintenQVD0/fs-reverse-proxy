package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(target)
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy, target *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ PROXY SERVER ] Request received at %s at %s\n", r.URL, time.Now().UTC())

		// Update the request to forward it to the target server
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Host = target.Host

		// Proxy the request
		log.Printf("[ PROXY SERVER ] Proxying request to %s at %s\n", r.URL, time.Now().UTC())
		proxy.ServeHTTP(w, r)
	}
}
