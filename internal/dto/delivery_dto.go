package dto

import "GoFrioCalor/internal/models"

// DeliveryResponse es el DTO para las respuestas de delivery (sin created_at y updated_at)
type DeliveryResponse struct {
	ID          int    `json:"id"`
	NroCta      string `json:"nro_cta"`
	NroRto      string `json:"nro_rto"`
	NroSerie    string `json:"nro_serie"`
	Token       string `json:"token"`
	Estado      string `json:"estado"`
	TipoEntrega string `json:"tipo_entrega"`
	FechaAccion string `json:"fecha_accion"`
}

// ToDeliveryResponse convierte un modelo Delivery a DeliveryResponse
func ToDeliveryResponse(delivery *models.Delivery) DeliveryResponse {
	return DeliveryResponse{
		ID:          delivery.ID,
		NroCta:      delivery.NroCta,
		NroRto:      delivery.NroRto,
		NroSerie:    delivery.NroSerie,
		Token:       delivery.Token,
		Estado:      delivery.Estado,
		TipoEntrega: delivery.TipoEntrega,
		FechaAccion: delivery.FechaAccion.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToDeliveryResponseList convierte un slice de Delivery a DeliveryResponse
func ToDeliveryResponseList(deliveries []models.Delivery) []DeliveryResponse {
	responses := make([]DeliveryResponse, len(deliveries))
	for i, delivery := range deliveries {
		responses[i] = ToDeliveryResponse(&delivery)
	}
	return responses
}
