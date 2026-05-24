package main

import (
	"log"
	"net/http"
	"os"
	"photobooth-be/internal/config"
	"photobooth-be/internal/handler"
	"photobooth-be/internal/middleware"
	"photobooth-be/internal/repository"
	"photobooth-be/internal/service"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {

	config.Load()

	db, err := config.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(
		userRepo,
		oauthConfig,
	)

	authHandler := handler.NewAuthHandler(
		authService,
		oauthConfig,
	)

	http.HandleFunc("/auth/google", authHandler.GoogleLogin)

	http.HandleFunc(
		"/auth/google/callback",
		authHandler.GoogleCallback,
	)
	http.HandleFunc(
		"/auth/refresh",
		authHandler.RefreshToken,
	)

	protectedHandler := http.HandlerFunc(
		authHandler.Profile,
	)

	http.Handle(
		"/profile",
		middleware.AuthMiddleware(
			protectedHandler,
		),
	)

	log.Println("server running on :8080")

	http.ListenAndServe(":8080", nil)
}
