package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func RunHTTP(host, port, endpoint, destination string) error {
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
	log.Printf("Starting HTTP server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("could not start HTTP server: %v", err)
	}
	return nil
}

// RunHTTPS starts an HTTPS server and handles HTTP-to-HTTPS redirection.
func RunHTTPS(host, port, endpoint, destination, certFile, keyFile string) error {
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
			// Redirect to HTTPS if the scheme is HTTP
			if r.TLS == nil {
				redirectToHTTPS(w, r)
				return
			}

			if r.URL.Path == "/" {
				http.Redirect(w, r, "/index.html", http.StatusFound)
				return
			}
			ProxyRequestHandler(NewProxy(targetURL), targetURL)(w, r)
		})
	} else {
		mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			// Redirect to HTTPS if the scheme is HTTP
			if r.TLS == nil {
				redirectToHTTPS(w, r)
				return
			}
			ProxyRequestHandler(NewProxy(targetURL), targetURL)(w, r)
		})
	}

	// Start HTTPS server
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting HTTPS server on %s\n", addr)
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

// redirectToHTTPS sends a 301 redirect to the HTTPS version of the request.
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	httpsURL := fmt.Sprintf("https://%s%s", r.Host, r.URL.RequestURI())
	http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
}
