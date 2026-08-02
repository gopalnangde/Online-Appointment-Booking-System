package controllers

import (
	"net/http"
	"strconv"

	"backend/dto"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type AppointmentController struct {
	service services.AppointmentService
}

func NewAppointmentController(service services.AppointmentService) *AppointmentController {
	return &AppointmentController{service: service}
}

func (c *AppointmentController) CreateAppointment(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var req dto.CreateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid request payload: " + err.Error()})
		return
	}

	apt, err := c.service.CreateAppointment(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Appointment booked successfully",
		Data:    apt,
	})
}

func (c *AppointmentController) GetMyAppointments(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)
	roleVal, _ := ctx.Get("userRole")
	role := roleVal.(string)

	var appointments interface{}
	var err error

	if role == "ServiceProvider" {
		appointments, err = c.service.GetProviderAppointments(userID)
	} else {
		appointments, err = c.service.GetCustomerAppointments(userID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Appointments retrieved successfully",
		Data:    appointments,
	})
}

func (c *AppointmentController) UpdateAppointmentStatus(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid appointment ID"})
		return
	}

	var req dto.UpdateAppointmentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid payload"})
		return
	}

	err = c.service.UpdateStatus(userID, uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Appointment status updated to " + req.Status,
	})
}
