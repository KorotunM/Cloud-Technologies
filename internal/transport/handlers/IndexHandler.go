package handlers

import (
	"html/template"
	"log"
	"net/http"
	"pragma/internal/database"
	"pragma/internal/models"
	"strconv"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров из запроса
	city := r.URL.Query().Get("city")
	form := r.URL.Query().Get("form")
	totalscore, _ := strconv.Atoi(r.URL.Query().Get("totalScore"))
	dormitory := r.URL.Query().Get("dormitory") == "yes"
	military := r.URL.Query().Get("military") == "yes"
	budget := r.URL.Query().Get("budget") == "yes"
	paid := r.URL.Query().Get("paid") == "yes"

	// Параметр offset
	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset <= 0 {
		offset = 0
	}

	// Получение выбранных предметов
	subjectIDs := r.URL.Query()["subjects"] // Массив значений параметра `subjects`
	selectedSubjects := []int{}

	// Конвертируем `subjectIDs` в `[]int`
	for _, idStr := range subjectIDs {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			selectedSubjects = append(selectedSubjects, id)
		}
	}

	// Логирование для отладки
	log.Printf("Selected Subjects: %v", selectedSubjects)
	log.Printf("Offset: %d", offset)

	// Получение данных с учётом фильтров и пагинации
	universities, faculties, directions, hasMore := database.GetFilteredData(city, form, totalscore, dormitory, military, budget, paid, selectedSubjects, offset)

	// Получаем список предметов для отображения в форме
	subjects := database.GetSubject()

	// Создание структуры данных для страницы
	pageData := models.PageData{
		Universities:     universities,
		Faculties:        faculties,
		Directions:       directions,
		Subjects:         subjects,
		City:             city,
		Form:             form,
		TotalScore:       totalscore,
		Dormitory:        dormitory,
		Military:         military,
		Budget:           budget,
		Paid:             paid,
		SelectedSubjects: selectedSubjects,
		Offset:           offset,
		HasMore:          hasMore,
	}

	// Создание нового шаблона и загрузка всех необходимых файлов
	tmpl, err := template.ParseFiles("./web/templates/index.html", "./web/templates/university.html")
	if err != nil {
		http.Error(w, "Не удалось загрузить шаблоны", http.StatusInternalServerError)
		log.Printf("Ошибка загрузки шаблона: %v", err)
		return
	}

	// Выполнение шаблона "index" и передача данных об университетах
	err = tmpl.ExecuteTemplate(w, "index", pageData)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
		log.Printf("Ошибка выполнения шаблона: %v", err)
	}
}
