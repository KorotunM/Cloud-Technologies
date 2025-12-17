package database

import (
	"context"
	"log"
)

// SaveOAuthUser - сохраняет пользователя, авторизованного через сторонний сервис
func SaveOAuthUser(email, provider, providerID string) (string, error) {
	var userEmail string

	// Проверим, существует ли пользователь с данным email и provider
	query := `
		SELECT email FROM users 
		WHERE email = $1 AND auth_provider = $2
	`
	err := DB.QueryRow(context.Background(), query, email, provider).Scan(&userEmail)

	if err != nil {
		// Пользователь не найден, создаем нового
		insertQuery := `
   				 INSERT INTO users (email, auth_provider, provider_id)
   				 VALUES ($1, $2, $3) 
  				  RETURNING email
				`
		err = DB.QueryRow(context.Background(), insertQuery, email, provider, providerID).Scan(&userEmail)
		if err != nil {
			log.Printf("Ошибка при добавлении OAuth пользователя: %v", err)
			return "0", err
		}
		log.Printf("Новый OAuth пользователь создан с ID: %v", userEmail)
		return userEmail, nil

	} else {
		log.Printf("OAuth пользователь найден с ID: %v", userEmail)
	}

	return userEmail, nil
}
