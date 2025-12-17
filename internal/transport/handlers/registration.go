package handlers

import (
	"encoding/json"
	"net/http"
	"pragma/internal/database"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Email  string `json:"email,omitempty"`
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response := RegisterResponse{
			Status: "error",
			Error:  "Некорректные данные",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Email == "" || req.Password == "" {
		response := RegisterResponse{
			Status: "error",
			Error:  "Email и пароль обязательны",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	var userID int
	userID, err = database.InsertUser(req.Email, req.Password)
	if err != nil {
		response := RegisterResponse{
			Status: "error",
			Error:  "Ошибка при регистрации: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Сохранение сессии
	err = setUserSession(w, r, userID)
	if err != nil {
		http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		Email:  req.Email,
		Status: "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
