package service

import (
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"context"
)

type TruckService interface {
	FindAll(ctx context.Context) ([]models.Truck, error)
	FindByID(ctx context.Context, id int) (*models.Truck, error)
	FindByNroTruck(ctx context.Context, nroTruck string) (*models.Truck, error)
	Create(ctx context.Context, truck *models.Truck) error
	Update(ctx context.Context, truck *models.Truck) error
	Delete(ctx context.Context, id int) error
}

type truckService struct {
	store store.TruckStore
}

func NewTruckService(store store.TruckStore) TruckService {
	return &truckService{store: store}
}

func (s *truckService) FindAll(ctx context.Context) ([]models.Truck, error) {
	return s.store.FindAll(ctx)
}

func (s *truckService) FindByID(ctx context.Context, id int) (*models.Truck, error) {
	return s.store.FindByID(ctx, id)
}

func (s *truckService) FindByNroTruck(ctx context.Context, nroTruck string) (*models.Truck, error) {
	return s.store.FindByNroTruck(ctx, nroTruck)
}

func (s *truckService) Create(ctx context.Context, truck *models.Truck) error {
	return s.store.Create(ctx, truck)
}

func (s *truckService) Update(ctx context.Context, truck *models.Truck) error {
	return s.store.Update(ctx, truck)
}

func (s *truckService) Delete(ctx context.Context, id int) error {
	return s.store.Delete(ctx, id)
}
