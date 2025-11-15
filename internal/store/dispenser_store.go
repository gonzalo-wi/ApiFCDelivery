package store

import (
	"GoFrioCalor/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type DispenserStore interface {
	FindAll(ctx context.Context) ([]models.Dispenser, error)
	FindByID(ctx context.Context, id int) (*models.Dispenser, error)
	Create(ctx context.Context, dispenser *models.Dispenser) error
	Update(ctx context.Context, dispenser *models.Dispenser) error
	Delete(ctx context.Context, id int) error
	FindByDeliveryID(ctx context.Context, deliveryID int) ([]models.Dispenser, error)
}

type dispenserStore struct {
	db *gorm.DB
}

func NewDispenserStore(db *gorm.DB) DispenserStore {
	return &dispenserStore{db: db}
}

func (s *dispenserStore) FindAll(ctx context.Context) ([]models.Dispenser, error) {
	var dispensers []models.Dispenser
	if err := s.db.WithContext(ctx).Find(&dispensers).Error; err != nil {
		return nil, fmt.Errorf("failed to find all dispensers: %w", err)
	}
	return dispensers, nil
}

func (s *dispenserStore) FindByID(ctx context.Context, id int) (*models.Dispenser, error) {
	var dispenser models.Dispenser
	if err := s.db.WithContext(ctx).First(&dispenser, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find dispenser with id %d: %w", id, err)
	}
	return &dispenser, nil
}

func (s *dispenserStore) FindByDeliveryID(ctx context.Context, deliveryID int) ([]models.Dispenser, error) {
	var dispensers []models.Dispenser
	if err := s.db.WithContext(ctx).Where("delivery_id = ?", deliveryID).Find(&dispensers).Error; err != nil {
		return nil, fmt.Errorf("failed to find dispensers for delivery %d: %w", deliveryID, err)
	}
	return dispensers, nil
}

func (s *dispenserStore) Create(ctx context.Context, dispenser *models.Dispenser) error {
	if err := s.db.WithContext(ctx).Create(dispenser).Error; err != nil {
		return fmt.Errorf("failed to create dispenser: %w", err)
	}
	return nil
}

func (s *dispenserStore) Update(ctx context.Context, dispenser *models.Dispenser) error {
	if err := s.db.WithContext(ctx).Save(dispenser).Error; err != nil {
		return fmt.Errorf("failed to update dispenser: %w", err)
	}
	return nil
}

func (s *dispenserStore) Delete(ctx context.Context, id int) error {
	if err := s.db.WithContext(ctx).Delete(&models.Dispenser{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete dispenser with id %d: %w", id, err)
	}
	return nil
}
