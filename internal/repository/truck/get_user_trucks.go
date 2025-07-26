package truck

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetUserTrucks(ctx context.Context, userID int64) ([]*models.Truck, error) {
	query := `
		SELECT id, user_id, truck_name, license_plate, make, model, year, color,
			   length, width, height, max_weight, truck_type, climate_control,
			   liftgate, pallet_jack, security_system, refrigerated, furniture_pads,
			   created_at, updated_at
		FROM trucks
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trucks []*models.Truck
	for rows.Next() {
		truck := &models.Truck{}
		err := rows.Scan(
			&truck.ID, &truck.UserID, &truck.TruckName, &truck.LicensePlate,
			&truck.Make, &truck.Model, &truck.Year, &truck.Color,
			&truck.Length, &truck.Width, &truck.Height, &truck.MaxWeight,
			&truck.TruckType, &truck.ClimateControl, &truck.Liftgate,
			&truck.PalletJack, &truck.SecuritySystem, &truck.Refrigerated,
			&truck.FurniturePads, &truck.CreatedAt, &truck.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		trucks = append(trucks, truck)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trucks, nil
}
