package handlers

import (
	"net/http"
	"pragma/internal/database"
	"pragma/internal/models"
	"strconv"
	"strings"
	"text/template"
)

type ComparePageData struct {
	Option string
	Data   []models.ComparisonData
}

func CompareHandler(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		http.Error(w, "Не выбраны университеты для сравнения", http.StatusBadRequest)
		return
	}

	idStrings := strings.Split(idsParam, ",")
	var universityIDs []int

	for _, idStr := range idStrings {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			universityIDs = append(universityIDs, id)
		}
	}

	data, err := database.GetComparisonData(universityIDs)
	if err != nil {
		http.Error(w, "Ошибка при получении данных для сравнения", http.StatusInternalServerError)
		return
	}

	// Формируем структуру-обертку
	pageData := ComparePageData{
		Option: "compare",
		Data:   data,
	}

	tmpl, err := template.ParseFiles("./web/templates/main.html", "./web/templates/compare.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "main", pageData)
	if err != nil {
		http.Error(w, "Ошибка рендеринга страницы", http.StatusInternalServerError)
		return
	}
}
