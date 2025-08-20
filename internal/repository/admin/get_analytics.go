package admin

import (
	"context"
	"moveshare/internal/models"
	"time"
)

func (r *repository) GetTopCompanies(ctx context.Context, days int, limit int) ([]models.TopCompany, error) {
	query := `
		SELECT 
			c.company_name,
			COUNT(j.id) as jobs_count
		FROM companies c
		INNER JOIN users u ON c.user_id = u.id
		INNER JOIN jobs j ON j.contractor_id = u.id
		WHERE j.created_at >= $1
		  AND c.company_name IS NOT NULL 
		  AND c.company_name != ''
		GROUP BY c.company_name
		ORDER BY jobs_count DESC
		LIMIT $2
	`

	cutoffDate := time.Now().AddDate(0, 0, -days)
	rows, err := r.db.Query(ctx, query, cutoffDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.TopCompany
	for rows.Next() {
		var company models.TopCompany
		err := rows.Scan(&company.CompanyName, &company.JobsCount)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return companies, nil
}

func (r *repository) GetBusiestRoutes(ctx context.Context, days int, limit int) ([]models.BusyRoute, error) {
	query := `
		SELECT 
			CONCAT(pickup_city, ', ', pickup_state, ' → ', delivery_city, ', ', delivery_state) as route,
			pickup_address,
			delivery_address,
			COUNT(*) as jobs_count
		FROM jobs 
		WHERE created_at >= $1
		  AND pickup_city IS NOT NULL 
		  AND pickup_state IS NOT NULL
		  AND delivery_city IS NOT NULL 
		  AND delivery_state IS NOT NULL
		GROUP BY pickup_city, pickup_state, delivery_city, delivery_state, pickup_address, delivery_address
		ORDER BY jobs_count DESC
		LIMIT $2
	`

	cutoffDate := time.Now().AddDate(0, 0, -days)
	rows, err := r.db.Query(ctx, query, cutoffDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []models.BusyRoute
	for rows.Next() {
		var route models.BusyRoute
		err := rows.Scan(&route.Route, &route.PickupAddress, &route.DeliveryAddress, &route.JobsCount)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return routes, nil
}