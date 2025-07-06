package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"moveshare/internal/repository"
	"moveshare/internal/service"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
)

// TruckHandler handles HTTP requests for truck operations
type TruckHandler struct {
	truckService service.TruckService
	minioClient  *minio.Client
	minioBucket  string
}

// NewTruckHandler creates a new TruckHandler
func NewTruckHandler(truckService service.TruckService, minioClient *minio.Client, minioBucket string) *TruckHandler {
	return &TruckHandler{
		truckService: truckService,
		minioClient:  minioClient,
		minioBucket:  minioBucket,
	}
}

// TruckRequest represents the truck creation/update request payload
type TruckRequest struct {
	TruckName      string   `json:"truck_name" example:"Blue Thunder"`
	LicensePlate   string   `json:"license_plate" example:"ABC123"`
	Make           string   `json:"make" example:"Ford"`
	Model          string   `json:"model" example:"F-150"`
	Year           int      `json:"year" example:"2020"`
	Color          string   `json:"color" example:"Red"`
	Length         float64  `json:"length" example:"26"`
	Width          float64  `json:"width" example:"8.5"`
	Height         float64  `json:"height" example:"9.5"`
	MaxWeight      float64  `json:"max_weight" example:"10000"`
	TruckType      string   `json:"truck_type" example:"Large"`
	ClimateControl bool     `json:"climate_control" example:"true"`
	Liftgate       bool     `json:"liftgate" example:"true"`
	PalletJack     bool     `json:"pallet_jack" example:"false"`
	SecuritySystem bool     `json:"security_system" example:"true"`
	Refrigerated   bool     `json:"refrigerated" example:"false"`
	FurniturePads  bool     `json:"furniture_pads" example:"true"`
	Photos         []string `json:"photos,omitempty"`
}

// TruckResponse represents the truck response
type TruckResponse struct {
	Truck *repository.Truck `json:"truck"`
}

// TrucksResponse represents the trucks list response
type TrucksResponse struct {
	Trucks []repository.Truck `json:"trucks"`
}

// CreateTruck godoc
// @Summary      Create a new truck
// @Description  Create a new truck with optional photo uploads. Photos will be stored in MinIO and presigned URLs will be returned.
// @Tags         trucks
// @Accept       multipart/form-data
// @Produce      json
// @Param        truck_name      formData string  true  "Truck Name"
// @Param        license_plate   formData string  true  "License Plate"
// @Param        make            formData string  true  "Make"
// @Param        model           formData string  true  "Model"
// @Param        year            formData int     true  "Year"
// @Param        color           formData string  false "Color"
// @Param        length          formData number  false "Length (ft)"
// @Param        width           formData number  false "Width (ft)"
// @Param        height          formData number  false "Height (ft)"
// @Param        max_weight      formData number  false "Max Weight (lbs)"
// @Param        truck_type      formData string  true  "Truck Type" Enums(Small, Medium, Large)
// @Param        climate_control formData boolean false "Climate Control"
// @Param        liftgate        formData boolean false "Liftgate"
// @Param        pallet_jack     formData boolean false "Pallet Jack"
// @Param        security_system formData boolean false "Security System"
// @Param        refrigerated    formData boolean false "Refrigerated"
// @Param        furniture_pads  formData boolean false "Furniture Pads"
// @Success      201 {object} TruckResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /trucks [post]
func (h *TruckHandler) CreateTruck(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, `{"error": "failed to parse form"}`, http.StatusBadRequest)
		return
	}

	truck := &repository.Truck{}
	truck.TruckName = r.FormValue("truck_name")
	truck.LicensePlate = r.FormValue("license_plate")
	truck.Make = r.FormValue("make")
	truck.Model = r.FormValue("model")
	truck.Color = r.FormValue("color")
	truck.TruckType = r.FormValue("truck_type")

	if yearStr := r.FormValue("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			truck.Year = year
		}
	}
	if lengthStr := r.FormValue("length"); lengthStr != "" {
		if length, err := strconv.ParseFloat(lengthStr, 64); err == nil {
			truck.Length = length
		}
	}
	if widthStr := r.FormValue("width"); widthStr != "" {
		if width, err := strconv.ParseFloat(widthStr, 64); err == nil {
			truck.Width = width
		}
	}
	if heightStr := r.FormValue("height"); heightStr != "" {
		if height, err := strconv.ParseFloat(heightStr, 64); err == nil {
			truck.Height = height
		}
	}
	if maxWeightStr := r.FormValue("max_weight"); maxWeightStr != "" {
		if maxWeight, err := strconv.ParseFloat(maxWeightStr, 64); err == nil {
			truck.MaxWeight = maxWeight
		}
	}
	truck.ClimateControl = r.FormValue("climate_control") == "true"
	truck.Liftgate = r.FormValue("liftgate") == "true"
	truck.PalletJack = r.FormValue("pallet_jack") == "true"
	truck.SecuritySystem = r.FormValue("security_system") == "true"
	truck.Refrigerated = r.FormValue("refrigerated") == "true"
	truck.FurniturePads = r.FormValue("furniture_pads") == "true"

	photos, err := h.handlePhotoUploadsMinio(r.MultipartForm.File["photos"])
	fmt.Printf("Files received: %d\n", len(r.MultipartForm.File["photos"]))

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "failed to upload photos: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	truck.Photos = photos

	err = h.truckService.CreateTruck(r.Context(), userID, truck)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TruckResponse{Truck: truck})
}

// GetUserTrucks godoc
// @Summary      Get user trucks
// @Description  Get all trucks for the authenticated user
// @Tags         trucks
// @Produce      json
// @Success      200 {object} TrucksResponse
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /trucks [get]
func (h *TruckHandler) GetUserTrucks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	trucks, err := h.truckService.GetUserTrucks(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "failed to get trucks"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TrucksResponse{Trucks: trucks})
}

// GetTruckByID godoc
// @Summary      Get truck by ID
// @Description  Get a specific truck by ID for the authenticated user
// @Tags         trucks
// @Produce      json
// @Param        id path int true "Truck ID"
// @Success      200 {object} TruckResponse
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /trucks/{id} [get]
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
	if err != nil || truck == nil {
		http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TruckResponse{Truck: truck})
}

// UpdateTruck godoc
// @Summary      Update truck
// @Description  Update an existing truck
// @Tags         trucks
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path int true "Truck ID"
// @Param        truck_name      formData string  false  "Truck Name"
// @Param        license_plate   formData string  false  "License Plate"
// @Param        make            formData string  false  "Make"
// @Param        model           formData string  false  "Model"
// @Param        year            formData int     false  "Year"
// @Param        color           formData string  false  "Color"
// @Param        length          formData number  false  "Length (ft)"
// @Param        width           formData number  false  "Width (ft)"
// @Param        height          formData number  false  "Height (ft)"
// @Param        max_weight      formData number  false  "Max Weight (lbs)"
// @Param        truck_type      formData string  false  "Truck Type" Enums(Small, Medium, Large)
// @Param        climate_control formData boolean false  "Climate Control"
// @Param        liftgate        formData boolean false  "Liftgate"
// @Param        pallet_jack     formData boolean false  "Pallet Jack"
// @Param        security_system formData boolean false  "Security System"
// @Param        refrigerated    formData boolean false  "Refrigerated"
// @Param        furniture_pads  formData boolean false  "Furniture Pads"
// @Param        photos formData file false "Truck Photos (можно несколько, повторите параметр)"
// @Success      200 {object} TruckResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /trucks/{id} [put]
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

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, `{"error": "failed to parse form"}`, http.StatusBadRequest)
		return
	}

	existingTruck, err := h.truckService.GetTruckByID(r.Context(), userID, truckID)
	if err != nil || existingTruck == nil {
		http.Error(w, `{"error": "truck not found"}`, http.StatusNotFound)
		return
	}

	if truckName := r.FormValue("truck_name"); truckName != "" {
		existingTruck.TruckName = truckName
	}
	if licensePlate := r.FormValue("license_plate"); licensePlate != "" {
		existingTruck.LicensePlate = licensePlate
	}
	if make := r.FormValue("make"); make != "" {
		existingTruck.Make = make
	}
	if model := r.FormValue("model"); model != "" {
		existingTruck.Model = model
	}
	if color := r.FormValue("color"); color != "" {
		existingTruck.Color = color
	}
	if truckType := r.FormValue("truck_type"); truckType != "" {
		existingTruck.TruckType = truckType
	}
	if yearStr := r.FormValue("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			existingTruck.Year = year
		}
	}
	if lengthStr := r.FormValue("length"); lengthStr != "" {
		if length, err := strconv.ParseFloat(lengthStr, 64); err == nil {
			existingTruck.Length = length
		}
	}
	if widthStr := r.FormValue("width"); widthStr != "" {
		if width, err := strconv.ParseFloat(widthStr, 64); err == nil {
			existingTruck.Width = width
		}
	}
	if heightStr := r.FormValue("height"); heightStr != "" {
		if height, err := strconv.ParseFloat(heightStr, 64); err == nil {
			existingTruck.Height = height
		}
	}
	if maxWeightStr := r.FormValue("max_weight"); maxWeightStr != "" {
		if maxWeight, err := strconv.ParseFloat(maxWeightStr, 64); err == nil {
			existingTruck.MaxWeight = maxWeight
		}
	}

	if r.FormValue("climate_control") != "" {
		existingTruck.ClimateControl = r.FormValue("climate_control") == "true"
	}
	if r.FormValue("liftgate") != "" {
		existingTruck.Liftgate = r.FormValue("liftgate") == "true"
	}
	if r.FormValue("pallet_jack") != "" {
		existingTruck.PalletJack = r.FormValue("pallet_jack") == "true"
	}
	if r.FormValue("security_system") != "" {
		existingTruck.SecuritySystem = r.FormValue("security_system") == "true"
	}
	if r.FormValue("refrigerated") != "" {
		existingTruck.Refrigerated = r.FormValue("refrigerated") == "true"
	}
	if r.FormValue("furniture_pads") != "" {
		existingTruck.FurniturePads = r.FormValue("furniture_pads") == "true"
	}

	// Новые фото — добавляются к существующим
	if files := r.MultipartForm.File["photos"]; len(files) > 0 {
		newPhotos, err := h.handlePhotoUploadsMinio(files)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "failed to upload photos: %s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		existingTruck.Photos = append(existingTruck.Photos, newPhotos...)
	}

	err = h.truckService.UpdateTruck(r.Context(), userID, existingTruck)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TruckResponse{Truck: existingTruck})
}

// DeleteTruck godoc
// @Summary      Delete truck
// @Description  Delete a truck
// @Tags         trucks
// @Param        id path int true "Truck ID"
// @Success      204
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /trucks/{id} [delete]
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
		http.Error(w, `{"error": "failed to delete truck"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlePhotoUploadsMinio uploads files to MinIO and returns presigned URLs
func (h *TruckHandler) handlePhotoUploadsMinio(files []*multipart.FileHeader) ([]string, error) {
	var photoURLs []string
	ctx := context.Background()

	for _, fileHeader := range files {
		if !h.isValidImageType(fileHeader.Filename) {
			continue
		}
		src, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		_, err = h.minioClient.PutObject(ctx, h.minioBucket, filename, src, fileHeader.Size, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			continue
		}

		// Сгенерировать presigned URL (временный, например, на 24 часа)
		url, err := h.minioClient.PresignedGetObject(ctx, h.minioBucket, filename, 24*time.Hour, nil)
		if err != nil {
			continue
		}
		photoURLs = append(photoURLs, url.String())
	}

	return photoURLs, nil
}

func (h *TruckHandler) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validExts[ext]
}
