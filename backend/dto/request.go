package dto

// RegisterRequest is the DTO for user registration input.
// Validation tags are processed by go-playground/validator.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,len=10,numeric"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=Customer ServiceProvider"`
}

// LoginRequest is the DTO for user login input.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserProfileRequest is the DTO for updating user/customer profile information.
type UpdateUserProfileRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,len=10,numeric"`
	Password string `json:"password" validate:"omitempty,min=8"`
}
