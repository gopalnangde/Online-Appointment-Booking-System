package dto

type CreateAppointmentRequest struct {
	ProviderID   uint   `json:"provider_id" binding:"required"`
	ServiceTitle string `json:"service_title" binding:"required"`
	BookingDate  string `json:"booking_date" binding:"required"`
	BookingTime  string `json:"booking_time" binding:"required"`
	Notes        string `json:"notes"`
}

type UpdateAppointmentStatusRequest struct {
	Status string `json:"status" binding:"required"` // Confirmed, Completed, Cancelled
}
