package handle

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := "https://accounts.google.com/o/oauth2/v2/auth"

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = os.Getenv("REDIRECT_URI")
	}

	state := generateState()

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   false, // set to true in production with HTTPS
		Path:     "/",
	})

	q := url.Values{}
	q.Add("client_id", clientID)
	q.Add("redirect_uri", redirectURI)
	q.Add("response_type", "code")
	q.Add("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	q.Add("access_type", "offline")
	q.Add("state", state)

	http.Redirect(
		w,
		r,
		fmt.Sprintf("%s?%s", authURL, q.Encode()),
		http.StatusFound,
	)
}

// HandleCallback menerima callback dari Google OAuth
func HandleCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Validasi State
	stateParam := r.URL.Query().Get("state")
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateParam == "" || stateCookie.Value != stateParam {
		http.Error(w, "State tidak cocok / invalid state", http.StatusBadRequest)
		return
	}

	// 2. Ambil authorization code dari URL query params
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code tidak ditemukan", http.StatusBadRequest)
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = os.Getenv("REDIRECT_URI")
	}
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // Default frontend port
	}

	// 3. Tukar Code dengan Access Token
	tokenURL := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)
	data.Add("code", code)
	data.Add("redirect_uri", redirectURI)
	data.Add("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		http.Error(w, "Gagal request token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse response token dari Google
	var tokenRes map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		http.Error(w, "Gagal parse token response", http.StatusInternalServerError)
		return
	}

	accessToken, ok := tokenRes["access_token"].(string)
	if !ok {
		http.Error(w, "Access token tidak valid", http.StatusInternalServerError)
		return
	}

	// 4. Gunakan Access Token untuk mengambil data profile user
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		http.Error(w, "Gagal membuat request userinfo", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Gagal mengambil data user", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	userInfo, err := io.ReadAll(userResp.Body)
	if err != nil {
		http.Error(w, "Gagal membaca data user", http.StatusInternalServerError)
		return
	}

	var user map[string]any
	if err := json.Unmarshal(userInfo, &user); err != nil {
		http.Error(w, "Gagal parse data user", http.StatusInternalServerError)
		return
	}

	email, _ := user["email"].(string)
	name, _ := user["name"].(string)

	// Redirect kembali ke frontend dengan query parameters
	redirectBack := fmt.Sprintf("%s?email=%s&name=%s", frontendURL, url.QueryEscape(email), url.QueryEscape(name))
	http.Redirect(w, r, redirectBack, http.StatusTemporaryRedirect)
}
