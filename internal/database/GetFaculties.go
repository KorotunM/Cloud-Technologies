package database

import (
	"context"
	"log"
	"pragma/internal/models"
)

// GetFaculties - возвращает факультеты, подходящие по предметам и университетам
func GetFaculties(universityIDs []int, subjectIDs []int) []models.Faculty {
	// Если нет университетов для фильтрации, возвращаем пустой срез
	if len(universityIDs) == 0 {
		return []models.Faculty{}
	}

	query := `
		SELECT DISTINCT f.id, f.university_id, f.name
		FROM faculties f
	`

	args := []interface{}{universityIDs}

	// Если предметы указаны, подключаем таблицу faculty_subjects
	if len(subjectIDs) > 0 {
		query += " JOIN faculty_subjects fs ON fs.faculty_id = f.id"
		query += " WHERE f.university_id = ANY($1) AND fs.subject_id = ANY($2)"
		args = append(args, subjectIDs)
	} else {
		// Если предметы не указаны, просто фильтруем по университетам
		query += " WHERE f.university_id = ANY($1)"
	}

	// Выполняем запрос
	rows, err := DB.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса к факультетам: %v", err)
		return nil
	}
	defer rows.Close()

	var faculties []models.Faculty

	// Сканируем данные
	for rows.Next() {
		var faculty models.Faculty
		err := rows.Scan(&faculty.ID, &faculty.UniversityID, &faculty.Name)
		if err != nil {
			log.Printf("Ошибка при сканировании факультета: %v", err)
			continue
		}
		faculties = append(faculties, faculty)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка при обработке строк факультетов: %v", err)
	}

	return faculties
}
