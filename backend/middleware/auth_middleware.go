package middleware

import (
	"net/http"
	"strings"

	"backend/constants"
	"backend/dto"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token from the Authorization header.
// On success, it stores the authenticated user's claims in the Gin context.
// On failure, it aborts the request with a 401 Unauthorized response.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: constants.MsgTokenMissing,
			})
			c.Abort()
			return
		}

		// Validate the "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: constants.MsgInvalidTokenFormat,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate the JWT token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: constants.MsgTokenInvalid,
			})
			c.Abort()
			return
		}

		// Store the authenticated user's information in the Gin context
		// so downstream handlers can access it
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.Name)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RequireCustomer ensures that only users with the "Customer" role can proceed.
// Must be used after AuthMiddleware.
func RequireCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role.(string) != constants.RoleCustomer {
			c.JSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: constants.MsgForbidden,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireServiceProvider ensures that only users with the "ServiceProvider" role can proceed.
// Must be used after AuthMiddleware.
func RequireServiceProvider() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role.(string) != constants.RoleServiceProvider {
			c.JSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: constants.MsgForbidden,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
