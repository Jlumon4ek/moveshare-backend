package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

// POST /post-new-job
func (h *JobHandler) PostNewJob(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.jobService.CreateJob(userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job created successfully",
		"job":     job,
	})
}

// POST /claim-job/:id
func (h *JobHandler) ClaimJob(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	err = h.jobService.ClaimJob(jobID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job claimed successfully"})
}

// DELETE /delete-job/:id
func (h *JobHandler) DeleteJob(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	err = h.jobService.DeleteJob(jobID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}

// GET /job/:id
func (h *JobHandler) GetJobByID(c *gin.Context) {
	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.jobService.GetJobByID(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"job": job})
}

func (h *JobHandler) GetMyJobs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var pagination models.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobs, total, err := h.jobService.GetMyJobs(userID.(int64), pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + pagination.Limit - 1) / pagination.Limit

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
		"pagination": gin.H{
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *JobHandler) GetAvailableJobs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var filters models.JobFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	jobs, total, err := h.jobService.GetAvailableJobs(userID.(int64), &filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get available jobs",
			"details": err.Error(),
		})
		return
	}

	totalPages := (total + filters.Limit - 1) / filters.Limit

	response := gin.H{
		"jobs": jobs,
		"pagination": gin.H{
			"page":        filters.Page,
			"limit":       filters.Limit,
			"total":       total,
			"total_pages": totalPages,
		},
		"filters_applied": gin.H{
			"number_of_bedrooms": filters.NumberOfBedrooms,
			"origin":             filters.Origin,
			"destination":        filters.Destination,
			"max_distance":       filters.MaxDistance,
			"date_start":         filters.DateStart,
			"date_end":           filters.DateEnd,
			"truck_size":         filters.TruckSize,
			"payout_min":         filters.PayoutMin,
			"payout_max":         filters.PayoutMax,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetJobFilterOptions godoc
// @Summary Get filter options
// @Description Gets available filter options for jobs (unique values)
// @Tags Jobs
// @Security BearerAuth
// @Success 200 {object} models.JobFilterOptions
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/filter-options/ [get]
func (h *JobHandler) GetJobFilterOptions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	options, err := h.jobService.GetFilterOptions(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get filter options",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, options)
}

func (h *JobHandler) GetClaimedJobs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var pagination models.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobs, total, err := h.jobService.GetClaimedJobs(userID.(int64), pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + pagination.Limit - 1) / pagination.Limit

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
		"pagination": gin.H{
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// POST /mark-job-completed/:id
func (h *JobHandler) MarkJobCompleted(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	err = h.jobService.MarkJobCompleted(jobID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job marked as completed successfully"})
}

// POST /export-jobs/
func (h *JobHandler) ExportJobs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.ExportJobsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobs, err := h.jobService.GetJobsForExport(userID.(int64), req.JobIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(jobs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No jobs found with provided IDs"})
		return
	}

	csvData, err := h.generateCSV(jobs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSV"})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=jobs_export.csv")
	c.Header("Content-Length", fmt.Sprintf("%d", len(csvData)))
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "application/octet-stream", csvData)
}

func (h *JobHandler) generateCSV(jobs []models.Job) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{
		"job_type", "number_of_bedrooms", "packing_boxes", "bulky_items", "inventory_list", 
		"hoisting", "additional_services_description", "estimated_crew_assistants", "truck_size", 
		"pickup_address", "pickup_floor", "pickup_building_type", "pickup_walk_distance", 
		"delivery_address", "delivery_floor", "delivery_building_type", "delivery_walk_distance", 
		"distance_miles", "job_status", "pickup_date", "pickup_time_from", "pickup_time_to", 
		"delivery_date", "delivery_time_from", "delivery_time_to", "cut_amount", "payment_amount", 
		"weight_lbs", "volume_cu_ft",
	}

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, job := range jobs {
		record := []string{
			job.JobType,
			job.NumberOfBedrooms,
			fmt.Sprintf("%t", job.PackingBoxes),
			fmt.Sprintf("%t", job.BulkyItems),
			fmt.Sprintf("%t", job.InventoryList),
			fmt.Sprintf("%t", job.Hoisting),
			job.AdditionalServicesDescription,
			job.EstimatedCrewAssistants,
			job.TruckSize,
			job.PickupAddress,
			h.formatIntPtr(job.PickupFloor),
			job.PickupBuildingType,
			job.PickupWalkDistance,
			job.DeliveryAddress,
			h.formatIntPtr(job.DeliveryFloor),
			job.DeliveryBuildingType,
			job.DeliveryWalkDistance,
			fmt.Sprintf("%.2f", job.DistanceMiles),
			job.JobStatus,
			job.PickupDate.Format("2006-01-02"),
			job.PickupTimeFrom.Format("15:04"),
			job.PickupTimeTo.Format("15:04"),
			job.DeliveryDate.Format("2006-01-02"),
			job.DeliveryTimeFrom.Format("15:04"),
			job.DeliveryTimeTo.Format("15:04"),
			fmt.Sprintf("%.2f", job.CutAmount),
			fmt.Sprintf("%.2f", job.PaymentAmount),
			fmt.Sprintf("%.2f", job.WeightLbs),
			fmt.Sprintf("%.2f", job.VolumeCuFt),
		}

		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *JobHandler) formatIntPtr(ptr *int) string {
	if ptr == nil {
		return ""
	}
	return fmt.Sprintf("%d", *ptr)
}

// GET /jobs/stats
func (h *JobHandler) GetJobsStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	stats, err := h.jobService.GetJobsStats(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get jobs statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
