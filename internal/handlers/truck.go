package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"moveshare/internal/repository"
	"moveshare/internal/service"
)

// TruckHandler handles HTTP requests for truck operations
type TruckHandler struct {
	truckService service.TruckService
}

// NewTruckHandler creates a new TruckHandler
func NewTruckHandler(truckService service.TruckService) *TruckHandler {
	return &TruckHandler{truckService: truckService}
}

// TruckRequest represents the truck creation/update request payload
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

// TruckResponse represents the truck response
type TruckResponse struct {
	Truck repository.Truck `json:"truck"`
}

// TrucksResponse represents the trucks list response
type TrucksResponse struct {
	Trucks []repository.Truck `json:"trucks"`
}

// CreateTruck handles truck creation with photo upload
// @Summary Create a new truck
// @Description Creates a new truck with optional photo uploads
// @Tags trucks
// @Accept multipart/form-data
// @Produce json
// @Param truck_name formData string true "Truck name"
// @Param license_plate formData string true "License plate"
// @Param make formData string false "Make"
// @Param model formData string false "Model"
// @Param year formData int false "Year"
// @Param color formData string false "Color"
// @Param length formData number false "Length"
// @Param width formData number false "Width"
// @Param height formData number false "Height"
// @Param max_weight formData number false "Max weight"
// @Param truck_type formData string false "Truck type (Small, Medium, Large)"
// @Param climate_control formData boolean false "Climate control"
// @Param liftgate formData boolean false "Liftgate"
// @Param pallet_jack formData boolean false "Pallet jack"
// @Param security_system formData boolean false "Security system"
// @Param refrigerated formData boolean false "Refrigerated"
// @Param furniture_pads formData boolean false "Furniture pads"
// @Param photos formData file false "Photo files"
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

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, `{"error": "failed to parse form"}`, http.StatusBadRequest)
		return
	}

	// Extract form fields
	truckName := r.FormValue("truck_name")
	licensePlate := r.FormValue("license_plate")
	
	if truckName == "" || licensePlate == "" {
		http.Error(w, `{"error": "truck_name and license_plate are required"}`, http.StatusBadRequest)
		return
	}

	// Parse numeric fields
	year, _ := strconv.Atoi(r.FormValue("year"))
	length, _ := strconv.ParseFloat(r.FormValue("length"), 64)
	width, _ := strconv.ParseFloat(r.FormValue("width"), 64)
	height, _ := strconv.ParseFloat(r.FormValue("height"), 64)
	maxWeight, _ := strconv.ParseFloat(r.FormValue("max_weight"), 64)

	// Parse boolean fields
	climateControl := r.FormValue("climate_control") == "true"
	liftgate := r.FormValue("liftgate") == "true"
	palletJack := r.FormValue("pallet_jack") == "true"
	securitySystem := r.FormValue("security_system") == "true"
	refrigerated := r.FormValue("refrigerated") == "true"
	furniturePads := r.FormValue("furniture_pads") == "true"

	truck := &repository.Truck{
		UserID:         userID,
		TruckName:      truckName,
		LicensePlate:   licensePlate,
		Make:           r.FormValue("make"),
		Model:          r.FormValue("model"),
		Year:           year,
		Color:          r.FormValue("color"),
		Length:         length,
		Width:          width,
		Height:         height,
		MaxWeight:      maxWeight,
		TruckType:      r.FormValue("truck_type"),
		ClimateControl: climateControl,
		Liftgate:       liftgate,
		PalletJack:     palletJack,
		SecuritySystem: securitySystem,
		Refrigerated:   refrigerated,
		FurniturePads:  furniturePads,
	}

	// Extract photo files
	var photos []*multipart.FileHeader
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if files, ok := r.MultipartForm.File["photos"]; ok {
			photos = files
		}
	}

	err = h.truckService.CreateTruck(r.Context(), truck, photos)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := TruckResponse{Truck: *truck}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetUserTrucks handles fetching trucks owned by the current user
// @Summary Get user trucks
// @Description Fetches all trucks owned by the current user
// @Tags trucks
// @Produce json
// @Success 200 {object} TrucksResponse
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

	resp := TrucksResponse{Trucks: trucks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetTruckByID handles fetching a specific truck by ID
// @Summary Get truck by ID
// @Description Fetches a specific truck by ID owned by the current user
// @Tags trucks
// @Produce json
// @Param id path int true "Truck ID"
// @Success 200 {object} TruckResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id} [get]
func (h *TruckHandler) GetTruckByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	truckIDStr := chi.URLParam(r, "id")
	if truckIDStr == "" {
		http.Error(w, `{"error": "truck ID is required"}`, http.StatusBadRequest)
		return
	}

	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid truck ID"}`, http.StatusBadRequest)
		return
	}

	truck, err := h.truckService.GetTruckByID(r.Context(), userID, truckID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	resp := TruckResponse{Truck: *truck}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateTruck handles truck updates
// @Summary Update a truck
// @Description Updates an existing truck owned by the current user
// @Tags trucks
// @Accept json
// @Produce json
// @Param id path int true "Truck ID"
// @Param body body TruckRequest true "Truck update data"
// @Success 200 {object} TruckResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id} [put]
func (h *TruckHandler) UpdateTruck(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	truckIDStr := chi.URLParam(r, "id")
	if truckIDStr == "" {
		http.Error(w, `{"error": "truck ID is required"}`, http.StatusBadRequest)
		return
	}

	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid truck ID"}`, http.StatusBadRequest)
		return
	}

	var req TruckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.TruckName == "" || req.LicensePlate == "" {
		http.Error(w, `{"error": "truck_name and license_plate are required"}`, http.StatusBadRequest)
		return
	}

	truck := &repository.Truck{
		ID:             truckID,
		UserID:         userID,
		TruckName:      req.TruckName,
		LicensePlate:   req.LicensePlate,
		Make:           req.Make,
		Model:          req.Model,
		Year:           req.Year,
		Color:          req.Color,
		Length:         req.Length,
		Width:          req.Width,
		Height:         req.Height,
		MaxWeight:      req.MaxWeight,
		TruckType:      req.TruckType,
		ClimateControl: req.ClimateControl,
		Liftgate:       req.Liftgate,
		PalletJack:     req.PalletJack,
		SecuritySystem: req.SecuritySystem,
		Refrigerated:   req.Refrigerated,
		FurniturePads:  req.FurniturePads,
	}

	err = h.truckService.UpdateTruck(r.Context(), truck)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Fetch updated truck to return
	updatedTruck, err := h.truckService.GetTruckByID(r.Context(), userID, truckID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := TruckResponse{Truck: *updatedTruck}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteTruck handles truck deletion
// @Summary Delete a truck
// @Description Deletes a truck owned by the current user
// @Tags trucks
// @Produce json
// @Param id path int true "Truck ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id} [delete]
func (h *TruckHandler) DeleteTruck(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	truckIDStr := chi.URLParam(r, "id")
	if truckIDStr == "" {
		http.Error(w, `{"error": "truck ID is required"}`, http.StatusBadRequest)
		return
	}

	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid truck ID"}`, http.StatusBadRequest)
		return
	}

	err = h.truckService.DeleteTruck(r.Context(), userID, truckID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": "Truck deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}