package repository

import (
	"backend/models"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	Create(apt *models.Appointment) error
	GetByID(id uint) (*models.Appointment, error)
	GetByCustomerID(customerID uint) ([]models.Appointment, error)
	GetByProviderID(providerID uint) ([]models.Appointment, error)
	UpdateStatus(id uint, status string) error
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) Create(apt *models.Appointment) error {
	return r.db.Create(apt).Error
}

func (r *appointmentRepository) GetByID(id uint) (*models.Appointment, error) {
	var apt models.Appointment
	err := r.db.Preload("Customer").Preload("Provider").First(&apt, id).Error
	if err != nil {
		return nil, err
	}
	r.db.Where("user_id = ?", apt.ProviderID).First(&apt.Profile)
	return &apt, nil
}

func (r *appointmentRepository) GetByCustomerID(customerID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Where("customer_id = ?", customerID).
		Preload("Provider").
		Order("created_at desc").
		Find(&appointments).Error

	for i := range appointments {
		r.db.Where("user_id = ?", appointments[i].ProviderID).First(&appointments[i].Profile)
	}
	return appointments, err
}

func (r *appointmentRepository) GetByProviderID(providerID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	var profile models.ServiceProviderProfile
	r.db.Where("user_id = ?", providerID).First(&profile)

	query := r.db.Where("provider_id = ?", providerID)
	if profile.ID > 0 {
		query = r.db.Where("provider_id = ? OR provider_id = ?", providerID, profile.ID)
	}

	err := query.Preload("Customer").
		Order("created_at desc").
		Find(&appointments).Error

	for i := range appointments {
		r.db.Where("user_id = ?", appointments[i].ProviderID).First(&appointments[i].Profile)
	}
	return appointments, err
}

func (r *appointmentRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Appointment{}).Where("id = ?", id).Update("status", status).Error
}
