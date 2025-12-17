package handlers

import (
	"net/http"
	"pragma/internal/config"

	"golang.org/x/oauth2"
)

func GoogleRedirectHandler(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOAuthConfig().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
