package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"photobooth-be/internal/helper"
	"photobooth-be/internal/middleware"
	"photobooth-be/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	authService service.AuthService
	oauth       *oauth2.Config
}

func NewAuthHandler(
	authService service.AuthService,
	oauth *oauth2.Config,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		oauth:       oauth,
	}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {

	url := h.oauth.AuthCodeURL("state-token")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(
	w http.ResponseWriter,
	r *http.Request,
) {

	code := r.URL.Query().Get("code")

	user, err := h.authService.HandleGoogleLogin(code)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	accessToken, err := helper.GenerateAccessToken(
		user.ID,
		user.Email,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	refreshToken, err := helper.GenerateRefreshToken(
		user.ID,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   604800,
	})

	redirectURL :=
		"http://localhost:5173/oauth-success?token=" +
			accessToken

	http.Redirect(
		w,
		r,
		redirectURL,
		http.StatusTemporaryRedirect,
	)
}

func (h *AuthHandler) Profile(
	w http.ResponseWriter,
	r *http.Request,
) {

	user := r.Context().Value(
		middleware.UserContextKey,
	)

	claims := user.(*middleware.JWTClaims)

	json.NewEncoder(w).Encode(map[string]any{
		"user_id": claims.UserID,
		"email":   claims.Email,
	})
}

func (h *AuthHandler) RefreshToken(
	w http.ResponseWriter,
	r *http.Request,
) {

	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		http.Error(
			w,
			"refresh token missing",
			http.StatusUnauthorized,
		)
		return
	}

	refreshToken := cookie.Value

	claims := &helper.RefreshTokenClaims{}

	token, err := jwt.ParseWithClaims(
		refreshToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {

			return []byte(
				os.Getenv("SESSION_SECRET"),
			), nil
		},
	)

	if err != nil || !token.Valid {

		http.Error(
			w,
			"invalid refresh token",
			http.StatusUnauthorized,
		)

		return
	}

	newAccessToken, err :=
		helper.GenerateAccessToken(
			claims.UserID,
			"",
		)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}
