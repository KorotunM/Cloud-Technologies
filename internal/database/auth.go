package database

import (
	"context"
	"fmt"
	"log"
)

// InsertUser - вставка нового пользователя и возврат его ID
func InsertUser(email, password string) (int, error) {
	hashedPassword := hashPassword(password)
	var userID int

	// Проверим, существует ли уже пользователь с таким email
	checkQuery := `
		SELECT id FROM users 
		WHERE email = $1
	`
	err := DB.QueryRow(context.Background(), checkQuery, email).Scan(&userID)
	if err == nil {
		log.Printf("Пользователь с email %s уже существует с ID: %d", email, userID)
		return 0, fmt.Errorf("пользователь с таким email уже существует")
	}

	// Вставка нового пользователя
	insertQuery := `
        INSERT INTO users (email, password_hash, auth_provider, provider_id) 
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	err = DB.QueryRow(context.Background(), insertQuery, email, hashedPassword, "local", "").Scan(&userID)
	if err != nil {
		log.Printf("Ошибка при вставке пользователя: %v", err)
		return 0, fmt.Errorf("ошибка при регистрации")
	}

	log.Printf("Пользователь %s успешно зарегистрирован с ID: %d", email, userID)
	return userID, nil
}
