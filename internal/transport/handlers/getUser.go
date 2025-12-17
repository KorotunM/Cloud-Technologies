package handlers

import (
	"encoding/json"
	"net/http"
	"pragma/internal/database"
)

type UserResponse struct {
	Email string `json:"email"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении данных пользователя", http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
