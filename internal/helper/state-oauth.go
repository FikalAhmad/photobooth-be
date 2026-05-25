package helper

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 menit saja
		HttpOnly: true,
		Secure:   false, // Set true jika sudah menggunakan HTTPS
	})

	return state
}
