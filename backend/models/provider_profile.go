package models

import "time"

// ServiceProviderProfile represents the service_provider_profiles table.
// Each ServiceProvider user has exactly one profile linked via UserID.
type ServiceProviderProfile struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Specialization string    `gorm:"type:varchar(100);not null" json:"specialization"`
	Description    string    `gorm:"type:text" json:"description"`
	Address        string    `gorm:"type:varchar(255);not null" json:"address"`
	City           string    `gorm:"type:varchar(100);not null" json:"city"`
	State          string    `gorm:"type:varchar(100);not null" json:"state"`
	PinCode        string    `gorm:"type:varchar(10);not null" json:"pin_code"`
	Experience     int       `gorm:"not null;default:0" json:"experience"` // Years of experience
	AvailableSlots string    `gorm:"type:text" json:"available_slots"`     // JSON array of time slots
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Belongs-to relationship with User
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName overrides the default table name used by GORM
func (ServiceProviderProfile) TableName() string {
	return "service_provider_profiles"
}
