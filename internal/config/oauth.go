package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleOAuthConfig() *oauth2.Config {
	cfg := Get().OAuth
	return &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func VKOAuthConfig() *oauth2.Config {
	cfg := Get().OAuth
	return &oauth2.Config{
		ClientID:     cfg.VKClientID,
		ClientSecret: cfg.VKClientSecret,
		RedirectURL:  cfg.VKRedirectURL,
		Scopes:       []string{"email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.vk.com/authorize",
			TokenURL: "https://oauth.vk.com/access_token",
		},
	}
}

func VerifyTelegramAuth(query url.Values, botToken string) bool {
	checkHash := query.Get("hash")
	query.Del("hash")

	var keys []string
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var dataStrings []string
	for _, key := range keys {
		dataStrings = append(dataStrings, fmt.Sprintf("%s=%s", key, query.Get(key)))
	}
	dataCheckString := strings.Join(dataStrings, "\n")

	secretKey := sha256.Sum256([]byte(botToken))
	hmac := hmac.New(sha256.New, secretKey[:])
	hmac.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(hmac.Sum(nil))

	return calculatedHash == checkHash
}
