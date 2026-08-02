package dto

// CreateProviderProfileRequest is the DTO for creating a service provider profile.
type CreateProviderProfileRequest struct {
	Specialization string `json:"specialization" validate:"required"`
	Description    string `json:"description"`
	Address        string `json:"address" validate:"required"`
	City           string `json:"city" validate:"required"`
	State          string `json:"state" validate:"required"`
	PinCode        string `json:"pin_code" validate:"required,len=6,numeric"`
	Experience     int      `json:"experience" validate:"gte=0"`
	AvailableSlots []string `json:"available_slots"`
}

// UpdateProviderProfileRequest is the DTO for updating a service provider profile.
// All fields are optional — only non-zero values will be updated.
type UpdateProviderProfileRequest struct {
	Specialization string   `json:"specialization"`
	Description    string   `json:"description"`
	Address        string   `json:"address"`
	City           string   `json:"city"`
	State          string   `json:"state"`
	PinCode        string   `json:"pin_code" validate:"omitempty,len=6,numeric"`
	Experience     int      `json:"experience" validate:"gte=0"`
	AvailableSlots []string `json:"available_slots"`
}

// ProviderProfileResponse is the DTO returned when fetching a provider profile.
type ProviderProfileResponse struct {
	ID             uint     `json:"id"`
	UserID         uint     `json:"user_id"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	Specialization string   `json:"specialization"`
	Description    string   `json:"description"`
	Address        string   `json:"address"`
	City           string   `json:"city"`
	State          string   `json:"state"`
	PinCode        string   `json:"pin_code"`
	Experience     int      `json:"experience"`
	AvailableSlots []string `json:"available_slots"`
	AverageRating  float64  `json:"average_rating"`
	TotalReviews   int      `json:"total_reviews"`
}

// ProviderListResponse is the DTO for the paginated list of providers.
type ProviderListResponse struct {
	Providers []ProviderProfileResponse `json:"providers"`
	Total     int64                     `json:"total"`
	Page      int                       `json:"page"`
	Limit     int                       `json:"limit"`
}
