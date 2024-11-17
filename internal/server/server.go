package server

import (
	"fmt"
	"net/http"
	"net/url"
)

func Run(host, port, endpoint, destination string) error {
	// Parse the destination URL
	targetURL, err := url.Parse(destination)
	if err != nil {
		return fmt.Errorf("invalid destination URL: %v", err)
	}

	// Create a new HTTP handler
	mux := http.NewServeMux()

	// Register health check endpoint
	mux.HandleFunc("/ping", ping)

	// Handle the root endpoint redirection to index.html
	if endpoint == "/" {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/index.html", http.StatusFound)
				return
			}
			ProxyRequestHandler(NewProxy(targetURL), targetURL)(w, r)
		})
	} else {
		mux.HandleFunc(endpoint, ProxyRequestHandler(NewProxy(targetURL), targetURL))
	}

	// Start HTTP server
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Starting HTTP server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("could not start HTTP server: %v", err)
	}
	return nil
}

func RunHTTPS(host, port, endpoint, destination, certFile, keyFile string) error {
	// Parse the destination URL
	targetURL, err := url.Parse(destination)
	if err != nil {
		return fmt.Errorf("invalid destination URL: %v", err)
	}

	// Create a new HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, ProxyRequestHandler(NewProxy(targetURL), targetURL))

	// Start HTTPS server
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Starting HTTPS server on %s\n", addr)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start the HTTPS server
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		return fmt.Errorf("could not start HTTPS server: %v", err)
	}
	return nil
}
