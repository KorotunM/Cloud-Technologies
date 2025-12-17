package database

import (
	"context"
	"log"
	"pragma/internal/models"
	"strconv"
)

func GetDirections(facultyIDs []int, form string, totalScore int, budget, paid bool) []models.Direction {
	if len(facultyIDs) == 0 {
		return []models.Direction{}
	}

	query := `
		SELECT d.id, d.faculty_id, d.name, d.tuition_fee, d.min_score_budget, d.min_score_paid, d.budget_places, d.paid_places
		FROM directions d
		WHERE d.faculty_id = ANY($1)
	`

	args := []interface{}{facultyIDs}
	paramIndex := 2

	// Фильтрация по форме обучения
	if form != "all" && form != "" {
		query += " AND d.name ILIKE $" + strconv.Itoa(paramIndex)
		args = append(args, "%"+form+"%")
		paramIndex++
	}

	// Фильтрация по бюджету
	if budget {
		query += " AND d.min_score_budget <= $" + strconv.Itoa(paramIndex) + " AND d.budget_places > 0"
		args = append(args, totalScore)
		paramIndex++
	}

	// Фильтрация по платному
	if paid {
		query += " AND d.min_score_paid <= $" + strconv.Itoa(paramIndex) + " AND d.paid_places > 0"
		args = append(args, totalScore)
	}

	rows, err := DB.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса к направлениям: %v", err)
		return nil
	}
	defer rows.Close()

	var directions []models.Direction

	for rows.Next() {
		var dir models.Direction
		err := rows.Scan(&dir.ID, &dir.FacultyID, &dir.Name, &dir.TuitionFee, &dir.MinScoreBudget, &dir.MinScorePaid, &dir.BudgetPlaces, &dir.PaidPlaces)
		if err != nil {
			log.Printf("Ошибка при сканировании направления: %v", err)
			continue
		}
		directions = append(directions, dir)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка при обработке строк направлений: %v", err)
	}

	return directions
}
