package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	dir := filepath.Dir(path)

	if err := os.MkdirAll(fmt.Sprintf("extracted/%v", dir), 0755); err != nil {
		log.Printf("Error creating directory %v: %v", dir, err)
		return
	}

	// Create a new file named "example.txt"
	file, err := os.Create(fmt.Sprintf("extracted/%v", path))
	if err != nil {
		log.Printf("Error Creating file %v: %v", path, err)
		return
	}
	defer file.Close()

	all, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading data %v: %v", path, err)
		return
	}

	_, err = file.Write(all)
	if err != nil {
		log.Printf("Error writing file %v: %v", path, err)
		return
	}
}

func main() {
	http.HandleFunc("/", handler) // Register the handler for the root path
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil) // Start the server on port 8080
}
