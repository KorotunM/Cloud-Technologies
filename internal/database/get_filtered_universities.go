package database

import (
	"context"
	"log"
	"pragma/internal/models"
	"strconv"
)

const Limit = 4

func GetFilteredUniversities(city string, dormitory, military bool, offset int) []models.University {
	query := `
		SELECT 
			u.id,
			u.name,
			u.city,
			(u.has_dormitory <> 0) AS has_dormitory,
			(u.has_military <> 0) AS has_military,
			r.rank_position,
			u.photo_key
		FROM universities u
		LEFT JOIN university_ratings r ON u.id = r.university_id
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if city != "" {
		query += " AND u.city ILIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+city+"%")
		argIndex++
	}

	if dormitory {
		// has_dormitory stored as integer; non-zero means true
		query += " AND u.has_dormitory <> 0"
	}

	if military {
		// has_military stored as integer; non-zero means true
		query += " AND u.has_military <> 0"
	}

	query += " GROUP BY u.id, r.rank_position, u.photo_key"
	query += " ORDER BY r.rank_position ASC NULLS LAST, u.name ASC"

	query += " LIMIT $" + strconv.Itoa(argIndex)
	args = append(args, Limit)
	argIndex++

	query += " OFFSET $" + strconv.Itoa(argIndex)
	args = append(args, offset)

	rows, err := DB.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("failed to fetch universities: %v", err)
		return nil
	}
	defer rows.Close()

	var universities []models.University
	for rows.Next() {
		var (
			uni          models.University
			rankPosition *int
			photoKey     *string
		)

		err := rows.Scan(
			&uni.ID,
			&uni.Name,
			&uni.City,
			&uni.HasDormitory,
			&uni.HasMilitary,
			&rankPosition,
			&photoKey,
		)
		if err != nil {
			log.Printf("failed to scan university: %v", err)
			continue
		}

		if rankPosition != nil {
			uni.Rating = *rankPosition
		}

		if photoKey != nil {
			uni.ImageKey = *photoKey
		}

		universities = append(universities, uni)
	}

	return universities
}
