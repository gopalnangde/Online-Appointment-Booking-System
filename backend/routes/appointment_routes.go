package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAppointmentRoutes(router *gin.Engine, controller *controllers.AppointmentController) {
	api := router.Group("/api/appointments")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("", controller.CreateAppointment)
		api.POST("/", controller.CreateAppointment)
		api.GET("", controller.GetMyAppointments)
		api.GET("/", controller.GetMyAppointments)
		api.PATCH("/:id/status", controller.UpdateAppointmentStatus)
		api.PUT("/:id/status", controller.UpdateAppointmentStatus)
	}
}
