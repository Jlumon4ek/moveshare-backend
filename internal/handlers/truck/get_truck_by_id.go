package truck

import (
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTruckByID godoc
// @Summary      Get truck by ID
// @Description  Get a truck by its ID
// @Tags         Trucks
// @Param        truckId path int true "Truck ID"
// @Router       /trucks/{truckId}/ [get]
// @Security     BearerAuth
func GetTruckByID(truckService service.TruckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("truckId")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing truckId in path"})
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid truckId format"})
			return
		}

		truck, err := truckService.GetTruckByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if truck == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Truck not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"truck": truck})
	}
}
