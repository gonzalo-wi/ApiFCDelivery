package store

import (
	"GoFrioCalor/internal/models"
	"context"

	"gorm.io/gorm"
)

type TruckStore interface {
	FindAll(ctx context.Context) ([]models.Truck, error)
	FindByID(ctx context.Context, id int) (*models.Truck, error)
	FindByNroTruck(ctx context.Context, nroTruck string) (*models.Truck, error)
	Create(ctx context.Context, truck *models.Truck) error
	Update(ctx context.Context, truck *models.Truck) error
	Delete(ctx context.Context, id int) error
}

type truckStore struct {
	db *gorm.DB
}

func NewTruckStore(db *gorm.DB) TruckStore {
	return &truckStore{db: db}
}

func (s *truckStore) FindAll(ctx context.Context) ([]models.Truck, error) {
	var trucks []models.Truck
	if err := s.db.WithContext(ctx).Find(&trucks).Error; err != nil {
		return nil, err
	}
	return trucks, nil
}

func (s *truckStore) FindByID(ctx context.Context, id int) (*models.Truck, error) {
	var truck models.Truck
	if err := s.db.WithContext(ctx).First(&truck, id).Error; err != nil {
		return nil, err
	}
	return &truck, nil
}

func (s *truckStore) FindByNroTruck(ctx context.Context, nroTruck string) (*models.Truck, error) {
	var truck models.Truck
	if err := s.db.WithContext(ctx).Where("nro_truck = ?", nroTruck).First(&truck).Error; err != nil {
		return nil, err
	}
	return &truck, nil
}

func (s *truckStore) Create(ctx context.Context, truck *models.Truck) error {
	if err := s.db.WithContext(ctx).Create(truck).Error; err != nil {
		return err
	}
	return nil
}

func (s *truckStore) Update(ctx context.Context, truck *models.Truck) error {
	if err := s.db.WithContext(ctx).Save(truck).Error; err != nil {
		return err
	}
	return nil
}

func (s *truckStore) Delete(ctx context.Context, id int) error {
	if err := s.db.WithContext(ctx).Delete(&models.Truck{}, id).Error; err != nil {
		return err
	}
	return nil
}
