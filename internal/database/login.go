package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"pragma/internal/models"
)

// AuthenticateUser - проверка email и пароля
func AuthenticateUser(email, password string) (*models.User, error) {
	hashedPassword := hashPassword(password)

	query := `
        SELECT id, email, password_hash, auth_provider, provider_id 
        FROM users 
        WHERE email = $1 AND password_hash = $2
    `

	var user models.User

	err := DB.QueryRow(context.Background(), query, email, hashedPassword).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.AuthProvider, &user.ProviderID,
	)

	if err != nil {
		log.Printf("Ошибка при проверке пользователя: %v", err)
		return nil, fmt.Errorf("Неверные данные для входа")
	}

	log.Printf("Пользователь %s успешно аутентифицирован", user.Email)
	return &user, nil
}

// Хеширование пароля (уже было)
func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
