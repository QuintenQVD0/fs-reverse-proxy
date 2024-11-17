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

	flag.Parse()

	if *destination == "" {
		log.Fatal("The --destination-url argument is required.")
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

	// Set up logging if running in the foreground
	if *logFile != "" {
		logOutput, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Could not open log file: %v", err)
		}
		defer logOutput.Close()
		log.SetOutput(logOutput)
	}

	if err := server.Run(*host, *port, *endpoint, *destination); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}
