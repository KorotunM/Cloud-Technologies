package database

import (
	"log"
	"pragma/internal/models"
)

// GetFilteredData - объединяет данные об университетах, факультетах и направлениях
func GetFilteredData(city, form string, totalScore int, dormitory, military, budget, paid bool, selectedSubjects []int, offset int) ([]models.University, []models.Faculty, []models.Direction, bool) {
	// Университеты
	universities := GetFilteredUniversities(city, dormitory, military, offset)
	// Извлекаем IDs университетов для фильтрации факультетов
	universityIDs := []int{}
	for _, uni := range universities {
		universityIDs = append(universityIDs, uni.ID)
	}

	// Получение факультетов для выбранных университетов
	faculties := GetFaculties(universityIDs, selectedSubjects)

	// Извлекаем IDs факультетов для фильтрации направлений
	facultyIDs := []int{}
	for _, fac := range faculties {
		facultyIDs = append(facultyIDs, fac.ID)
	}

	// Получение направлений для выбранных факультетов
	directions := GetDirections(facultyIDs, form, totalScore, budget, paid)

	// Логирование для отладки
	log.Printf("Universities: %d, Faculties: %d, Directions: %d", len(universities), len(faculties), len(directions))

	hasMore := len(universities) >= Limit

	return universities, faculties, directions, hasMore
}
