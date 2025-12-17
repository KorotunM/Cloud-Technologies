package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"pragma/internal/config"
	"pragma/internal/database"
	"pragma/internal/models"
)

func TelegramCallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	botToken := config.Get().OAuth.TelegramBotToken

	if botToken == "" || !config.VerifyTelegramAuth(query, botToken) {
		http.Redirect(w, r, "/?error=telegram_signature_invalid", http.StatusFound)
		return
	}

	userID := query.Get("id")
	username := query.Get("username")

	if username == "" {
		http.Redirect(w, r, "/?error=telegram_username_missing", http.StatusFound)
		return
	}

	userEmail := fmt.Sprintf("%s@telegram", username)

	user := models.User{
		Email:        userEmail,
		AuthProvider: "telegram",
		ProviderID:   userID,
	}

	email, err := database.SaveOAuthUser(user.Email, user.AuthProvider, user.ProviderID)
	if err != nil {
		http.Redirect(w, r, "/?error=telegram_user_save_failed", http.StatusFound)
		return
	}

	redirectURL := "/?email=" + url.QueryEscape(email) + "&status=success"
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
