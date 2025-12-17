package database

import (
	"context"
	"encoding/json"
	"log"
	"pragma/internal/models"
)

func GetSearchHistory(userID int) []models.SearchHistory {
	query := `
		SELECT id, filters, searched_at 
		FROM search_history 
		WHERE user_id = $1 
		ORDER BY searched_at DESC 
		LIMIT 20
	`

	rows, err := DB.Query(context.Background(), query, userID)
	if err != nil {
		log.Printf("Ошибка при получении истории поиска: %v", err)
		return nil
	}
	defer rows.Close()

	var history []models.SearchHistory

	for rows.Next() {
		var entry models.SearchHistory

		err := rows.Scan(&entry.ID, &entry.Filters, &entry.SearchedAt)
		if err != nil {
			log.Printf("Ошибка сканирования истории поиска: %v", err)
			continue
		}

		// Парсим JSON-строку `Filters` в структуру `FiltersData`
		var filters models.Filters
		if err := json.Unmarshal([]byte(entry.Filters), &filters); err != nil {
			log.Printf("Ошибка парсинга фильтров: %v", err)
			continue
		}

		entry.FiltersData = filters

		history = append(history, entry)
	}

	return history
}
