package truck

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetTruckByID(ctx context.Context, id int64) (*models.Truck, error) {
	query := `
		SELECT id, user_id, truck_name, license_plate, make, model, year, color,
			   length, width, height, max_weight, truck_type, climate_control,
			   liftgate, pallet_jack, security_system, refrigerated, furniture_pads,
			   created_at, updated_at
		FROM trucks
		WHERE id = $1`

	truck := &models.Truck{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&truck.ID,
		&truck.UserID,
		&truck.TruckName,
		&truck.LicensePlate,
		&truck.Make,
		&truck.Model,
		&truck.Year,
		&truck.Color,
		&truck.Length,
		&truck.Width,
		&truck.Height,
		&truck.MaxWeight,
		&truck.TruckType,
		&truck.ClimateControl,
		&truck.Liftgate,
		&truck.PalletJack,
		&truck.SecuritySystem,
		&truck.Refrigerated,
		&truck.FurniturePads,
		&truck.CreatedAt,
		&truck.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return truck, nil
}
