package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes registers all authentication-related routes under /api/auth.
// Public routes (register, login) have no middleware.
// Protected routes (profile) require JWT authentication.
func SetupAuthRoutes(router *gin.Engine, authController *controllers.AuthController) {
	// Group all auth routes under /api/auth
	auth := router.Group("/api/auth")
	{
		// Public routes — no authentication required
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)

		// Protected routes — require valid JWT token
		protected := auth.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", authController.GetProfile)
			protected.PUT("/profile", authController.UpdateProfile)
			protected.GET("/users/:id", authController.GetUserByID)
		}
	}
}
