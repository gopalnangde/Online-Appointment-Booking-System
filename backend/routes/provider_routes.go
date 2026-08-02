package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupProviderRoutes registers all service provider profile routes.
// Provider-specific routes require authentication + ServiceProvider role.
// Public routes allow anyone to browse providers.
func SetupProviderRoutes(router *gin.Engine, providerController *controllers.ProviderController) {
	// Protected routes — require authentication + ServiceProvider role
	provider := router.Group("/api/provider")
	provider.Use(middleware.AuthMiddleware())
	provider.Use(middleware.RequireServiceProvider())
	{
		provider.POST("/profile", providerController.CreateProfile)
		provider.GET("/profile", providerController.GetOwnProfile)
		provider.PUT("/profile", providerController.UpdateProfile)
	}

	// Public routes — anyone can browse service providers
	providers := router.Group("/api/providers")
	{
		providers.GET("", providerController.GetAllProviders)
		providers.GET("/", providerController.GetAllProviders)
		providers.GET("/:id", providerController.GetProviderByID)
	}
}
