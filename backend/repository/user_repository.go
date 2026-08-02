package repository

import (
	"backend/models"

	"gorm.io/gorm"
)

// UserRepository defines the contract for user-related database operations.
type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
}

// userRepository is the concrete implementation of UserRepository.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance with the given database connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser inserts a new user record into the database.
func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// UpdateUser persists changes to an existing user.
func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// FindByEmail retrieves a user by their email address.
// Returns nil (without error) if no user is found.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByPhone retrieves a user by their phone number.
// Returns nil (without error) if no user is found.
func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	result := r.db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByID retrieves a user by their primary key ID.
// Returns nil (without error) if no user is found.
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
