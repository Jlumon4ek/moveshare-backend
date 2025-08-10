// internal/handlers/job/create_job.go
package job

// import (
// 	"encoding/json"
// 	"fmt"
// 	"moveshare/internal/models"
// 	"moveshare/internal/service"
// 	"moveshare/internal/utils"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// // CreateJob godoc
// // @Summary      Create a new job
// // @Description  Creates a new job posting with images and detailed information
// // @Tags         Jobs
// // @Security     BearerAuth
// // @Accept       multipart/form-data
// // @Produce      json
// // @Param        job_type               formData string   true  "Job Type" Enums(Residential Move, Office Relocation, Warehouse Transfer, Other)
// // @Param        job_type_other         formData string   false "Custom job type if Other is selected"
// // @Param        number_of_bedrooms     formData string   true  "Number of bedrooms" Enums(1 Bedroom, 2 Bedrooms, 3 Bedrooms, 4 Bedrooms, 5+ Bedrooms, Office)
// // @Param        additional_services    formData string   false "Additional services as JSON array"
// // @Param        additional_services_description formData string false "Description of additional services"
// // @Param        truck_size             formData string   true  "Truck size" Enums(Small, Medium, Large)
// // @Param        pickup_location        formData string   true  "Pickup address"
// // @Param        pickup_location_type   formData string   true  "Pickup location type" Enums(House, Stairs, Elevator)
// // @Param        pickup_floor           formData int      false "Floor number if Stairs selected"
// // @Param        pickup_walk_distance   formData string   true  "Walking distance at pickup"
// // @Param        pickup_date            formData string   true  "Pickup date (YYYY-MM-DD)"
// // @Param        pickup_time_window     formData string   true  "Pickup time window (HH:MM-HH:MM)"
// // @Param        delivery_location      formData string   true  "Delivery address"
// // @Param        delivery_location_type formData string   true  "Delivery location type" Enums(House, Stairs, Elevator)
// // @Param        delivery_floor         formData int      false "Floor number if Stairs selected"
// // @Param        delivery_walk_distance formData string   true  "Walking distance at delivery"
// // @Param        delivery_date          formData string   true  "Delivery date (YYYY-MM-DD)"
// // @Param        delivery_time_window   formData string   true  "Delivery time window (HH:MM-HH:MM)"
// // @Param        payment_amount         formData number   true  "Payment amount in USD"
// // @Param        images                 formData []file   false "Job images or inventory PDFs"
// // @Success      201  {object}  models.CreateJobResponse
// // @Failure      400  {object}  map[string]string
// // @Failure      401  {object}  map[string]string
// // @Failure      500  {object}  map[string]string
// // @Router       /jobs [post]
// func CreateJob(jobService service.JobService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID, err := utils.GetUserIDFromContext(c)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 			return
// 		}

// 		// Парсим multipart form
// 		form, err := c.MultipartForm()
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
// 			return
// 		}

// 		// Получаем файлы
// 		files := form.File["images"]

// 		// Парсим основные поля
// 		var req models.CreateJobRequest
// 		if err := c.ShouldBind(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error":   "Invalid request data",
// 				"details": err.Error(),
// 			})
// 			return
// 		}

// 		// Валидируем данные
// 		if err := validateJobRequest(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Парсим дополнительные услуги (JSON array)
// 		var additionalServices []string
// 		if req.AdditionalServices != "" {
// 			if err := json.Unmarshal([]byte(req.AdditionalServices), &additionalServices); err != nil {
// 				// Если не JSON, пробуем разделить по запятым
// 				additionalServices = strings.Split(req.AdditionalServices, ",")
// 				for i := range additionalServices {
// 					additionalServices[i] = strings.TrimSpace(additionalServices[i])
// 				}
// 			}
// 		}

// 		// Парсим даты
// 		pickupDate, err := time.Parse("2006-01-02", req.PickupDate)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pickup date format. Use YYYY-MM-DD"})
// 			return
// 		}

// 		deliveryDate, err := time.Parse("2006-01-02", req.DeliveryDate)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery date format. Use YYYY-MM-DD"})
// 			return
// 		}

// 		// Создаем объект Job
// 		job := &models.Job{
// 			UserID:                 userID,
// 			JobType:                req.JobType,
// 			JobTypeOther:           req.JobTypeOther,
// 			NumberOfBedrooms:       req.NumberOfBedrooms,
// 			AdditionalServices:     additionalServices,
// 			AdditionalServicesDesc: req.AdditionalServicesDesc,
// 			TruckSize:              req.TruckSize,
// 			PickupLocation:         req.PickupLocation,
// 			PickupLocationType:     req.PickupLocationType,
// 			PickupFloor:            req.PickupFloor,
// 			PickupWalkDistance:     req.PickupWalkDistance,
// 			PickupDate:             pickupDate,
// 			PickupTimeWindow:       req.PickupTimeWindow,
// 			DeliveryLocation:       req.DeliveryLocation,
// 			DeliveryLocationType:   req.DeliveryLocationType,
// 			DeliveryFloor:          req.DeliveryFloor,
// 			DeliveryWalkDistance:   req.DeliveryWalkDistance,
// 			DeliveryDate:           deliveryDate,
// 			DeliveryTimeWindow:     req.DeliveryTimeWindow,
// 			PaymentAmount:          req.PaymentAmount,
// 			Images:                 files,
// 			Status:                 "available",
// 		}

// 		// Создаем работу
// 		jobID, err := jobService.CreateJob(c.Request.Context(), job)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error":   "Failed to create job",
// 				"details": err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusCreated, models.CreateJobResponse{
// 			JobID:   jobID,
// 			Message: "Job created successfully",
// 			Success: true,
// 		})
// 	}
// }

// func validateJobRequest(req *models.CreateJobRequest) error {
// 	// Валидация Job Type
// 	if !contains(models.ValidJobTypes, req.JobType) {
// 		return fmt.Errorf("invalid job type. Must be one of: %v", models.ValidJobTypes)
// 	}

// 	// Если выбран Other, должен быть указан job_type_other
// 	if req.JobType == "Other" && strings.TrimSpace(req.JobTypeOther) == "" {
// 		return fmt.Errorf("job_type_other is required when job_type is 'Other'")
// 	}

// 	// Валидация Number of Bedrooms
// 	if !contains(models.ValidBedroomOptions, req.NumberOfBedrooms) {
// 		return fmt.Errorf("invalid number of bedrooms. Must be one of: %v", models.ValidBedroomOptions)
// 	}

// 	// Валидация Truck Size
// 	if !contains(models.ValidTruckSizes, req.TruckSize) {
// 		return fmt.Errorf("invalid truck size. Must be one of: %v", models.ValidTruckSizes)
// 	}

// 	// Валидация Location Types
// 	if !contains(models.ValidLocationTypes, req.PickupLocationType) {
// 		return fmt.Errorf("invalid pickup location type. Must be one of: %v", models.ValidLocationTypes)
// 	}

// 	if !contains(models.ValidLocationTypes, req.DeliveryLocationType) {
// 		return fmt.Errorf("invalid delivery location type. Must be one of: %v", models.ValidLocationTypes)
// 	}

// 	// Валидация Floor - должен быть указан если выбран Stairs
// 	if req.PickupLocationType == "Stairs" && req.PickupFloor == nil {
// 		return fmt.Errorf("pickup_floor is required when pickup_location_type is 'Stairs'")
// 	}

// 	if req.DeliveryLocationType == "Stairs" && req.DeliveryFloor == nil {
// 		return fmt.Errorf("delivery_floor is required when delivery_location_type is 'Stairs'")
// 	}

// 	// Валидация дат - delivery должен быть после или равен pickup
// 	pickupDate, err1 := time.Parse("2006-01-02", req.PickupDate)
// 	deliveryDate, err2 := time.Parse("2006-01-02", req.DeliveryDate)

// 	if err1 != nil || err2 != nil {
// 		return fmt.Errorf("invalid date format. Use YYYY-MM-DD")
// 	}

// 	if deliveryDate.Before(pickupDate) {
// 		return fmt.Errorf("delivery date cannot be before pickup date")
// 	}

// 	// Валидация времени
// 	if !isValidTimeWindow(req.PickupTimeWindow) {
// 		return fmt.Errorf("invalid pickup time window format. Use HH:MM-HH:MM")
// 	}

// 	if !isValidTimeWindow(req.DeliveryTimeWindow) {
// 		return fmt.Errorf("invalid delivery time window format. Use HH:MM-HH:MM")
// 	}

// 	return nil
// }

// func contains(slice []string, item string) bool {
// 	for _, s := range slice {
// 		if s == item {
// 			return true
// 		}
// 	}
// 	return false
// }

// func isValidTimeWindow(timeWindow string) bool {
// 	parts := strings.Split(timeWindow, "-")
// 	if len(parts) != 2 {
// 		return false
// 	}

// 	// Проверяем формат времени HH:MM
// 	for _, part := range parts {
// 		if _, err := time.Parse("15:04", strings.TrimSpace(part)); err != nil {
// 			return false
// 		}
// 	}

// 	return true
// }
