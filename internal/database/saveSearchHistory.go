package database

import (
	"context"
	"log"
	"time"
)

func SaveSearchHistory(userID int, filters string) {
	query := `
		INSERT INTO search_history (user_id, filters, searched_at)
		VALUES ($1, $2, $3)
	`
	_, err := DB.Exec(context.Background(), query, userID, filters, time.Now())
	if err != nil {
		log.Printf("Ошибка при сохранении истории поиска: %v", err)
	}
}
