package router

import (
	"moveshare/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupLocationRoutes(r gin.IRouter, locationHandler *handlers.LocationHandler) {
	r.GET("/states", locationHandler.GetStates)
	r.GET("/cities", locationHandler.GetCities)
}
