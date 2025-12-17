package database

import (
	"context"
	"pragma/internal/models"
)

func GetComparisonData(universityIDs []int) ([]models.ComparisonData, error) {
	var data []models.ComparisonData

	query := `
		SELECT 
			u.id, u.name, u.city, u.has_dormitory, u.has_military,
			COALESCE(SUM(d.budget_places), 0) AS total_budget_places,
			COALESCE(SUM(d.paid_places), 0) AS total_paid_places,
			COALESCE(MIN(d.tuition_fee), 0) AS min_tuition_fee,
			COALESCE(MIN(d.min_score_budget), 0) AS min_score_budget,
			COALESCE(MIN(d.min_score_paid), 0) AS min_score_paid,
			COUNT(DISTINCT d.id) AS programs_count,
			COUNT(DISTINCT f.id) AS specialties_count
		FROM universities u
		LEFT JOIN faculties f ON f.university_id = u.id
		LEFT JOIN directions d ON d.faculty_id = f.id
		WHERE u.id = ANY($1)
		GROUP BY u.id
	`

	rows, err := DB.Query(context.Background(), query, universityIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cd models.ComparisonData
		err := rows.Scan(
			&cd.ID, &cd.Name, &cd.City, &cd.HasDormitory, &cd.HasMilitary,
			&cd.TotalBudgetPlaces, &cd.TotalPaidPlaces, &cd.MinTuitionFee,
			&cd.MinScoreBudget, &cd.MinScorePaid, &cd.ProgramsCount, &cd.SpecialtiesCount,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, cd)
	}

	return data, nil
}
