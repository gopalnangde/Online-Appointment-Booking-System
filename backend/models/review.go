package models

import "time"

// Review represents customer reviews and ratings for service providers.
type Review struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AppointmentID uint      `gorm:"index" json:"appointment_id"`
	CustomerID    uint      `gorm:"not null;index" json:"customer_id"`
	ProviderID    uint      `gorm:"not null;index" json:"provider_id"`
	Rating        int       `gorm:"not null" json:"rating"` // 1 to 5 stars
	Comment       string    `gorm:"type:text" json:"comment"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Customer User `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"customer,omitempty"`
}

// TableName overrides the default table name used by GORM
func (Review) TableName() string {
	return "reviews"
}
