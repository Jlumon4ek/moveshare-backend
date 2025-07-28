package router

import (
	"moveshare/internal/handlers/truck"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func TruckRouter(r gin.IRouter, truckService service.TruckService, jwtAuth service.JWTAuth) {
	truckGroup := r.Group("/trucks")
	truckGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		truckGroup.POST("/", truck.CreateTruck(truckService))
		truckGroup.GET("/", truck.GetUserTrucks(truckService))
		truckGroup.GET("/:truckId/", truck.GetTruckByID(truckService))
		truckGroup.DELETE("/:truckId/", truck.DeleteTruck(truckService))

	}
}
