package main

import (
	"log"
	"net/http"
)

func Webserver() {

	port := ":80"

	fs := http.FileServer(http.Dir("/app"))

	http.Handle("/", fs)

	log.Printf("Starting server on port %s\n", port)
	log.Println("Serving files from the current directory.")
	log.Println("Access the site at: :80")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
