package truck

import (
	"moveshare/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteTruck godoc
// @Summary      Delete a truck
// @Description  Deletes a truck by its ID
// @Tags         Trucks
// @Param        truckId path int true "Truck ID"
// @Router       /trucks/{truckId}/ [delete]
// @Security     BearerAuth
func DeleteTruck(truckService service.TruckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		truckIdStr := c.Param("truckId")
		truckId, err := strconv.ParseInt(truckIdStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid truckId"})
			return
		}
		if err := truckService.DeleteTruck(c, truckId); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Truck deleted successfully"})
	}
}
