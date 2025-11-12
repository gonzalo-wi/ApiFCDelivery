package dto

import "GoFrioCalor/internal/models"

type DeliveryResponse struct {
	NroCta      string              `json:"nro_cta"`
	NroRto      string              `json:"nro_rto"`
	Dispensers  []DispenserResponse `json:"dispensers"`
	Token       string              `json:"token"`
	Estado      string              `json:"estado"`
	TipoEntrega string              `json:"tipo_entrega"`
	FechaAccion string              `json:"fecha_accion"`
}

func ToDeliveryResponse(delivery *models.Delivery) DeliveryResponse {
	dispensers := make([]DispenserResponse, len(delivery.Dispensers))
	for i, d := range delivery.Dispensers {
		dispensers[i] = DispenserResponse{
			Marca:    d.Marca,
			NroSerie: d.NroSerie,
		}
	}

	return DeliveryResponse{
		NroCta:      delivery.NroCta,
		NroRto:      delivery.NroRto,
		Dispensers:  dispensers,
		Token:       delivery.Token,
		Estado:      delivery.Estado,
		TipoEntrega: delivery.TipoEntrega,
		FechaAccion: delivery.FechaAccion.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToDeliveryResponseList(deliveries []models.Delivery) []DeliveryResponse {
	responses := make([]DeliveryResponse, len(deliveries))
	for i, delivery := range deliveries {
		responses[i] = ToDeliveryResponse(&delivery)
	}
	return responses
}
