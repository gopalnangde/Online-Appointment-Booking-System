package controllers

import (
	"net/http"
	"strconv"

	"backend/constants"
	"backend/dto"
	"backend/services"
	"backend/validators"

	"github.com/gin-gonic/gin"
)

// ProviderController handles HTTP requests for provider profile endpoints.
type ProviderController struct {
	providerService services.ProviderService
}

// NewProviderController creates a new ProviderController with the given service dependency.
func NewProviderController(providerService services.ProviderService) *ProviderController {
	return &ProviderController{providerService: providerService}
}

// CreateProfile handles POST /api/provider/profile
// Creates a new service provider profile for the authenticated ServiceProvider user.
func (ctrl *ProviderController) CreateProfile(c *gin.Context) {
	var req dto.CreateProviderProfileRequest

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgInvalidRequestBody,
		})
		return
	}

	// Validate request fields
	if err := validators.Validate.Struct(req); err != nil {
		errors := validators.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgValidationFailed,
			Data:    errors,
		})
		return
	}

	// Get authenticated user ID from context
	userID, _ := c.Get("userID")

	// Delegate to service layer
	if err := ctrl.providerService.CreateProfile(userID.(uint), req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == constants.MsgProfileExists {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: constants.MsgProfileCreated,
	})
}

// GetOwnProfile handles GET /api/provider/profile
// Returns the authenticated provider's own profile.
func (ctrl *ProviderController) GetOwnProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	profile, err := ctrl.providerService.GetOwnProfile(userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == constants.MsgProviderNotFound {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: constants.MsgProfileFetched,
		Data:    profile,
	})
}

// UpdateProfile handles PUT /api/provider/profile
// Updates the authenticated provider's profile with the provided fields.
func (ctrl *ProviderController) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProviderProfileRequest

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgInvalidRequestBody,
		})
		return
	}

	// Validate request fields
	if err := validators.Validate.Struct(req); err != nil {
		errors := validators.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgValidationFailed,
			Data:    errors,
		})
		return
	}

	userID, _ := c.Get("userID")

	// Delegate to service layer
	if err := ctrl.providerService.UpdateProfile(userID.(uint), req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == constants.MsgProviderNotFound {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: constants.MsgProfileUpdated,
	})
}

// GetAllProviders handles GET /api/providers
// Returns a paginated list of all service providers (public endpoint).
func (ctrl *ProviderController) GetAllProviders(c *gin.Context) {
	// Parse pagination query params with defaults
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := ctrl.providerService.GetAllProviders(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: constants.MsgProvidersFetched,
		Data:    result,
	})
}

// GetProviderByID handles GET /api/providers/:id
// Returns a single provider's profile details (public endpoint).
func (ctrl *ProviderController) GetProviderByID(c *gin.Context) {
	// Parse profile ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid provider ID",
		})
		return
	}

	profile, err := ctrl.providerService.GetProviderByID(uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == constants.MsgProviderNotFound {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: constants.MsgProvidersFetched,
		Data:    profile,
	})
}
