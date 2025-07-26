package truck

// import (
// 	"encoding/json"
// 	"moveshare/internal/domain"
// 	"moveshare/internal/dto"
// 	"moveshare/internal/service"
// 	"net/http"
// )

// func CreateTruck(truckService service.TruckService) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userID, ok := r.Context().Value("userID").(int64)
// 		if !ok {
// 			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
// 			return
// 		}

// 		var req dto.TruckRequest
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
// 			return
// 		}

// 		truck := &domain.Truck{
// 			UserID:         userID,
// 			TruckName:      req.TruckName,
// 			LicensePlate:   req.LicensePlate,
// 			Make:           req.Make,
// 			Model:          req.Model,
// 			Year:           req.Year,
// 			Color:          req.Color,
// 			Length:         req.Length,
// 			Width:          req.Width,
// 			Height:         req.Height,
// 			MaxWeight:      req.MaxWeight,
// 			TruckType:      req.TruckType,
// 			ClimateControl: req.ClimateControl,
// 			Liftgate:       req.Liftgate,
// 			PalletJack:     req.PalletJack,
// 			SecuritySystem: req.SecuritySystem,
// 			Refrigerated:   req.Refrigerated,
// 			FurniturePads:  req.FurniturePads,
// 		}

// 		return truck

// 	}
// }
