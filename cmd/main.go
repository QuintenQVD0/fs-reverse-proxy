package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"reverse-proxy-learn/internal/server"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Host to bind the server to")
	port := flag.String("listen-port", "8080", "Port to listen on")
	endpoint := flag.String("endpoint", "/", "Endpoint to proxy")
	destination := flag.String("destination-url", "http://localhost:8080", "Destination URL for the proxy")
	background := flag.Bool("background", false, "Run the server in the background")
	logFile := flag.String("log-file", "farming-dashboard-reverse-server.log", "Log file path if running in the background")
	httpsFlag := flag.Bool("https", false, "Enable HTTPS server")
	tlsCert := flag.String("tls-cert", "", "Path to TLS certificate file (required for HTTPS)")
	tlsKey := flag.String("tls-key", "", "Path to TLS key file (required for HTTPS)")

	flag.Parse()

	if *destination == "" {
		log.Fatal("The --destination-url argument is required.")
	}

	if *httpsFlag && (*tlsCert == "" || *tlsKey == "") {
		log.Fatal("Both --tls-cert and --tls-key must be provided for HTTPS.")
	}

	if *background {
		// Fork the process to run in the background
		args := append([]string{}, os.Args[1:]...)
		for i, arg := range args {
			if arg == "--background=true" || arg == "--background" {
				args[i] = "--background=false"
			}
		}

		cmd := exec.Command(os.Args[0], args...)
		logOutput, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Could not open log file: %v", err)
		}
		defer logOutput.Close()
		cmd.Stdout = logOutput
		cmd.Stderr = logOutput

		if err := cmd.Start(); err != nil {
			log.Fatalf("Failed to start the server in the background: %v", err)
		}
		log.Printf("Server is running in the background. Logs are being written to %s\n", *logFile)
		os.Exit(0)
	}

	// Start the server in the foreground
	if *httpsFlag {
		// Run HTTPS
		if err := server.RunHTTPS(*host, *port, *endpoint, *destination, *tlsCert, *tlsKey); err != nil {
			log.Fatalf("Could not start HTTPS server: %v", err)
		}
	} else {
		// Run HTTP
		if err := server.Run(*host, *port, *endpoint, *destination); err != nil {
			log.Fatalf("Could not start HTTP server: %v", err)
		}
	}
}
