package router

import (
	"moveshare/internal/handlers/company"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func CompanyRouter(r *gin.Engine, companyService service.CompanyService, jwtAuth service.JWTAuth) {
	companyGroup := r.Group("/company")
	companyGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		companyGroup.GET("/", company.GetCompany(companyService))
		companyGroup.PATCH("/", company.PatchCompany(companyService))
	}
}
