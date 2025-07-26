package router

// import (
// 	"moveshare/internal/auth"
// 	"moveshare/internal/handlers/job"
// 	"moveshare/internal/middleware"
// 	"moveshare/internal/service"

// 	"github.com/gin-gonic/gin"
// )

// func JobRouter(r *gin.Engine, jobService service.JobService, jwtAuth auth.JWTAuth) {
// 	jobGroup := r.Group("/jobs")
// 	jobGroup.Use(middleware.AuthMiddleware(jwtAuth))
// 	{
// 		jobGroup.POST("", job.CreateJob(jobService))
// 		jobGroup.GET("/available", job.GetAvailableJobs(jobService))
// 		jobGroup.GET("/my", job.GetUserJobs(jobService))
// 		jobGroup.DELETE("/:id", job.DeleteJob(jobService))
// 		jobGroup.POST("/:id/apply", job.ApplyForJob(jobService))
// 		jobGroup.GET("/applications/my", job.GetMyApplications(jobService))
// 	}
// }
