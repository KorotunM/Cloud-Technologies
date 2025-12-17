package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"pragma/internal/cache"
	"pragma/internal/database"
	"pragma/internal/models"
	"pragma/internal/service"
	"pragma/internal/storage"
)

const cacheDuration = 5 * time.Minute

func GetUniversitiesHandler(w http.ResponseWriter, r *http.Request, redisClient *cache.RedisClient, storageClient *storage.Client) {
	ctx := r.Context()

	city := r.URL.Query().Get("city")
	form := r.URL.Query().Get("form")
	totalscoreStr := r.URL.Query().Get("totalScore")
	totalscore, _ := strconv.Atoi(totalscoreStr)
	dormitory := r.URL.Query().Get("dormitory") == "yes"
	military := r.URL.Query().Get("military") == "yes"
	budget := r.URL.Query().Get("budget") == "yes"
	paid := r.URL.Query().Get("paid") == "yes"

	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	subjectIDs := r.URL.Query()["subjects"]
	selectedSubjects := []int{}
	for _, idStr := range subjectIDs {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			selectedSubjects = append(selectedSubjects, id)
		}
	}

	cacheKey := "universities:" + city + ":" + form + ":" + strconv.Itoa(totalscore) + ":" +
		strconv.FormatBool(dormitory) + ":" + strconv.FormatBool(military) + ":" +
		strconv.FormatBool(budget) + ":" + strconv.FormatBool(paid) + ":" + strconv.Itoa(offset)

	if cachedData, err := redisClient.Get(cacheKey); err == nil && cachedData != "" {
		var pageData models.PageData
		if err := json.Unmarshal([]byte(cachedData), &pageData); err == nil {
			pageData.Universities = enrichImageURLs(ctx, pageData.Universities, storageClient)
			renderUniversities(w, pageData)
			return
		}
	}

	validationErrors := service.ValidateInput(city, totalscoreStr)
	if service.HasErrors(validationErrors) {
		subjects := database.GetSubject()

		pageData := models.PageData{
			Universities: nil,
			Subjects:     subjects,
			Errors:       validationErrors,
			City:         city,
			Form:         form,
			TotalScore:   totalscore,
			Dormitory:    dormitory,
			Military:     military,
			Budget:       budget,
			Paid:         paid,
		}

		renderUniversities(w, pageData)
		return
	}

	if userID, err := getUserIDFromSession(r); err == nil {
		filterData := map[string]interface{}{
			"city":       city,
			"form":       form,
			"totalScore": totalscore,
			"dormitory":  dormitory,
			"military":   military,
			"budget":     budget,
			"paid":       paid,
			"subjects":   selectedSubjects,
		}
		if filterJSON, err := json.Marshal(filterData); err == nil {
			database.SaveSearchHistory(userID, string(filterJSON))
		}
	}

	universities, faculties, directions, hasMore := database.GetFilteredData(city, form, totalscore, dormitory, military, budget, paid, selectedSubjects, offset)
	universities = enrichImageURLs(ctx, universities, storageClient)
	subjects := database.GetSubject()

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

	if data, err := json.Marshal(pageData); err == nil {
		if err := redisClient.Set(cacheKey, data, cacheDuration); err != nil {
			log.Printf("cache set failed: %v", err)
		}
	}

	renderUniversities(w, pageData)
}

func renderUniversities(w http.ResponseWriter, pageData models.PageData) {
	tmpl, err := template.ParseFiles("./web/templates/index.html", "./web/templates/university.html")
	if err != nil {
		log.Printf("failed to parse templates: %v", err)
		http.Error(w, "template parsing error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "index", pageData); err != nil {
		log.Printf("failed to render template: %v", err)
		http.Error(w, "template render error", http.StatusInternalServerError)
	}
}

func enrichImageURLs(ctx context.Context, universities []models.University, storageClient *storage.Client) []models.University {
	if storageClient == nil {
		return universities
	}

	for i := range universities {
		url, err := storageClient.ImageURL(ctx, universities[i].ImageKey)
		if err != nil {
			log.Printf("failed to build image URL for university %d: %v", universities[i].ID, err)
			continue
		}
		universities[i].ImageURL = url
	}

	return universities
}
