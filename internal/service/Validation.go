package service

import (
	"pragma/internal/models"
	"regexp"
	"strings"
)

// ValidateInput выполняет проверку данных
func ValidateInput(city, score string) models.Errors {
	errors := models.Errors{}

	// Регулярные выражения для проверки
	rusLetters := regexp.MustCompile(`^[а-яА-ЯёЁ\s\-]*$`)
	numeric := regexp.MustCompile(`^[0-9]*$`)

	if city != "" && !rusLetters.MatchString(city) {
		errors.CityError = "Город должен содержать только русские буквы, пробелы и тире."
	}

	if score != "" && !numeric.MatchString(score) {
		errors.ScoreError = "Сумма баллов должна содержать только цифры."
	}

	return errors
}

func HasErrors(errors models.Errors) bool {
	return strings.TrimSpace(errors.CityError) != "" ||
		strings.TrimSpace(errors.ScoreError) != ""
}
