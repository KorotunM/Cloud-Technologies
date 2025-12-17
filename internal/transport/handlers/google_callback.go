package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"pragma/internal/config"
	"pragma/internal/database"
	"pragma/internal/models"
)

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	oauthCfg := config.GoogleOAuthConfig()

	token, err := oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		http.Redirect(w, r, "/?error=google_auth_failed", http.StatusFound)
		return
	}

	client := oauthCfg.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Redirect(w, r, "/?error=google_profile_request_failed", http.StatusFound)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Redirect(w, r, "/?error=google_profile_decode_failed", http.StatusFound)
		return
	}

	user := models.User{
		Email:        userInfo.Email,
		AuthProvider: "google",
		ProviderID:   userInfo.ID,
	}

	userEmail, err := database.SaveOAuthUser(user.Email, user.AuthProvider, user.ProviderID)
	if err != nil {
		http.Redirect(w, r, "/?error=google_user_persist_failed", http.StatusFound)
		return
	}

	redirectURL := "/?email=" + url.QueryEscape(userEmail) + "&status=success"
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
