package services

import (
	"encoding/json"
	"errors"

	"backend/constants"
	"backend/dto"
	"backend/models"
	"backend/repository"
)

var defaultTimeSlots = []string{
	"09:00 AM", "10:00 AM", "11:00 AM", "12:00 PM",
	"02:00 PM", "03:00 PM", "04:00 PM", "05:00 PM", "06:00 PM",
}

// ProviderService defines the contract for provider profile business logic.
type ProviderService interface {
	CreateProfile(userID uint, req dto.CreateProviderProfileRequest) error
	GetOwnProfile(userID uint) (*dto.ProviderProfileResponse, error)
	UpdateProfile(userID uint, req dto.UpdateProviderProfileRequest) error
	GetAllProviders(page, limit int) (*dto.ProviderListResponse, error)
	GetProviderByID(profileID uint) (*dto.ProviderProfileResponse, error)
}

// providerService is the concrete implementation of ProviderService.
type providerService struct {
	providerRepo repository.ProviderRepository
	reviewRepo   repository.ReviewRepository
}

// NewProviderService creates a new ProviderService with the given repository dependencies.
func NewProviderService(providerRepo repository.ProviderRepository, reviewRepo repository.ReviewRepository) ProviderService {
	return &providerService{
		providerRepo: providerRepo,
		reviewRepo:   reviewRepo,
	}
}

// CreateProfile creates a new provider profile for the authenticated user.
// Ensures only one profile per provider by checking for duplicates.
func (s *providerService) CreateProfile(userID uint, req dto.CreateProviderProfileRequest) error {
	// Check if profile already exists for this user
	existing, err := s.providerRepo.FindByUserID(userID)
	if err != nil {
		return errors.New(constants.MsgInternalServerError)
	}
	if existing != nil {
		return errors.New(constants.MsgProfileExists)
	}

	slots := req.AvailableSlots
	if len(slots) == 0 {
		slots = defaultTimeSlots
	}
	slotsJSON, _ := json.Marshal(slots)

	// Create the profile model from the request DTO
	profile := &models.ServiceProviderProfile{
		UserID:         userID,
		Specialization: req.Specialization,
		Description:    req.Description,
		Address:        req.Address,
		City:           req.City,
		State:          req.State,
		PinCode:        req.PinCode,
		Experience:     req.Experience,
		AvailableSlots: string(slotsJSON),
	}

	if err := s.providerRepo.CreateProfile(profile); err != nil {
		return errors.New(constants.MsgInternalServerError)
	}

	return nil
}

// GetOwnProfile retrieves the authenticated provider's own profile.
func (s *providerService) GetOwnProfile(userID uint) (*dto.ProviderProfileResponse, error) {
	profile, err := s.providerRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}
	if profile == nil {
		return nil, errors.New(constants.MsgProviderNotFound)
	}

	return s.mapProfileToResponse(profile), nil
}

// UpdateProfile updates the authenticated provider's profile.
// Only non-zero/non-empty fields from the request are applied.
func (s *providerService) UpdateProfile(userID uint, req dto.UpdateProviderProfileRequest) error {
	// Fetch the existing profile
	profile, err := s.providerRepo.FindByUserID(userID)
	if err != nil {
		return errors.New(constants.MsgInternalServerError)
	}
	if profile == nil {
		return errors.New(constants.MsgProviderNotFound)
	}

	// Apply only the fields that were provided (non-zero values)
	if req.Specialization != "" {
		profile.Specialization = req.Specialization
	}
	if req.Description != "" {
		profile.Description = req.Description
	}
	if req.Address != "" {
		profile.Address = req.Address
	}
	if req.City != "" {
		profile.City = req.City
	}
	if req.State != "" {
		profile.State = req.State
	}
	if req.PinCode != "" {
		profile.PinCode = req.PinCode
	}
	if req.Experience > 0 {
		profile.Experience = req.Experience
	}
	if req.AvailableSlots != nil {
		slotsJSON, err := json.Marshal(req.AvailableSlots)
		if err == nil {
			profile.AvailableSlots = string(slotsJSON)
		}
	}

	if err := s.providerRepo.UpdateProfile(profile); err != nil {
		return errors.New(constants.MsgInternalServerError)
	}

	return nil
}

// GetAllProviders retrieves a paginated list of all service providers.
func (s *providerService) GetAllProviders(page, limit int) (*dto.ProviderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	profiles, total, err := s.providerRepo.FindAll(page, limit)
	if err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}

	var providers []dto.ProviderProfileResponse
	for _, profile := range profiles {
		providers = append(providers, *s.mapProfileToResponse(&profile))
	}

	return &dto.ProviderListResponse{
		Providers: providers,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// GetProviderByID retrieves a single provider's profile by profile ID or user ID.
func (s *providerService) GetProviderByID(profileID uint) (*dto.ProviderProfileResponse, error) {
	profile, err := s.providerRepo.FindByID(profileID)
	if err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}
	if profile == nil {
		profile, err = s.providerRepo.FindByUserID(profileID)
		if err != nil {
			return nil, errors.New(constants.MsgInternalServerError)
		}
	}
	if profile == nil {
		return nil, errors.New(constants.MsgProviderNotFound)
	}

	return s.mapProfileToResponse(profile), nil
}

// mapProfileToResponse converts a ServiceProviderProfile model (with preloaded User)
// into a ProviderProfileResponse DTO — merging user details and review statistics into the response.
func (s *providerService) mapProfileToResponse(profile *models.ServiceProviderProfile) *dto.ProviderProfileResponse {
	var slots []string
	if profile.AvailableSlots != "" {
		_ = json.Unmarshal([]byte(profile.AvailableSlots), &slots)
	}
	if len(slots) == 0 {
		slots = defaultTimeSlots
	}

	resp := &dto.ProviderProfileResponse{
		ID:             profile.ID,
		UserID:         profile.UserID,
		Name:           profile.User.Name,
		Email:          profile.User.Email,
		Phone:          profile.User.Phone,
		Specialization: profile.Specialization,
		Description:    profile.Description,
		Address:        profile.Address,
		City:           profile.City,
		State:          profile.State,
		PinCode:        profile.PinCode,
		Experience:     profile.Experience,
		AvailableSlots: slots,
	}

	if s.reviewRepo != nil {
		avg, total, err := s.reviewRepo.GetStatsByProviderID(profile.UserID)
		if err == nil {
			resp.AverageRating = avg
			resp.TotalReviews = total
		}
	}

	return resp
}
