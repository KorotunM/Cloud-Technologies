package handlers

import (
	"net/http"
	"pragma/internal/config"
)

func VKRedirectHandler(w http.ResponseWriter, r *http.Request) {
	url := config.VKOAuthConfig().AuthCodeURL("state-token")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
