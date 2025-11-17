package store

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DeliveryStore interface {
	FindAll(ctx context.Context) ([]models.Delivery, error)
	FindByID(ctx context.Context, id int) (*models.Delivery, error)
	FindByFilters(ctx context.Context, nroCta string, fechaAccion *time.Time) ([]models.Delivery, error)
	Create(ctx context.Context, delivery *models.Delivery) error
	Update(ctx context.Context, delivery *models.Delivery) error
	Delete(ctx context.Context, id int) error
}

type deliveryStore struct {
	db *gorm.DB
}

func NewDeliveryStore(db *gorm.DB) DeliveryStore {
	return &deliveryStore{db: db}
}

func (s *deliveryStore) FindAll(ctx context.Context) ([]models.Delivery, error) {
	var deliveries []models.Delivery
	if err := s.db.WithContext(ctx).Preload("Dispensers").Find(&deliveries).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFindAllDeliveries, err)
	}
	return deliveries, nil
}
func (s *deliveryStore) FindByID(ctx context.Context, id int) (*models.Delivery, error) {
	var delivery models.Delivery
	if err := s.db.WithContext(ctx).Preload("Dispensers").First(&delivery, id).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFindDeliveryByID, id, err)
	}
	return &delivery, nil
}

func (s *deliveryStore) FindByFilters(ctx context.Context, nroCta string, fechaAccion *time.Time) ([]models.Delivery, error) {
	var deliveries []models.Delivery
	query := s.db.WithContext(ctx).Preload("Dispensers")
	if nroCta != "" {
		query = query.Where("nro_cta = ?", nroCta)
	}
	if fechaAccion != nil {
		startOfDay := time.Date(fechaAccion.Year(), fechaAccion.Month(), fechaAccion.Day(), 0, 0, 0, 0, fechaAccion.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("fecha_accion >= ? AND fecha_accion < ?", startOfDay, endOfDay)
	}

	if err := query.Find(&deliveries).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFindDeliveriesFilters, err)
	}
	return deliveries, nil
}

func (s *deliveryStore) Create(ctx context.Context, delivery *models.Delivery) error {
	if err := s.db.WithContext(ctx).Create(delivery).Error; err != nil {
		return fmt.Errorf(constants.ErrCreateDelivery, err)
	}
	return nil
}

func (s *deliveryStore) Update(ctx context.Context, delivery *models.Delivery) error {
	if err := s.db.WithContext(ctx).Save(delivery).Error; err != nil {
		return fmt.Errorf(constants.ErrUpdateDelivery, err)
	}
	return nil
}

func (s *deliveryStore) Delete(ctx context.Context, id int) error {
	if err := s.db.WithContext(ctx).Delete(&models.Delivery{}, id).Error; err != nil {
		return fmt.Errorf(constants.ErrDeleteDelivery, id, err)
	}
	return nil
}
