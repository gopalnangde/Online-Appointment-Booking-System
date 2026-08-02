package repository

import (
	"backend/models"

	"gorm.io/gorm"
)

// ProviderRepository defines the contract for provider profile database operations.
type ProviderRepository interface {
	CreateProfile(profile *models.ServiceProviderProfile) error
	FindByUserID(userID uint) (*models.ServiceProviderProfile, error)
	FindByID(id uint) (*models.ServiceProviderProfile, error)
	UpdateProfile(profile *models.ServiceProviderProfile) error
	FindAll(page, limit int) ([]models.ServiceProviderProfile, int64, error)
}

// providerRepository is the concrete implementation of ProviderRepository.
type providerRepository struct {
	db *gorm.DB
}

// NewProviderRepository creates a new ProviderRepository instance.
func NewProviderRepository(db *gorm.DB) ProviderRepository {
	return &providerRepository{db: db}
}

// CreateProfile inserts a new provider profile record into the database.
func (r *providerRepository) CreateProfile(profile *models.ServiceProviderProfile) error {
	return r.db.Create(profile).Error
}

// FindByUserID retrieves a provider profile by the associated user ID.
// Preloads the User relation to include provider's name, email, and phone.
// Returns nil (without error) if no profile is found.
func (r *providerRepository) FindByUserID(userID uint) (*models.ServiceProviderProfile, error) {
	var profile models.ServiceProviderProfile
	result := r.db.Preload("User").Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

// FindByID retrieves a provider profile by its primary key.
// Preloads the User relation.
// Returns nil (without error) if no profile is found.
func (r *providerRepository) FindByID(id uint) (*models.ServiceProviderProfile, error) {
	var profile models.ServiceProviderProfile
	result := r.db.Preload("User").First(&profile, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &profile, nil
}

// UpdateProfile persists changes to an existing provider profile.
func (r *providerRepository) UpdateProfile(profile *models.ServiceProviderProfile) error {
	return r.db.Save(profile).Error
}

// FindAll retrieves a paginated list of all provider profiles.
// Preloads the User relation and returns total count for pagination.
func (r *providerRepository) FindAll(page, limit int) ([]models.ServiceProviderProfile, int64, error) {
	var profiles []models.ServiceProviderProfile
	var total int64

	// Count total records
	if err := r.db.Model(&models.ServiceProviderProfile{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results with User preloaded
	offset := (page - 1) * limit
	result := r.db.Preload("User").Offset(offset).Limit(limit).Order("id DESC").Find(&profiles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return profiles, total, nil
}
