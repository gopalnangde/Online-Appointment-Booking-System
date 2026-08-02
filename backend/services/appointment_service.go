package services

import (
	"errors"

	"backend/dto"
	"backend/models"
	"backend/repository"
)

type AppointmentService interface {
	CreateAppointment(customerID uint, req *dto.CreateAppointmentRequest) (*models.Appointment, error)
	GetCustomerAppointments(customerID uint) ([]models.Appointment, error)
	GetProviderAppointments(providerID uint) ([]models.Appointment, error)
	UpdateStatus(userID uint, appointmentID uint, req *dto.UpdateAppointmentStatusRequest) error
}

type appointmentService struct {
	repo repository.AppointmentRepository
}

func NewAppointmentService(repo repository.AppointmentRepository) AppointmentService {
	return &appointmentService{repo: repo}
}

func (s *appointmentService) CreateAppointment(customerID uint, req *dto.CreateAppointmentRequest) (*models.Appointment, error) {
	apt := &models.Appointment{
		CustomerID:   customerID,
		ProviderID:   req.ProviderID,
		ServiceTitle: req.ServiceTitle,
		BookingDate:  req.BookingDate,
		BookingTime:  req.BookingTime,
		Notes:        req.Notes,
		Status:       "Pending",
	}

	err := s.repo.Create(apt)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(apt.ID)
}

func (s *appointmentService) GetCustomerAppointments(customerID uint) ([]models.Appointment, error) {
	return s.repo.GetByCustomerID(customerID)
}

func (s *appointmentService) GetProviderAppointments(providerID uint) ([]models.Appointment, error) {
	return s.repo.GetByProviderID(providerID)
}

func (s *appointmentService) UpdateStatus(userID uint, appointmentID uint, req *dto.UpdateAppointmentStatusRequest) error {
	apt, err := s.repo.GetByID(appointmentID)
	if err != nil {
		return errors.New("appointment not found")
	}

	// Validate authorization: customer, direct provider user ID, or matching profile user ID
	isAuthorized := apt.CustomerID == userID || apt.ProviderID == userID || (apt.Profile.UserID > 0 && apt.Profile.UserID == userID)
	if !isAuthorized {
		return errors.New("unauthorized to update this appointment")
	}

	validStatuses := map[string]bool{
		"Pending":   true,
		"Confirmed": true,
		"Completed": true,
		"Cancelled": true,
	}

	if !validStatuses[req.Status] {
		return errors.New("invalid status value")
	}

	return s.repo.UpdateStatus(appointmentID, req.Status)
}
