package transport

import (
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *AuthHandler, jwt services.JWTService) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		users := api.Group("/users")
		users.Use(AuthMiddleware(jwt))
		{
			users.GET("/me", authHandler.Me)
			users.PUT("/me", authHandler.UpdateMe)
		}
	}

	return r
}
