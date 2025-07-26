package truck

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) CreateTruck(ctx context.Context, truck *models.Truck) (*models.Truck, error) {
	query := `
		INSERT INTO trucks (
			user_id, truck_name, license_plate, make, model, year, color,
			length, width, height, max_weight, truck_type,
			climate_control, liftgate, pallet_jack, security_system,
			refrigerated, furniture_pads, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, NOW(), NOW()
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx, query,
		truck.UserID, truck.TruckName, truck.LicensePlate, truck.Make, truck.Model, truck.Year, truck.Color,
		truck.Length, truck.Width, truck.Height, truck.MaxWeight, truck.TruckType,
		truck.ClimateControl, truck.Liftgate, truck.PalletJack, truck.SecuritySystem,
		truck.Refrigerated, truck.FurniturePads,
	).Scan(&truck.ID, &truck.CreatedAt, &truck.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return truck, nil
}
