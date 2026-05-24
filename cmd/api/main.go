package main

import (
	"fmt"
	"log"
	"net/http"
	"photobooth-be/internal/handle"

	"github.com/joho/godotenv"
)

func main() {
	// Memuat file .env menggunakan godotenv
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Tidak dapat memuat file .env, menggunakan environment variables sistem")
	}

	http.HandleFunc("/login", handle.HandleLogin)
	http.HandleFunc("/callback", handle.HandleCallback)

	fmt.Println("starting server on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
