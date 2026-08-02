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

// AuthController handles HTTP requests for authentication endpoints.
type AuthController struct {
	authService services.AuthService
}

// NewAuthController creates a new AuthController with the given service dependency.
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register handles POST /api/auth/register
// It validates the request body, delegates to the service layer, and returns a structured response.
func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind JSON body to the request DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgInvalidRequestBody,
		})
		return
	}

	// Validate the request fields using go-playground/validator
	if err := validators.Validate.Struct(req); err != nil {
		errors := validators.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgValidationFailed,
			Data:    errors,
		})
		return
	}

	// Delegate to service layer
	if err := ctrl.authService.Register(req); err != nil {
		// Determine appropriate status code based on error message
		statusCode := http.StatusInternalServerError
		if err.Error() == constants.MsgEmailExists || err.Error() == constants.MsgPhoneExists {
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
		Message: constants.MsgRegistrationSuccess,
	})
}

// Login handles POST /api/auth/login
// It validates credentials and returns a JWT token on success.
func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind JSON body to the request DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgInvalidRequestBody,
		})
		return
	}

	// Validate the request fields
	if err := validators.Validate.Struct(req); err != nil {
		errors := validators.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgValidationFailed,
			Data:    errors,
		})
		return
	}

	// Delegate to service layer
	loginResponse, err := ctrl.authService.Login(req)
	if err != nil {
		// Return 401 for invalid credentials, 500 for other errors
		statusCode := http.StatusInternalServerError
		if err.Error() == constants.MsgInvalidCredentials {
			statusCode = http.StatusUnauthorized
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: constants.MsgLoginSuccess,
		Data:    loginResponse,
	})
}

// GetProfile handles GET /api/auth/profile
// It reads the authenticated user's ID from the Gin context (set by auth middleware)
// and returns their profile information.
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	// Extract user ID set by the authentication middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: constants.MsgUnauthorized,
		})
		return
	}

	// Fetch user profile from service layer
	userResponse, err := ctrl.authService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: constants.MsgProfileFetched,
		Data:    userResponse,
	})
}

// GetUserByID handles GET /api/auth/users/:id
// Returns profile information for any user by their ID.
func (ctrl *AuthController) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	userResponse, err := ctrl.authService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: constants.MsgProfileFetched,
		Data:    userResponse,
	})
}

// UpdateProfile handles PUT /api/auth/profile
// Updates the authenticated user's profile information (Name, Email, Phone, Password).
func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: constants.MsgUnauthorized,
		})
		return
	}

	var req dto.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgInvalidRequestBody,
		})
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		errors := validators.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: constants.MsgValidationFailed,
			Data:    errors,
		})
		return
	}

	userResponse, err := ctrl.authService.UpdateProfile(userID.(uint), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == constants.MsgEmailExists || err.Error() == constants.MsgPhoneExists {
			statusCode = http.StatusConflict
		} else if err.Error() == constants.MsgUserNotFound {
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
		Data:    userResponse,
	})
}
