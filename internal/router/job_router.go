package router

import (
	"moveshare/internal/handlers/job"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func JobRouter(r *gin.Engine, jobService service.JobService, jwtAuth service.JWTAuth) {
	jobGroup := r.Group("/jobs")
	jobGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		jobGroup.POST("/post-new-job/", job.PostNewJob(jobService))
		jobGroup.GET("/available/", job.GetAvailableJobs(jobService))
		jobGroup.GET("/my", job.GetMyJobs(jobService))
		jobGroup.DELETE("/:jobID", job.DeleteJob(jobService))
		jobGroup.POST("/:jobID/apply", job.ApplyForJob(jobService))
		jobGroup.GET("/applications/my", job.GetMyApplications(jobService))
	}
}
