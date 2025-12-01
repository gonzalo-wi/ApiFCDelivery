package service

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type DeliveryService interface {
	FindAll(ctx context.Context) ([]models.Delivery, error)
	FindByID(ctx context.Context, id int) (*models.Delivery, error)
	FindByFilters(ctx context.Context, nroCta string, fechaAccion *time.Time) ([]models.Delivery, error)
	Create(ctx context.Context, delivery *models.Delivery) error
	Update(ctx context.Context, delivery *models.Delivery) error
	Delete(ctx context.Context, id int) error
}
type deliveryService struct {
	store store.DeliveryStore
}

func NewDeliveryService(store store.DeliveryStore) DeliveryService {
	return &deliveryService{store: store}
}

func (s *deliveryService) FindAll(ctx context.Context) ([]models.Delivery, error) {
	return s.store.FindAll(ctx)
}

func (s *deliveryService) FindByID(ctx context.Context, id int) (*models.Delivery, error) {
	return s.store.FindByID(ctx, id)
}

func (s *deliveryService) FindByFilters(ctx context.Context, nroCta string, fechaAccion *time.Time) ([]models.Delivery, error) {
	return s.store.FindByFilters(ctx, nroCta, fechaAccion)
}

// Al momento de crear la entrega se genera el token para el cliente
func (s *deliveryService) Create(ctx context.Context, delivery *models.Delivery) error {
	delivery.Token = s.generateToken()
	if delivery.FechaAccion.IsZero() {
		delivery.FechaAccion = models.CustomDate{Time: time.Now()}
	}
	return s.store.Create(ctx, delivery)
}

func (s *deliveryService) Update(ctx context.Context, delivery *models.Delivery) error {
	return s.store.Update(ctx, delivery)
}

func (s *deliveryService) Delete(ctx context.Context, id int) error {
	return s.store.Delete(ctx, id)
}

// generar token de 4 digitos
func (s *deliveryService) generateToken() string {
	rangeSize := int64(constants.TOKEN_MAX - constants.TOKEN_MIN + 1)
	n, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return fmt.Sprintf("%04d", time.Now().UnixNano()%10000)
	}
	token := int(n.Int64()) + constants.TOKEN_MIN
	return fmt.Sprintf("%04d", token)
}
