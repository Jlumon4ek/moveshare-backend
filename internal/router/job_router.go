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
		protected.GET("/available-jobs/", jobHandler.GetAvailableJobs)    // уже обновлен
		protected.GET("/filter-options/", jobHandler.GetJobFilterOptions) // новый эндпоинт
		protected.GET("/stats/", jobHandler.GetJobsStats)                 // статистика работ
		protected.DELETE("/delete-job/:id/", jobHandler.DeleteJob)
		protected.GET("/my-jobs/", jobHandler.GetMyJobs)
		protected.GET("/:id/details/", jobHandler.GetJobByID)
		protected.GET("/claimed-jobs/", jobHandler.GetClaimedJobs)
		protected.GET("/pending-jobs/", jobHandler.GetPendingJobs)
		protected.GET("/today-schedule/", jobHandler.GetTodayScheduleJobs)
		protected.GET("/user-work-stats/", jobHandler.GetUserWorkStats)
		protected.POST("/mark-job-completed/:id/", jobHandler.MarkJobCompleted)
		protected.POST("/export-jobs/", jobHandler.ExportJobs)
		protected.POST("/upload-files/:id/", jobHandler.UploadJobFiles)
		protected.POST("/upload-verification-documents/:id/", jobHandler.UploadVerificationDocuments)
		protected.POST("/upload-work-photos/:id/", jobHandler.UploadWorkPhotos)
		protected.GET("/:id/files/", jobHandler.GetJobFiles)
		protected.GET("/:id/files/by-type/", jobHandler.GetJobFilesByType)
	}
}
