package handlers

import (
	"context"
	"fmt"
	"net/http"

	"pragma/internal/config"
	"pragma/internal/database"
	"pragma/internal/models"
)

func VKCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	oauthCfg := config.VKOAuthConfig()
	token, err := oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "vk auth failed", http.StatusInternalServerError)
		return
	}

	email, _ := token.Extra("email").(string)
	userID, _ := token.Extra("user_id").(float64)

	if email == "" || userID == 0 {
		http.Error(w, "vk profile is missing email", http.StatusBadRequest)
		return
	}

	user := models.User{
		Email:        email,
		AuthProvider: "vk",
		ProviderID:   fmt.Sprintf("%.0f", userID),
	}

	if _, err := database.SaveOAuthUser(user.Email, user.AuthProvider, user.ProviderID); err != nil {
		http.Error(w, "vk user save failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
