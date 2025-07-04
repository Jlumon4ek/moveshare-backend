package handlers

import (
	"encoding/json"
	"moveshare/internal/repository"
	"moveshare/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5" // Правильный импорт
)

// JobHandler handles HTTP requests for job operations
type JobHandler struct {
	jobService service.JobService
}

// NewJobHandler creates a new JobHandler
func NewJobHandler(jobService service.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

// JobRequest represents the job creation request payload
type JobRequest struct {
	JobTitle           string  `json:"job_title"`
	Description        string  `json:"description"`
	CargoType          string  `json:"cargo_type"`
	Urgency            string  `json:"urgency"`
	TruckSize          string  `json:"truck_size"`
	LoadingAssistance  bool    `json:"loading_assistance"`
	PickupDate         string  `json:"pickup_date"`
	PickupTimeWindow   string  `json:"pickup_time_window"`
	DeliveryDate       string  `json:"delivery_date"`
	DeliveryTimeWindow string  `json:"delivery_time_window"`
	PickupLocation     string  `json:"pickup_location"`
	DeliveryLocation   string  `json:"delivery_location"`
	PayoutAmount       float64 `json:"payout_amount"`
	EarlyDeliveryBonus float64 `json:"early_delivery_bonus"`
	PaymentTerms       string  `json:"payment_terms"`
	WeightLb           float64 `json:"weight_lb"`
	VolumeCuFt         float64 `json:"volume_cu_ft"`
	Liftgate           bool    `json:"liftgate"`
	FragileItems       bool    `json:"fragile_items"`
	ClimateControl     bool    `json:"climate_control"`
	AssemblyRequired   bool    `json:"assembly_required"`
	ExtraInsurance     bool    `json:"extra_insurance"`
	AdditionalPacking  bool    `json:"additional_packing"`
}

// JobResponse represents the job creation response
type JobResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// AvailableJobsResponse represents the response for available jobs
type AvailableJobsResponse struct {
	Jobs []repository.Job `json:"jobs"`
}

// MyJobsResponse represents the response for user jobs
type MyJobsResponse struct {
	Jobs []repository.Job `json:"jobs"`
}

// ApplicationResponse represents the application response
type ApplicationResponse struct {
	Message string `json:"message"`
}

// CreateJob handles job creation
// @Summary Create a new job
// @Description Creates a new job posting
// @Tags jobs
// @Accept json
// @Produce json
// @Param body body JobRequest true "Job creation data"
// @Success 201 {object} JobResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs [post]
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Parse dates
	pickupDate, err := time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		http.Error(w, `{"error": "invalid pickup_date format, use YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}
	deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
	if err != nil {
		http.Error(w, `{"error": "invalid delivery_date format, use YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	job := &repository.Job{
		JobTitle:           req.JobTitle,
		Description:        req.Description,
		CargoType:          req.CargoType,
		Urgency:            req.Urgency,
		TruckSize:          req.TruckSize,
		LoadingAssistance:  req.LoadingAssistance,
		PickupDate:         pickupDate,
		PickupTimeWindow:   req.PickupTimeWindow,
		DeliveryDate:       deliveryDate,
		DeliveryTimeWindow: req.DeliveryTimeWindow,
		PickupLocation:     req.PickupLocation,
		DeliveryLocation:   req.DeliveryLocation,
		PayoutAmount:       req.PayoutAmount,
		EarlyDeliveryBonus: req.EarlyDeliveryBonus,
		PaymentTerms:       req.PaymentTerms,
		WeightLb:           req.WeightLb,
		VolumeCuFt:         req.VolumeCuFt,
		Liftgate:           req.Liftgate,
		FragileItems:       req.FragileItems,
		ClimateControl:     req.ClimateControl,
		AssemblyRequired:   req.AssemblyRequired,
		ExtraInsurance:     req.ExtraInsurance,
		AdditionalPacking:  req.AdditionalPacking,
	}

	err = h.jobService.CreateJob(r.Context(), userID, job)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	resp := JobResponse{ID: job.ID, Message: "Job created successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetAvailableJobs handles fetching available jobs
// @Summary Get available jobs
// @Description Fetches jobs not created by the current user with optional filters and pagination
// @Tags jobs
// @Produce json
// @Param pickup_location query string false "Pickup location filter"
// @Param delivery_location query string false "Delivery location filter"
// @Param pickup_date_start query string false "Start date for pickup (YYYY-MM-DD)"
// @Param pickup_date_end query string false "End date for pickup (YYYY-MM-DD)"
// @Param truck_size query string false "Truck size filter"
// @Param limit query int false "Number of results per page (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} AvailableJobsResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/available [get]
func (h *JobHandler) GetAvailableJobs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	filters := make(map[string]string)
	values := r.URL.Query()
	if val := values.Get("pickup_location"); val != "" {
		filters["pickup_location"] = val
	}
	if val := values.Get("delivery_location"); val != "" {
		filters["delivery_location"] = val
	}
	if val := values.Get("pickup_date_start"); val != "" {
		filters["pickup_date_start"] = val
	}
	if val := values.Get("pickup_date_end"); val != "" {
		filters["pickup_date_end"] = val
	}
	if val := values.Get("truck_size"); val != "" {
		filters["truck_size"] = val
	}

	// Parse pagination parameters
	limit := 10
	if val := values.Get("limit"); val != "" {
		if l, err := strconv.Atoi(val); err == nil && l > 0 {
			limit = l
		}
	}
	offset := 0
	if val := values.Get("offset"); val != "" {
		if o, err := strconv.Atoi(val); err == nil && o >= 0 {
			offset = o
		}
	}

	jobs, err := h.jobService.GetAvailableJobs(r.Context(), userID, filters, limit, offset)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := AvailableJobsResponse{Jobs: jobs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetUserJobs handles fetching jobs created by the current user
// @Summary Get user jobs
// @Description Fetches all jobs created by the current user
// @Tags jobs
// @Produce json
// @Success 200 {object} MyJobsResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/my [get]
func (h *JobHandler) GetUserJobs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	jobs, err := h.jobService.GetUserJobs(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := MyJobsResponse{Jobs: jobs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteJob handles job deletion
// @Summary Delete a job
// @Description Deletes a job created by the current user
// @Tags jobs
// @Produce json
// @Param id path int true "Job ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id} [delete]
func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	jobIDStr := chi.URLParam(r, "id")
	if jobIDStr == "" {
		http.Error(w, `{"error": "job ID is required"}`, http.StatusBadRequest)
		return
	}

	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid job ID"}`, http.StatusBadRequest)
		return
	}

	err = h.jobService.DeleteJob(r.Context(), userID, jobID)
	if err != nil {
		if err.Error() == "job not found or unauthorized" {
			http.Error(w, `{"error": "job not found or unauthorized"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Job deleted successfully"})
}

// ApplyForJob handles job application
// @Summary Apply for a job
// @Description Allows a user to apply for a job
// @Tags jobs
// @Param id path int true "Job ID"
// @Success 201 {object} ApplicationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id}/apply [post]
func (h *JobHandler) ApplyForJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	jobIDStr := chi.URLParam(r, "id")
	if jobIDStr == "" {
		http.Error(w, `{"error": "job ID is required"}`, http.StatusBadRequest)
		return
	}

	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid job ID"}`, http.StatusBadRequest)
		return
	}

	err = h.jobService.ApplyForJob(r.Context(), userID, jobID)
	if err != nil {
		if err.Error() == "cannot apply to your own job" {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusForbidden)
		} else if err.Error() == "application already exists or job not found" {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusConflict)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	resp := ApplicationResponse{Message: "Application submitted successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetMyApplications handles fetching applications submitted by the current user
// @Summary Get my applications
// @Description Fetches all applications submitted by the current user
// @Tags jobs
// @Produce json
// @Success 200 {object} MyJobsResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/applications/my [get]
func (h *JobHandler) GetMyApplications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// This is a placeholder; you'll need to implement the repository and service logic to fetch applications
	// For now, assuming a new method in jobRepo to get applications by userID
	jobs, err := h.jobService.GetMyApplications(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := MyJobsResponse{Jobs: jobs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
