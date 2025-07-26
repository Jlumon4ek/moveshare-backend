package truck

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserTrucks godoc
// @Summary      Get trucks for the authenticated user
// @Description  Retrieves all trucks associated with the authenticated user.
// @Tags         Trucks
// @Security     BearerAuth
// @Router       /trucks/ [get]
func GetUserTrucks(truckService service.TruckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		trucks, err := truckService.GetUserTrucks(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"trucks": trucks})
	}
}
