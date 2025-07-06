package handlers

import (
	"encoding/json"
	"moveshare/internal/repository"
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// UserTrucksResponse represents the response for user trucks
type UserTrucksResponse struct {
	Trucks []repository.Truck `json:"trucks"`
}

// TruckPhotosResponse represents the response for truck photos
type TruckPhotosResponse struct {
	Photos []repository.TruckPhoto `json:"photos"`
}

// PhotoUploadResponse represents the response for photo upload
type PhotoUploadResponse struct {
	ID      int64  `json:"id"`
	FileURL string `json:"file_url"`
	Message string `json:"message"`
}

// CreateTruck handles truck creation
// @Summary Create a new truck
// @Description Creates a new truck for the authenticated user
// @Tags trucks
// @Accept json
// @Produce json
// @Param body body TruckRequest true "Truck creation data"
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

	var req TruckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	truck := &repository.Truck{
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

	if err := h.truckService.CreateTruck(r.Context(), userID, truck); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	response := TruckResponse{
		ID:      truck.ID,
		Message: "Truck created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetUserTrucks handles fetching trucks for the current user
// @Summary Get user trucks
// @Description Fetches all trucks owned by the authenticated user
// @Tags trucks
// @Accept json
// @Produce json
// @Success 200 {object} UserTrucksResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/my [get]
func (h *TruckHandler) GetUserTrucks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	trucks, err := h.truckService.GetUserTrucks(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch trucks"}`, http.StatusInternalServerError)
		return
	}

	response := UserTrucksResponse{
		Trucks: trucks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteTruck handles truck deletion
// @Summary Delete a truck
// @Description Deletes a truck owned by the authenticated user
// @Tags trucks
// @Accept json
// @Produce json
// @Param id path int true "Truck ID"
// @Success 200 {object} map[string]string
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
	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid truck ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.truckService.DeleteTruck(r.Context(), userID, truckID); err != nil {
		if err == repository.ErrTruckNotFound {
			http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "failed to delete truck"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Truck deleted successfully"})
}

// UploadTruckPhoto handles photo upload for a truck
// @Summary Upload truck photo
// @Description Uploads a photo for a specific truck
// @Tags trucks
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Truck ID"
// @Param photo formData file true "Photo file"
// @Success 201 {object} PhotoUploadResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id}/photos [post]
func (h *TruckHandler) UploadTruckPhoto(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	truckIDStr := chi.URLParam(r, "id")
	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid truck ID"}`, http.StatusBadRequest)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		http.Error(w, `{"error": "failed to parse form"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, `{"error": "no photo file provided"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	photo, err := h.truckService.UploadTruckPhoto(r.Context(), userID, truckID, file, header)
	if err != nil {
		if err.Error() == "truck not found" || err.Error() == "unauthorized access to truck" {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	response := PhotoUploadResponse{
		ID:      photo.ID,
		FileURL: photo.FileURL,
		Message: "Photo uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetTruckPhotos handles fetching photos for a truck
// @Summary Get truck photos
// @Description Fetches all photos for a specific truck
// @Tags trucks
// @Accept json
// @Produce json
// @Param id path int true "Truck ID"
// @Success 200 {object} TruckPhotosResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id}/photos [get]
func (h *TruckHandler) GetTruckPhotos(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	truckIDStr := chi.URLParam(r, "id")
	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid truck ID"}`, http.StatusBadRequest)
		return
	}

	photos, err := h.truckService.GetTruckPhotos(r.Context(), truckID)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch photos"}`, http.StatusInternalServerError)
		return
	}

	response := TruckPhotosResponse{
		Photos: photos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}