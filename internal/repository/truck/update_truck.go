package truck

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) UpdateTruck(ctx context.Context, userID int64, truck *models.Truck) error {
	query := `
		UPDATE trucks 
		SET truck_name = $1, license_plate = $2, make = $3, model = $4, year = $5, 
			color = $6, length = $7, width = $8, height = $9, max_weight = $10, 
			truck_type = $11, climate_control = $12, liftgate = $13, pallet_jack = $14, 
			security_system = $15, refrigerated = $16, furniture_pads = $17, updated_at = NOW()
		WHERE id = $18 AND user_id = $19`

	_, err := r.db.Exec(ctx, query,
		truck.TruckName, truck.LicensePlate, truck.Make, truck.Model, truck.Year,
		truck.Color, truck.Length, truck.Width, truck.Height, truck.MaxWeight,
		truck.TruckType, truck.ClimateControl, truck.Liftgate, truck.PalletJack,
		truck.SecuritySystem, truck.Refrigerated, truck.FurniturePads,
		truck.ID, userID)

	return err
}
