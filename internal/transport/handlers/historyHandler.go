package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"pragma/internal/cache"
	"pragma/internal/database"
	"pragma/internal/models"
	"strconv"
	"time"
)

const historyCacheDuration = 5 * time.Minute

type HistoryPageData struct {
	Option string
	Data   []models.SearchHistory
}

func HistoryHandler(w http.ResponseWriter, r *http.Request, redisClient *cache.RedisClient) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	// Создаём ключ для кеша
	cacheKey := "user_history:" + strconv.Itoa(userID)

	// Попытка загрузить данные из кеша
	cachedData, err := redisClient.Get(cacheKey)
	if err == nil && cachedData != "" {
		// Если данные найдены в кеше, десериализуем их
		var cachedHistory []models.SearchHistory
		if err := json.Unmarshal([]byte(cachedData), &cachedHistory); err == nil {
			log.Println("Данные загружены из кеша")
			renderHistoryPage(w, cachedHistory)
			return
		}
	}

	// Если данных нет в кеше — получаем из базы
	history := database.GetSearchHistory(userID)

	// Сохранение данных в кеш
	historyData, err := json.Marshal(history)
	if err == nil {
		redisClient.Set(cacheKey, historyData, historyCacheDuration)
	}

	renderHistoryPage(w, history)
}

// renderHistoryPage - рендерит страницу истории
func renderHistoryPage(w http.ResponseWriter, history []models.SearchHistory) {
	pageData := HistoryPageData{
		Option: "history",
		Data:   history,
	}

	tmpl, err := template.ParseFiles(
		"./web/templates/main.html",
		"./web/templates/history.html",
		"./web/templates/compare.html",
		"./web/templates/university.html",
	)
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Printf("Ошибка загрузки шаблона: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "main", pageData)
	if err != nil {
		log.Printf("Ошибка рендеринга страницы: %v", err)
		http.Error(w, "Ошибка рендеринга страницы", http.StatusInternalServerError)
	}
}
