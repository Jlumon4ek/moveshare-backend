package handlers

import (
	"encoding/json"
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
	TruckName      string                  `json:"truck_name"`
	LicensePlate   string                  `json:"license_plate"`
	Make           string                  `json:"make"`
	Model          string                  `json:"model"`
	Year           int                     `json:"year"`
	Color          string                  `json:"color"`
	Length         float64                 `json:"length"`
	Width          float64                 `json:"width"`
	Height         float64                 `json:"height"`
	MaxWeight      float64                 `json:"max_weight"`
	TruckType      repository.TruckType    `json:"truck_type"`
	ClimateControl bool                    `json:"climate_control"`
	Liftgate       bool                    `json:"liftgate"`
	PalletJack     bool                    `json:"pallet_jack"`
	SecuritySystem bool                    `json:"security_system"`
	Refrigerated   bool                    `json:"refrigerated"`
	FurniturePads  bool                    `json:"furniture_pads"`
}

// TruckResponse represents a single truck response
type TruckResponse struct {
	Truck *repository.Truck `json:"truck"`
}

// TrucksResponse represents multiple trucks response
type TrucksResponse struct {
	Trucks []repository.Truck `json:"trucks"`
}

// PhotosResponse represents truck photos response
type PhotosResponse struct {
	Photos []repository.TruckPhoto `json:"photos"`
}

// PhotoUploadResponse represents photo upload response
type PhotoUploadResponse struct {
	Message string                    `json:"message"`
	Photos  []repository.TruckPhoto   `json:"photos"`
}

// GetUserTrucks handles fetching all trucks for the current user
// @Summary Get user's trucks
// @Description Fetches all trucks belonging to the current user
// @Tags trucks
// @Produce json
// @Success 200 {object} TrucksResponse
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
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := TrucksResponse{Trucks: trucks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// CreateTruck handles creating a new truck
// @Summary Create a new truck
// @Description Creates a new truck for the current user
// @Tags trucks
// @Accept json
// @Produce json
// @Param body body TruckRequest true "Truck data"
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
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
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

	err := h.truckService.CreateTruck(r.Context(), userID, truck)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	resp := TruckResponse{Truck: truck}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetTruckByID handles fetching a specific truck by ID
// @Summary Get truck by ID
// @Description Fetches a specific truck by ID for the current user
// @Tags trucks
// @Produce json
// @Param id path int true "Truck ID"
// @Success 200 {object} TruckResponse
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
	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid truck ID"}`, http.StatusBadRequest)
		return
	}

	truck, err := h.truckService.GetTruckByID(r.Context(), userID, truckID)
	if err != nil {
		if err.Error() == "truck not found" {
			http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	resp := TruckResponse{Truck: truck}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UpdateTruck handles updating an existing truck
// @Summary Update truck
// @Description Updates an existing truck for the current user
// @Tags trucks
// @Accept json
// @Produce json
// @Param id path int true "Truck ID"
// @Param body body TruckRequest true "Truck data"
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
	truckID, err := strconv.ParseInt(truckIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid truck ID"}`, http.StatusBadRequest)
		return
	}

	var req TruckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
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

	err = h.truckService.UpdateTruck(r.Context(), userID, truckID, truck)
	if err != nil {
		if err.Error() == "truck not found" {
			http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		}
		return
	}

	// Fetch updated truck to return it
	updatedTruck, err := h.truckService.GetTruckByID(r.Context(), userID, truckID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := TruckResponse{Truck: updatedTruck}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteTruck handles deleting a truck
// @Summary Delete truck
// @Description Deletes a truck for the current user
// @Tags trucks
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

	err = h.truckService.DeleteTruck(r.Context(), userID, truckID)
	if err != nil {
		if err.Error() == "truck not found" {
			http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Truck deleted successfully"})
}

// UploadTruckPhotos handles uploading photos for a truck
// @Summary Upload truck photos
// @Description Uploads one or more photos for a specific truck
// @Tags trucks
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Truck ID"
// @Param photos formData file true "Photo files"
// @Success 201 {object} PhotoUploadResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id}/photos [post]
func (h *TruckHandler) UploadTruckPhotos(w http.ResponseWriter, r *http.Request) {
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

	// Parse multipart form (max 50MB total)
	err = r.ParseMultipartForm(50 * 1024 * 1024)
	if err != nil {
		http.Error(w, `{"error": "failed to parse multipart form"}`, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		http.Error(w, `{"error": "no photos provided"}`, http.StatusBadRequest)
		return
	}

	photos, err := h.truckService.UploadTruckPhotos(r.Context(), userID, truckID, files)
	if err != nil {
		if err.Error() == "truck not found" {
			http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		}
		return
	}

	resp := PhotoUploadResponse{
		Message: "Photos uploaded successfully",
		Photos:  photos,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetTruckPhotos handles fetching photos for a truck
// @Summary Get truck photos
// @Description Fetches all photos for a specific truck
// @Tags trucks
// @Produce json
// @Param id path int true "Truck ID"
// @Success 200 {object} PhotosResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id}/photos [get]
func (h *TruckHandler) GetTruckPhotos(w http.ResponseWriter, r *http.Request) {
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

	photos, err := h.truckService.GetTruckPhotos(r.Context(), userID, truckID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if photos == nil {
		http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
		return
	}

	resp := PhotosResponse{Photos: photos}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteTruckPhoto handles deleting a specific truck photo
// @Summary Delete truck photo
// @Description Deletes a specific photo for a truck
// @Tags trucks
// @Produce json
// @Param id path int true "Truck ID"
// @Param photoId path int true "Photo ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trucks/{id}/photos/{photoId} [delete]
func (h *TruckHandler) DeleteTruckPhoto(w http.ResponseWriter, r *http.Request) {
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

	photoIDStr := chi.URLParam(r, "photoId")
	photoID, err := strconv.ParseInt(photoIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid photo ID"}`, http.StatusBadRequest)
		return
	}

	err = h.truckService.DeleteTruckPhoto(r.Context(), userID, truckID, photoID)
	if err != nil {
		if err.Error() == "photo not found" {
			http.Error(w, `{"error": "photo not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Photo deleted successfully"})
}