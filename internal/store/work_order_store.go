package store

import (
	"GoFrioCalor/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type WorkOrderStore interface {
	Create(workOrder *models.WorkOrder) error
	GetNextOrderNumber() (string, error)
}

type workOrderStore struct {
	db *gorm.DB
}

func NewWorkOrderStore(db *gorm.DB) WorkOrderStore {
	return &workOrderStore{db: db}
}

func (s *workOrderStore) Create(workOrder *models.WorkOrder) error {
	return s.db.Create(workOrder).Error
}

func (s *workOrderStore) GetNextOrderNumber() (string, error) {
	var count int64
	if err := s.db.Model(&models.WorkOrder{}).Count(&count).Error; err != nil {
		return "", err
	}

	orderNumber := fmt.Sprintf("OT-%06d", count+1)
	return orderNumber, nil
}
