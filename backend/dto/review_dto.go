package dto

type CreateReviewRequest struct {
	AppointmentID uint   `json:"appointment_id"`
	ProviderID    uint   `json:"provider_id" binding:"required"`
	Rating        int    `json:"rating" binding:"required,min=1,max=5"`
	Comment       string `json:"comment"`
}

type ProviderReviewStats struct {
	AverageRating float64 `json:"average_rating"`
	TotalReviews  int     `json:"total_reviews"`
}
