package main

import (
	"fmt"
	"log"

	"backend/config"
	"backend/controllers"
	"backend/database"
	"backend/middleware"
	"backend/repository"
	"backend/routes"
	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Step 1: Load configuration from .env
	cfg := config.LoadConfig()

	// Step 2: Connect to MySQL and run auto-migrations
	database.Connect(cfg)

	// Step 3: Initialize JWT secret from configuration
	utils.SetJWTSecret(cfg.JWTSecret)

	// Step 4: Wire up dependencies (Repository → Service → Controller)
	db := database.GetDB()

	// Auth module
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// Review module
	reviewRepo := repository.NewReviewRepository(db)
	reviewService := services.NewReviewService(reviewRepo)
	reviewController := controllers.NewReviewController(reviewService)

	// Provider profile module
	providerRepo := repository.NewProviderRepository(db)
	providerService := services.NewProviderService(providerRepo, reviewRepo)
	providerController := controllers.NewProviderController(providerService)

	// Appointment module
	aptRepo := repository.NewAppointmentRepository(db)
	aptService := services.NewAppointmentService(aptRepo)
	aptController := controllers.NewAppointmentController(aptService)

	// Step 5: Create the Gin router and register routes
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.SetupAuthRoutes(router, authController)
	routes.SetupProviderRoutes(router, providerController)
	routes.SetupAppointmentRoutes(router, aptController)
	routes.SetupReviewRoutes(router, reviewController)

	// Step 6: Start the server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}