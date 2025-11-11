package store

import (
	"GoFrioCalor/internal/models"

	"gorm.io/gorm"
)

type DispenserStore interface {
	FindAll() ([]models.Dispenser, error)
	FindByID(id int) (*models.Dispenser, error)
	Create(dispenser *models.Dispenser) error
	Update(dispenser *models.Dispenser) error
	Delete(id int) error
	FindByDeliveryID(deliveryID int) ([]models.Dispenser, error)
}

type dispenserStore struct {
	db *gorm.DB
}

func NewDispenserStore(db *gorm.DB) DispenserStore {
	return &dispenserStore{db: db}
}

func (s *dispenserStore) FindAll() ([]models.Dispenser, error) {
	var dispensers []models.Dispenser
	if err := s.db.Find(&dispensers).Error; err != nil {
		return nil, err
	}
	return dispensers, nil
}

func (s *dispenserStore) FindByID(id int) (*models.Dispenser, error) {
	var dispenser models.Dispenser
	if err := s.db.First(&dispenser, id).Error; err != nil {
		return nil, err
	}
	return &dispenser, nil
}

func (s *dispenserStore) FindByDeliveryID(deliveryID int) ([]models.Dispenser, error) {
	var dispensers []models.Dispenser
	if err := s.db.Where("delivery_id = ?", deliveryID).Find(&dispensers).Error; err != nil {
		return nil, err
	}
	return dispensers, nil
}

func (s *dispenserStore) Create(dispenser *models.Dispenser) error {
	return s.db.Create(dispenser).Error
}

func (s *dispenserStore) Update(dispenser *models.Dispenser) error {
	return s.db.Save(dispenser).Error
}

func (s *dispenserStore) Delete(id int) error {
	return s.db.Delete(&models.Dispenser{}, id).Error
}
