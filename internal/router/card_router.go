package router

import (
	// "moveshare/internal/handlers/card"
	"moveshare/internal/middleware"
	"moveshare/internal/service"

	"github.com/gin-gonic/gin"
)

func CardRouter(r *gin.Engine, cardService service.CardService, jwtAuth service.JWTAuth) {
	cardGroup := r.Group("/cards")
	cardGroup.Use(middleware.AuthMiddleware(jwtAuth))
	{
		// cardGroup.POST("", card.CreateCard(cardService))
		// cardGroup.GET("", card.GetUserCards(cardService))
		// cardGroup.GET("/:id", card.GetCardByID(cardService))
		// cardGroup.PUT("/:id", card.UpdateCard(cardService))
		// cardGroup.DELETE("/:id", card.DeleteCard(cardService))
		// cardGroup.POST("/:id/default", card.SetDefaultCard(cardService))
	}
}
