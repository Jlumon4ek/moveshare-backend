package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"moveshare/internal/repository"
	"moveshare/internal/service"

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

// TruckResponse represents the truck response payload
type TruckResponse struct {
	*repository.Truck
}

// TrucksResponse represents multiple trucks response
type TrucksResponse struct {
	Trucks []repository.Truck `json:"trucks"`
}

// TruckPhotosResponse represents truck photos response
type TruckPhotosResponse struct {
	Photos []repository.TruckPhoto `json:"photos"`
}

// PhotoUploadResponse represents photo upload response
type PhotoUploadResponse struct {
	Message     string `json:"message"`
	PhotosCount int    `json:"photos_count"`
}

// GetMyTrucks handles fetching user's trucks
// @Summary Get user's trucks
// @Description Fetches all trucks owned by the current user
// @Tags trucks
// @Accept json
// @Produce json
// @Success 200 {object} TrucksResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/my [get]
func (h *TruckHandler) GetMyTrucks(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(w, r)
	if !ok {
		return
	}
	
	trucks, err := h.truckService.GetUserTrucks(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch trucks: %v"}`, err), http.StatusInternalServerError)
		return
	}
	
	response := TrucksResponse{Trucks: trucks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
	userID, ok := getUserIDFromContext(w, r)
	if !ok {
		return
	}
	
	var req TruckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
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
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TruckResponse{Truck: truck})
}

// UploadTruckPhotos handles truck photo uploads
// @Summary Upload truck photos
// @Description Uploads photos for a specific truck (max 10 photos per truck)
// @Tags trucks
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Truck ID"
// @Param photos formData file true "Photo files (max 10 total)"
// @Success 200 {object} PhotoUploadResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id}/photos [post]
func (h *TruckHandler) UploadTruckPhotos(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(w, r)
	if !ok {
		return
	}
	
	truckIDStr := chi.URLParam(r, "id")
	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid truck ID"}`, http.StatusBadRequest)
		return
	}
	
	// Parse multipart form (32MB max)
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, `{"error": "Failed to parse multipart form"}`, http.StatusBadRequest)
		return
	}
	
	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		http.Error(w, `{"error": "No photos provided"}`, http.StatusBadRequest)
		return
	}
	
	// Create upload directory if it doesn't exist
	uploadDir := fmt.Sprintf("uploads/trucks/%d", truckID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, `{"error": "Failed to create upload directory"}`, http.StatusInternalServerError)
		return
	}
	
	// Save files to disk
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to open file %s"}`, fileHeader.Filename), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		
		// Create destination file
		dst, err := os.Create(filepath.Join(uploadDir, fileHeader.Filename))
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to create file %s"}`, fileHeader.Filename), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		
		// Copy file content
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to save file %s"}`, fileHeader.Filename), http.StatusInternalServerError)
			return
		}
	}
	
	// Upload to service
	if err := h.truckService.UploadTruckPhotos(r.Context(), userID, truckID, files); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}
	
	response := PhotoUploadResponse{
		Message:     fmt.Sprintf("Successfully uploaded %d photos", len(files)),
		PhotosCount: len(files),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTruckPhotos handles fetching truck photos
// @Summary Get truck photos
// @Description Fetches all photos for a specific truck
// @Tags trucks
// @Accept json
// @Produce json
// @Param id path int true "Truck ID"
// @Success 200 {object} TruckPhotosResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id}/photos [get]
func (h *TruckHandler) GetTruckPhotos(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(w, r)
	if !ok {
		return
	}
	
	truckIDStr := chi.URLParam(r, "id")
	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid truck ID"}`, http.StatusBadRequest)
		return
	}
	
	photos, err := h.truckService.GetTruckPhotos(r.Context(), userID, truckID)
	if err != nil {
		if err.Error() == "truck not found or access denied" {
			http.Error(w, `{"error": "Truck not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch photos: %v"}`, err), http.StatusInternalServerError)
		return
	}
	
	response := TruckPhotosResponse{Photos: photos}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to get user ID from context with error handling
func getUserIDFromContext(w http.ResponseWriter, r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return 0, false
	}
	return userID, true
}