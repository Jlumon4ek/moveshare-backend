package truck

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTruck godoc
// @Summary      Create a new truck
// @Description  Creates a truck associated with the authenticated user and uploads photos.
// @Tags         Trucks
// @Accept       multipart/form-data
// @Param        truck_name      formData string   true  "Truck Name"
// @Param        license_plate   formData string   true  "License Plate"
// @Param        make            formData string   true  "Make"
// @Param        model           formData string   true  "Model"
// @Param        year            formData int      true  "Year"
// @Param        color           formData string   true  "Color"
// @Param        length          formData number   true  "Length (ft)"
// @Param        width           formData number   true  "Width (ft)"
// @Param        height          formData number   true  "Height (ft)"
// @Param        max_weight      formData number   true  "Max Weight (lbs)"
// @Param        truck_type      formData string   true  "Truck Type" Enums(Small, Medium, Large)
// @Param        climate_control formData boolean  false "Climate Control"
// @Param        liftgate        formData boolean  false "Liftgate"
// @Param        pallet_jack     formData boolean  false "Pallet Jack"
// @Param        security_system formData boolean  false "Security System"
// @Param        refrigerated    formData boolean  false "Refrigerated"
// @Param        furniture_pads  formData boolean  false "Furniture Pads"
// @Param        photo           formData []file   false "Truck photo(s)"
// @Security     BearerAuth
// @Router       /trucks/ [post]
func CreateTruck(truckService service.TruckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
			return
		}
		files := form.File["photo"]

		year, _ := strconv.Atoi(c.PostForm("year"))
		length, _ := strconv.ParseFloat(c.PostForm("length"), 64)
		width, _ := strconv.ParseFloat(c.PostForm("width"), 64)
		height, _ := strconv.ParseFloat(c.PostForm("height"), 64)
		maxWeight, _ := strconv.ParseFloat(c.PostForm("max_weight"), 64)

		truck := &models.Truck{
			UserID:         userID,
			TruckName:      c.PostForm("truck_name"),
			LicensePlate:   c.PostForm("license_plate"),
			Make:           c.PostForm("make"),
			Model:          c.PostForm("model"),
			Year:           year,
			Color:          c.PostForm("color"),
			Length:         length,
			Width:          width,
			Height:         height,
			MaxWeight:      maxWeight,
			TruckType:      c.PostForm("truck_type"),
			ClimateControl: c.PostForm("climate_control") == "true",
			Liftgate:       c.PostForm("liftgate") == "true",
			PalletJack:     c.PostForm("pallet_jack") == "true",
			SecuritySystem: c.PostForm("security_system") == "true",
			Refrigerated:   c.PostForm("refrigerated") == "true",
			FurniturePads:  c.PostForm("furniture_pads") == "true",
			Photos:         files,
		}

		if err := truckService.CreateTruck(c.Request.Context(), truck); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Truck created successfully"})
	}
}
