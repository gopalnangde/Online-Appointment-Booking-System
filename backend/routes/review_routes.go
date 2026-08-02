package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupReviewRoutes(router *gin.Engine, controller *controllers.ReviewController) {
	api := router.Group("/api/reviews")
	{
		// Public route to view provider reviews
		api.GET("/provider/:providerId", controller.GetProviderReviews)

		// Protected route to post a review
		api.POST("", middleware.AuthMiddleware(), controller.CreateReview)
		api.POST("/", middleware.AuthMiddleware(), controller.CreateReview)
	}
}
