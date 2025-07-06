package handlers

import (
	"encoding/json"
	"mime/multipart"
	"moveshare/internal/repository"
	"moveshare/internal/service"
	"net/http"
	"strconv"
)

// TruckHandler handles HTTP requests for truck operations
type TruckHandler struct {
	truckService service.TruckService
}

// NewTruckHandler creates a new TruckHandler
func NewTruckHandler(truckService service.TruckService) *TruckHandler {
	return &TruckHandler{truckService: truckService}
}

// TruckRequest represents the truck creation request payload
type TruckRequest struct {
	TruckName      string  `json:"truck_name"`
	LicensePlate   string  `json:"license_plate"`
	Make           string  `json:"make"`
	Model          string  `json:"model"`
	Year           int     `json:"year"`
	Color          string  `json:"color"`
	Length         float64 `json:"length"`
	Width          float64 `json:"width"`
	Height         float64 `json:"height"`
	MaxWeight      float64 `json:"max_weight"`
	TruckType      string  `json:"truck_type"`
	ClimateControl bool    `json:"climate_control"`
	Liftgate       bool    `json:"liftgate"`
	PalletJack     bool    `json:"pallet_jack"`
	SecuritySystem bool    `json:"security_system"`
	Refrigerated   bool    `json:"refrigerated"`
	FurniturePads  bool    `json:"furniture_pads"`
}

// TruckResponse represents the truck creation response
type TruckResponse struct {
	Truck *repository.Truck `json:"truck"`
}

// GetUserTrucks handles fetching trucks for the current user
// @Summary Get user trucks
// @Description Fetches all trucks owned by the current user
// @Tags trucks
// @Produce json
// @Success 200 {object} TruckResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks [get]
func (h *TruckHandler) GetUserTrucks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	trucks, err := h.truckService.GetUserTrucks(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := TruckResponse{Truck: nil} // Not used here, but for consistency
	w.Header().Set("Content-Type", "application/json")
	if len(trucks) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "no trucks found"})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string][]repository.Truck{"trucks": trucks})
	}
}

// CreateTruck handles creating a new truck with photos
// @Summary Create a new truck
// @Description Creates a new truck with optional photo uploads
// @Tags trucks
// @Accept multipart/form-data
// @Produce json
// @Param truck_name formData string true "Truck Name"
// @Param license_plate formData string true "License Plate"
// @Param make formData string false "Make"
// @Param model formData string false "Model"
// @Param year formData int false "Year"
// @Param color formData string false "Color"
// @Param length formData float64 false "Length"
// @Param width formData float64 false "Width"
// @Param height formData float64 false "Height"
// @Param max_weight formData float64 false "Max Weight"
// @Param truck_type formData string false "Truck Type (Small, Medium, Large)"
// @Param climate_control formData boolean false "Climate Control"
// @Param liftgate formData boolean false "Liftgate"
// @Param pallet_jack formData boolean false "Pallet Jack"
// @Param security_system formData boolean false "Security System"
// @Param refrigerated formData boolean false "Refrigerated"
// @Param furniture_pads formData boolean false "Furniture Pads"
// @Param photos formData file false "Truck Photos (multiple allowed)"
// @Success 201 {object} TruckResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks [post]
func (h *TruckHandler) CreateTruck(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, `{"error": "unable to parse form"}`, http.StatusBadRequest)
		return
	}

	// Parse truck data
	truck := &repository.Truck{
		UserID:       userID,
		TruckName:    r.FormValue("truck_name"),
		LicensePlate: r.FormValue("license_plate"),
		Make:         r.FormValue("make"),
		Model:        r.FormValue("model"),
	}
	if yearStr := r.FormValue("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, `{"error": "invalid year"}`, http.StatusBadRequest)
			return
		}
		truck.Year = year
	}
	if color := r.FormValue("color"); color != "" {
		truck.Color = color
	}
	if lengthStr := r.FormValue("length"); lengthStr != "" {
		length, err := strconv.ParseFloat(lengthStr, 64)
		if err != nil {
			http.Error(w, `{"error": "invalid length"}`, http.StatusBadRequest)
			return
		}
		truck.Length = length
	}
	if widthStr := r.FormValue("width"); widthStr != "" {
		width, err := strconv.ParseFloat(widthStr, 64)
		if err != nil {
			http.Error(w, `{"error": "invalid width"}`, http.StatusBadRequest)
			return
		}
		truck.Width = width
	}
	if heightStr := r.FormValue("height"); heightStr != "" {
		height, err := strconv.ParseFloat(heightStr, 64)
		if err != nil {
			http.Error(w, `{"error": "invalid height"}`, http.StatusBadRequest)
			return
		}
		truck.Height = height
	}
	if maxWeightStr := r.FormValue("max_weight"); maxWeightStr != "" {
		maxWeight, err := strconv.ParseFloat(maxWeightStr, 64)
		if err != nil {
			http.Error(w, `{"error": "invalid max_weight"}`, http.StatusBadRequest)
			return
		}
		truck.MaxWeight = maxWeight
	}
	if truckType := r.FormValue("truck_type"); truckType != "" {
		truck.TruckType = truckType
	}
	if climateControlStr := r.FormValue("climate_control"); climateControlStr != "" {
		climateControl, err := strconv.ParseBool(climateControlStr)
		if err != nil {
			http.Error(w, `{"error": "invalid climate_control"}`, http.StatusBadRequest)
			return
		}
		truck.ClimateControl = climateControl
	}
	if liftgateStr := r.FormValue("liftgate"); liftgateStr != "" {
		liftgate, err := strconv.ParseBool(liftgateStr)
		if err != nil {
			http.Error(w, `{"error": "invalid liftgate"}`, http.StatusBadRequest)
			return
		}
		truck.Liftgate = liftgate
	}
	if palletJackStr := r.FormValue("pallet_jack"); palletJackStr != "" {
		palletJack, err := strconv.ParseBool(palletJackStr)
		if err != nil {
			http.Error(w, `{"error": "invalid pallet_jack"}`, http.StatusBadRequest)
			return
		}
		truck.PalletJack = palletJack
	}
	if securitySystemStr := r.FormValue("security_system"); securitySystemStr != "" {
		securitySystem, err := strconv.ParseBool(securitySystemStr)
		if err != nil {
			http.Error(w, `{"error": "invalid security_system"}`, http.StatusBadRequest)
			return
		}
		truck.SecuritySystem = securitySystem
	}
	if refrigeratedStr := r.FormValue("refrigerated"); refrigeratedStr != "" {
		refrigerated, err := strconv.ParseBool(refrigeratedStr)
		if err != nil {
			http.Error(w, `{"error": "invalid refrigerated"}`, http.StatusBadRequest)
			return
		}
		truck.Refrigerated = refrigerated
	}
	if furniturePadsStr := r.FormValue("furniture_pads"); furniturePadsStr != "" {
		furniturePads, err := strconv.ParseBool(furniturePadsStr)
		if err != nil {
			http.Error(w, `{"error": "invalid furniture_pads"}`, http.StatusBadRequest)
			return
		}
		truck.FurniturePads = furniturePads
	}

	// Handle photo uploads
	var files []multipart.FileHeader
	if photoFiles := r.MultipartForm.File["photos"]; len(photoFiles) > 0 {
		files = photoFiles
	}

	err = h.truckService.CreateTruck(r.Context(), userID, truck, files)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := TruckResponse{Truck: truck}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
