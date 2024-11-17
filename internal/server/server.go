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

	// Create a new router
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

	// Start the server
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Starting server at %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}

	return nil
}
