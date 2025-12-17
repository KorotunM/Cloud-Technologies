package database

import (
	"context"
	"log"
	"pragma/internal/models"
)

// GetUserByID - получение пользователя по ID
func GetUserByID(userID int) (models.User, error) {
	query := `
		SELECT id, email, password_hash, auth_provider, provider_id 
		FROM users 
		WHERE id = $1
	`

	var user models.User

	err := DB.QueryRow(context.Background(), query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.AuthProvider,
		&user.ProviderID,
	)

	if err != nil {
		log.Printf("Ошибка при получении пользователя с ID %d: %v", userID, err)
		return models.User{}, err
	}

	return user, nil
}
