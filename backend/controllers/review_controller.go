package controllers

import (
	"net/http"
	"strconv"

	"backend/dto"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	service services.ReviewService
}

func NewReviewController(service services.ReviewService) *ReviewController {
	return &ReviewController{service: service}
}

func (c *ReviewController) CreateReview(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var req dto.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid request: " + err.Error()})
		return
	}

	review, err := c.service.CreateReview(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Review submitted successfully",
		Data:    review,
	})
}

func (c *ReviewController) GetProviderReviews(ctx *gin.Context) {
	idStr := ctx.Param("providerId")
	providerID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid provider ID"})
		return
	}

	reviews, stats, err := c.service.GetProviderReviews(uint(providerID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Reviews retrieved successfully",
		Data: gin.H{
			"reviews": reviews,
			"stats":   stats,
		},
	})
}
