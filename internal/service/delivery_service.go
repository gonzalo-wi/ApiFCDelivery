package service

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"fmt"
	"math/rand"
	"time"
)

type DeliveryService interface {
	FindAll() ([]models.Delivery, error)
	FindByID(id int) (*models.Delivery, error)
	Create(delivery *models.Delivery) error
	Update(delivery *models.Delivery) error
	Delete(id int) error
}
type deliveryService struct {
	store store.DeliveryStore
}

func NewDeliveryService(store store.DeliveryStore) DeliveryService {
	return &deliveryService{store: store}
}

func (s *deliveryService) FindAll() ([]models.Delivery, error) {
	return s.store.FindAll()
}
func (s *deliveryService) FindByID(id int) (*models.Delivery, error) {
	return s.store.FindByID(id)
}

func (s *deliveryService) Create(delivery *models.Delivery) error {
	delivery.Token = s.generateToken()
	if delivery.FechaAccion.IsZero() {
		delivery.FechaAccion = models.CustomDate{Time: time.Now()}
	}
	return s.store.Create(delivery)
}

func (s *deliveryService) Update(delivery *models.Delivery) error {
	return s.store.Update(delivery)
}
func (s *deliveryService) Delete(id int) error {
	return s.store.Delete(id)
}

func (s *deliveryService) generateToken() string {
	rand.Seed(time.Now().UnixNano())
	token := rand.Intn(constants.TOKEN_MAX-constants.TOKEN_MIN+1) + constants.TOKEN_MIN
	return fmt.Sprintf("%04d", token)
}
