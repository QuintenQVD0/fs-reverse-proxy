package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reverse-proxy-learn/internal/server"
)

const version = "1.0.1"

func main() {
	// Command-line flags
	host := flag.String("host", "0.0.0.0", "Host to bind the server to")
	port := flag.String("listen-port", "8081", "Port to listen on")
	endpoint := flag.String("endpoint", "/", "Endpoint to proxy")
	destination := flag.String("destination-url", "http://localhost:8080", "Destination URL for the proxy")
	background := flag.Bool("background", false, "Run the server in the background")
	logFile := flag.String("log-file", "/home/container/farming-dashboard-reverse-server.log", "Log file path")
	httpsFlag := flag.Bool("https", false, "Enable HTTPS server")
	tlsCert := flag.String("tls-cert", "", "Path to TLS certificate file (required for HTTPS)")
	tlsKey := flag.String("tls-key", "", "Path to TLS key file (required for HTTPS)")

	flag.Parse()

	// Print version information
	fmt.Printf("Reverse Proxy Server Version: %s\n", version)

	// Validate destination URL
	if *destination == "" {
		log.Fatal("The --destination-url argument is required.")
	}

	// Validate TLS flags if HTTPS is enabled
	if *httpsFlag && (*tlsCert == "" || *tlsKey == "") {
		log.Fatal("Both --tls-cert and --tls-key must be provided for HTTPS.")
	}

	// Setup logging
	logOutput, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	defer logOutput.Close()

	if *background {
		// Fork the process to run in the background
		args := append([]string{}, os.Args[1:]...)
		for i, arg := range args {
			if arg == "--background=true" || arg == "--background" {
				args[i] = "--background=false"
			}
		}

		cmd := exec.Command(os.Args[0], args...)
		cmd.Stdout = logOutput
		cmd.Stderr = logOutput

		if err := cmd.Start(); err != nil {
			log.Fatalf("Failed to start the reverse proxy server in the background: %v", err)
		}
		log.Printf("The reverse proxy server is running in the background. Logs are being written to %s\n", *logFile)
		os.Exit(0)
	} else {
		// Combine log output for foreground process
		multiWriter := io.MultiWriter(os.Stdout, logOutput)
		log.SetOutput(multiWriter)

		// Log startup with version
		log.Printf("Starting Reverse Proxy Server (Version: %s)\n", version)

		// Start the server in the foreground
		if *httpsFlag {
			// Run HTTPS
			if err := server.RunHTTPS(*host, *port, *endpoint, *destination, *tlsCert, *tlsKey); err != nil {
				log.Fatalf("Could not start HTTPS server: %v", err)
			}
		} else {
			// Run HTTP
			if err := server.RunHTTP(*host, *port, *endpoint, *destination); err != nil {
				log.Fatalf("Could not start HTTP server: %v", err)
			}
		}
	}
}
