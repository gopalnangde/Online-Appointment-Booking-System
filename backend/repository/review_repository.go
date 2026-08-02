package repository

import (
	"backend/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	GetByProviderID(providerID uint) ([]models.Review, error)
	GetStatsByProviderID(providerID uint) (float64, int, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetByProviderID(providerID uint) ([]models.Review, error) {
	var reviews []models.Review
	var profile models.ServiceProviderProfile
	targetUserID := providerID
	if err := r.db.Where("id = ? OR user_id = ?", providerID, providerID).First(&profile).Error; err == nil {
		targetUserID = profile.UserID
	}

	err := r.db.Where("provider_id = ? OR provider_id = ?", providerID, targetUserID).
		Preload("Customer").
		Order("created_at desc").
		Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetStatsByProviderID(providerID uint) (float64, int, error) {
	type Result struct {
		AvgRating float64
		Total     int
	}
	var res Result
	var profile models.ServiceProviderProfile
	targetUserID := providerID
	if err := r.db.Where("id = ? OR user_id = ?", providerID, providerID).First(&profile).Error; err == nil {
		targetUserID = profile.UserID
	}

	err := r.db.Model(&models.Review{}).
		Select("ROUND(COALESCE(AVG(rating), 0), 1) as avg_rating, COUNT(id) as total").
		Where("provider_id = ? OR provider_id = ?", providerID, targetUserID).
		Scan(&res).Error
	return res.AvgRating, res.Total, err
}
