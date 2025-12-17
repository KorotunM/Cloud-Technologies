package database

import (
	"context"
	"log"
	"pragma/internal/models"
)

// GetSubject - возвращает все предметы, которые можно сдать на ЕГЭ
func GetSubject() []models.Subject {
	query := `SELECT id, name FROM subjects`

	// Выполняем запрос
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("Ошибка выполнения запроса к предметам: %v", err)
		return nil
	}
	defer rows.Close()

	var subjects []models.Subject

	// Сканируем данные
	for rows.Next() {
		var sub models.Subject
		err := rows.Scan(&sub.ID, &sub.Name)
		if err != nil {
			log.Printf("Ошибка при сканировании предмета: %v", err)
			continue
		}
		subjects = append(subjects, sub)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка при обработке строк предметов: %v", err)
	}

	return subjects
}
