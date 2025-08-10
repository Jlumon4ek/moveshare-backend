package repository

import (
	"context"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) GetAllStates(ctx context.Context) ([]models.State, error) {
	query := `SELECT id, name FROM states ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []models.State
	for rows.Next() {
		var state models.State
		err := rows.Scan(&state.ID, &state.Name)
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return states, nil
}

func (r *LocationRepository) GetCities(ctx context.Context, stateID *int64) ([]models.CityWithState, error) {
	var query string
	var args []interface{}

	if stateID != nil {
		query = `
			SELECT c.id, c.name, c.state_id, s.name as state_name 
			FROM cities c 
			JOIN states s ON c.state_id = s.id 
			WHERE c.state_id = $1 
			ORDER BY c.name`
		args = append(args, *stateID)
	} else {
		query = `
			SELECT c.id, c.name, c.state_id, s.name as state_name 
			FROM cities c 
			JOIN states s ON c.state_id = s.id 
			ORDER BY s.name, c.name`
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.CityWithState
	for rows.Next() {
		var city models.CityWithState
		err := rows.Scan(&city.ID, &city.Name, &city.StateID, &city.StateName)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}
