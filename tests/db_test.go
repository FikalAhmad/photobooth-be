package tests

import (
	database "photobooth-be/internal/config"
	"testing"
)

func TestDbPing(t *testing.T) {
	db, err := database.NewPostgresDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Expected database to ping successfully, but got an error: %v", err)
	}
}
