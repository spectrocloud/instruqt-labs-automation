package main

import (
	"log"
	"net/http"
	"os"
)

// Webserver starts an HTTP server to serve static files
func Webserver() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Use assets directory for local testing, /app for container
	dir := os.Getenv("SERVE_DIR")
	if dir == "" {
		dir = "assets"
	}

	fs := http.FileServer(http.Dir(dir))

	http.Handle("/", fs)

	log.Printf("Starting server on port %s\n", port)
	log.Printf("Serving files from %s directory\n", dir)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
