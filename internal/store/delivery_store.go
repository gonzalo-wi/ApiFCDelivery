package store

import (
	"GoFrioCalor/internal/models"

	"gorm.io/gorm"
)

type DeliveryStore interface {
	FindAll() ([]models.Delivery, error)
	FindByID(id int) (*models.Delivery, error)
	Create(delivery *models.Delivery) error
	Update(delivery *models.Delivery) error
	Delete(id int) error
}

type deliveryStore struct {
	db *gorm.DB
}

func NewDeliveryStore(db *gorm.DB) DeliveryStore {
	return &deliveryStore{db: db}
}

func (s *deliveryStore) FindAll() ([]models.Delivery, error) {
	var deliveries []models.Delivery
	if err := s.db.Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}
func (s *deliveryStore) FindByID(id int) (*models.Delivery, error) {
	var delivery models.Delivery
	if err := s.db.First(&delivery, id).Error; err != nil {
		return nil, err
	}
	return &delivery, nil
}
func (s *deliveryStore) Create(delivery *models.Delivery) error {
	return s.db.Create(delivery).Error
}
func (s *deliveryStore) Update(delivery *models.Delivery) error {
	return s.db.Save(delivery).Error
}
func (s *deliveryStore) Delete(id int) error {
	return s.db.Delete(&models.Delivery{}, id).Error
}
