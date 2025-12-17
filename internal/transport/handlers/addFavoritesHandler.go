package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pragma/internal/database"
)

type FavoriteRequest struct {
	UniversityID int `json:"university_id"`
}

// Добавление в избранное
func AddToFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения данных", http.StatusBadRequest)
		return
	}

	// Логируем тело запроса
	log.Printf("Тело запроса: %s", body)

	// Сохраняем тело запроса для повторного чтения
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Декодируем JSON
	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded data: UniversityID=%d", req.UniversityID)

	if req.UniversityID == 0 {
		http.Error(w, "Некорректные данные: UniversityID не может быть 0", http.StatusBadRequest)
		return
	}

	// Получение userID из сессии
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	// Добавление в избранное
	err = database.AddToFavorites(userID, req.UniversityID)
	if err != nil {
		http.Error(w, "Ошибка при добавлении в избранное", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешное добавление: userID=%d, universityID=%d", userID, req.UniversityID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
