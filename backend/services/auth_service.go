package services

import (
	"errors"

	"backend/constants"
	"backend/dto"
	"backend/models"
	"backend/repository"
	"backend/utils"
)

// AuthService defines the contract for authentication business logic.
type AuthService interface {
	Register(req dto.RegisterRequest) error
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(userID uint) (*dto.UserResponse, error)
	GetUserByID(userID uint) (*dto.UserResponse, error)
	UpdateProfile(userID uint, req dto.UpdateUserProfileRequest) (*dto.UserResponse, error)
}

// authService is the concrete implementation of AuthService.
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new AuthService with the given repository dependency.
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Register handles user registration business logic:
// 1. Check for duplicate email
// 2. Check for duplicate phone
// 3. Hash the password
// 4. Persist the user
func (s *authService) Register(req dto.RegisterRequest) error {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return errors.New(constants.MsgInternalServerError)
	}
	if existingUser != nil {
		return errors.New(constants.MsgEmailExists)
	}

	// Check if phone number already exists
	existingUser, err = s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		return errors.New(constants.MsgInternalServerError)
	}
	if existingUser != nil {
		return errors.New(constants.MsgPhoneExists)
	}

	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New(constants.MsgInternalServerError)
	}

	// Create the user model from the request DTO
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     req.Role,
	}

	// Persist the user to the database
	if err := s.userRepo.CreateUser(user); err != nil {
		return errors.New(constants.MsgInternalServerError)
	}

	return nil
}

// Login handles user authentication:
// 1. Look up user by email
// 2. Verify password
// 3. Generate JWT token
// 4. Return token and user data
func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}
	if user == nil {
		return nil, errors.New(constants.MsgInvalidCredentials)
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New(constants.MsgInvalidCredentials)
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}

	// Build and return the login response
	loginResponse := &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Role:  user.Role,
		},
	}

	return loginResponse, nil
}

// GetProfile retrieves a user's profile by their ID.
// Returns a UserResponse DTO (without password).
func (s *authService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}
	if user == nil {
		return nil, errors.New(constants.MsgUserNotFound)
	}

	userResponse := &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	return userResponse, nil
}

// GetUserByID retrieves any user's profile by their ID (used by providers to view patient profiles).
func (s *authService) GetUserByID(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}
	if user == nil {
		return nil, errors.New(constants.MsgUserNotFound)
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateProfile updates user profile information (Name, Email, Phone, optional Password).
func (s *authService) UpdateProfile(userID uint, req dto.UpdateUserProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}
	if user == nil {
		return nil, errors.New(constants.MsgUserNotFound)
	}

	// Check email uniqueness if email is changed
	if req.Email != user.Email {
		existing, err := s.userRepo.FindByEmail(req.Email)
		if err != nil {
			return nil, errors.New(constants.MsgInternalServerError)
		}
		if existing != nil && existing.ID != userID {
			return nil, errors.New(constants.MsgEmailExists)
		}
	}

	// Check phone uniqueness if phone is changed
	if req.Phone != user.Phone {
		existing, err := s.userRepo.FindByPhone(req.Phone)
		if err != nil {
			return nil, errors.New(constants.MsgInternalServerError)
		}
		if existing != nil && existing.ID != userID {
			return nil, errors.New(constants.MsgPhoneExists)
		}
	}

	// Apply updates
	user.Name = req.Name
	user.Email = req.Email
	user.Phone = req.Phone

	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New(constants.MsgInternalServerError)
		}
		user.Password = hashed
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, errors.New(constants.MsgInternalServerError)
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  user.Role,
	}, nil
}
