package router

import (
	"moveshare/internal/handlers"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupJobRoutes(r gin.IRouter, jobHandler *handlers.JobHandler, jwtAuth service.JWTAuth) {
	protected := r.Group("/jobs")
	protected.Use(middleware.AuthMiddleware(jwtAuth))
	{

		protected.POST("/post-new-job/", jobHandler.PostNewJob)
		protected.POST("/claim-job/:id/", jobHandler.ClaimJob)
		protected.GET("/available-jobs/", jobHandler.GetAvailableJobs)
		protected.DELETE("/delete-job/:id/", jobHandler.DeleteJob)
		protected.GET("/my-jobs/", jobHandler.GetMyJobs)
		protected.GET("/job/:id/", jobHandler.GetJobByID)
		protected.GET("/claimed-jobs/", jobHandler.GetClaimedJobs)
	}
}
