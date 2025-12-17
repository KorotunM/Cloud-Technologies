package handlers

import (
	"encoding/json"
	"net/http"
	"pragma/internal/database"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Email string `json:"email"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	user, err := database.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Неверные данные для входа", http.StatusUnauthorized)
		return
	}

	err = setUserSession(w, r, user.ID)
	if err != nil {
		http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Email: user.Email}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
