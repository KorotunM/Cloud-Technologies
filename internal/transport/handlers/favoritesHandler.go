package handlers

import (
	"log"
	"net/http"
	"pragma/internal/database"
	"text/template"
)

func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	// Получение данных для избранного
	favorites := database.GetFavorites(userID)

	favorites.Option = "favorite"

	tmpl, err := template.ParseFiles("./web/templates/main.html", "./web/templates/university.html")
	if err != nil {
		http.Error(w, "Не удалось загрузить шаблоны", http.StatusInternalServerError)
		log.Printf("Ошибка загрузки шаблона: %v", err)
		return
	}

	// Рендеринг страницы
	err = tmpl.ExecuteTemplate(w, "main", favorites)
	if err != nil {
		http.Error(w, "Ошибка рендеринга страницы", http.StatusInternalServerError)
		log.Printf("Ошибка рендеринга страницы %v", err)
		return
	}
}
