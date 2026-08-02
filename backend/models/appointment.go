package models

import "time"

// Appointment represents the appointments table in the database.
type Appointment struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID   uint      `gorm:"not null;index" json:"customer_id"`
	ProviderID   uint      `gorm:"not null;index" json:"provider_id"`
	ServiceTitle string    `gorm:"type:varchar(150);not null" json:"service_title"`
	BookingDate  string    `gorm:"type:varchar(20);not null" json:"booking_date"`
	BookingTime  string    `gorm:"type:varchar(20);not null" json:"booking_time"`
	Notes        string    `gorm:"type:text" json:"notes"`
	Status       string    `gorm:"type:varchar(20);default:'Pending';not null" json:"status"` // Pending, Confirmed, Completed, Cancelled
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Customer User                   `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"customer,omitempty"`
	Provider User                   `gorm:"foreignKey:ProviderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"provider,omitempty"`
	Profile  ServiceProviderProfile `gorm:"foreignKey:ProviderID;references:UserID" json:"provider_profile,omitempty"`
}

// TableName overrides the default table name used by GORM
func (Appointment) TableName() string {
	return "appointments"
}
