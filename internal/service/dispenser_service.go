package service

import (
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"context"
)

type DispenserService interface {
	FindAll(ctx context.Context) ([]models.Dispenser, error)
	FindByID(ctx context.Context, id int) (*models.Dispenser, error)
	Create(ctx context.Context, dispenser *models.Dispenser) error
	Update(ctx context.Context, dispenser *models.Dispenser) error
	Delete(ctx context.Context, id int) error
}

type dispenserService struct {
	store store.DispenserStore
}

func NewDispenserService(store store.DispenserStore) DispenserService {
	return &dispenserService{store: store}
}

func (s *dispenserService) FindAll(ctx context.Context) ([]models.Dispenser, error) {
	return s.store.FindAll(ctx)
}

func (s *dispenserService) FindByID(ctx context.Context, id int) (*models.Dispenser, error) {
	return s.store.FindByID(ctx, id)
}

func (s *dispenserService) Create(ctx context.Context, dispenser *models.Dispenser) error {
	return s.store.Create(ctx, dispenser)
}

func (s *dispenserService) Update(ctx context.Context, dispenser *models.Dispenser) error {
	return s.store.Update(ctx, dispenser)
}

func (s *dispenserService) Delete(ctx context.Context, id int) error {
	return s.store.Delete(ctx, id)
}
