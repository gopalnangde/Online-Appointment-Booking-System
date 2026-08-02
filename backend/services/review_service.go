package services

import (
	"errors"

	"backend/dto"
	"backend/models"
	"backend/repository"
)

type ReviewService interface {
	CreateReview(customerID uint, req *dto.CreateReviewRequest) (*models.Review, error)
	GetProviderReviews(providerID uint) ([]models.Review, *dto.ProviderReviewStats, error)
}

type reviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) ReviewService {
	return &reviewService{repo: repo}
}

func (s *reviewService) CreateReview(customerID uint, req *dto.CreateReviewRequest) (*models.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	review := &models.Review{
		AppointmentID: req.AppointmentID,
		CustomerID:    customerID,
		ProviderID:    req.ProviderID,
		Rating:        req.Rating,
		Comment:       req.Comment,
	}

	err := s.repo.Create(review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) GetProviderReviews(providerID uint) ([]models.Review, *dto.ProviderReviewStats, error) {
	reviews, err := s.repo.GetByProviderID(providerID)
	if err != nil {
		return nil, nil, err
	}

	avg, total, err := s.repo.GetStatsByProviderID(providerID)
	if err != nil {
		return nil, nil, err
	}

	stats := &dto.ProviderReviewStats{
		AverageRating: avg,
		TotalReviews:  total,
	}

	return reviews, stats, nil
}
