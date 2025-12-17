package database

import (
	"context"
	"log"
	"pragma/internal/models"
)

func AddToFavorites(userID, universityID int) error {
	query := `
		INSERT INTO favorites (user_id, university_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := DB.Exec(context.Background(), query, userID, universityID)
	if err != nil {
		log.Println("Ошибка при добавлении в избранное:", err)
		return err
	}
	return nil
}

func RemoveFromFavorites(userID, universityID int) error {
	query := `
		DELETE FROM favorites WHERE user_id = $1 AND university_id = $2
	`
	_, err := DB.Exec(context.Background(), query, userID, universityID)
	if err != nil {
		log.Println("Ошибка при удалении из избранного:", err)
		return err
	}
	return nil
}

// GetFavorites - получение данных об избранных университетах
func GetFavorites(userID int) models.Favorites {
	var universityIDs []int

	// Получение IDs университетов из избранного
	query := `
		SELECT university_id FROM favorites WHERE user_id = $1
	`
	rows, err := DB.Query(context.Background(), query, userID)
	if err != nil {
		log.Printf("Ошибка при получении избранного: %v", err)
		return models.Favorites{}
	}
	defer rows.Close()

	for rows.Next() {
		var universityID int
		if err := rows.Scan(&universityID); err != nil {
			log.Printf("Ошибка при сканировании university_id: %v", err)
			continue
		}
		universityIDs = append(universityIDs, universityID)
	}

	// Если нет избранных университетов, возвращаем пустой PageData
	if len(universityIDs) == 0 {
		return models.Favorites{}
	}

	// Получаем данные об университетах
	universities := GetFilteredUniversities("", false, false, 0)

	// Извлекаем факультеты для этих университетов
	faculties := GetFaculties(universityIDs, []int{})

	// Извлекаем направления для этих факультетов
	var facultyIDs []int
	for _, faculty := range faculties {
		facultyIDs = append(facultyIDs, faculty.ID)
	}

	directions := GetDirections(facultyIDs, "all", 0, true, true)

	return models.Favorites{
		Universities: universities,
		Faculties:    faculties,
		Directions:   directions,
	}
}
