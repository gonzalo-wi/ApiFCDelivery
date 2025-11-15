package store

import (
	"GoFrioCalor/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type WorkOrderStore interface {
	Create(ctx context.Context, workOrder *models.WorkOrder) error
	GetNextOrderNumber(ctx context.Context) (string, error)
}

type workOrderStore struct {
	db *gorm.DB
}

func NewWorkOrderStore(db *gorm.DB) WorkOrderStore {
	return &workOrderStore{db: db}
}

func (s *workOrderStore) Create(ctx context.Context, workOrder *models.WorkOrder) error {
	if err := s.db.WithContext(ctx).Create(workOrder).Error; err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}
	return nil
}

func (s *workOrderStore) GetNextOrderNumber(ctx context.Context) (string, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.WorkOrder{}).Count(&count).Error; err != nil {
		return "", fmt.Errorf("failed to count work orders: %w", err)
	}

	orderNumber := fmt.Sprintf("OT-%06d", count+1)
	return orderNumber, nil
}
