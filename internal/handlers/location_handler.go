package handlers

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationService *service.LocationService
}

func NewLocationHandler(locationService *service.LocationService) *LocationHandler {
	return &LocationHandler{locationService: locationService}
}

// GetStates godoc
// @Summary Get all states
// @Description Retrieves a list of all available states
// @Tags Location
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]models.State "List of states"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /location/states [get]
func (h *LocationHandler) GetStates(c *gin.Context) {
	states, err := h.locationService.GetAllStates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"states": states,
	})
}

// GetCities godoc
// @Summary Get cities
// @Description Retrieves a list of cities, optionally filtered by state
// @Tags Location
// @Accept json
// @Produce json
// @Param state_id query int false "State ID to filter cities"
// @Success 200 {object} map[string][]models.CityWithState "List of cities"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /location/cities [get]
func (h *LocationHandler) GetCities(c *gin.Context) {
	var query models.CitiesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cities, err := h.locationService.GetCities(query.StateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cities": cities,
	})
}
