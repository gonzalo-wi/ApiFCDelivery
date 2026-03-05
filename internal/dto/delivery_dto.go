package dto

import "GoFrioCalor/internal/models"

type ItemDispenserResponse struct {
	Tipo     models.TipoDispenser `json:"tipo"`
	Cantidad uint                 `json:"cantidad"`
}

type DeliveryResponse struct {
	ID             int                     `json:"id,omitempty"`
	NroCta         string                  `json:"nro_cta"`
	NroRto         string                  `json:"nro_rto"`
	ItemDispensers []ItemDispenserResponse `json:"item_dispensers"`
	Cantidad       uint                    `json:"cantidad"`
	Token          string                  `json:"token"`
	Estado         models.EstadoEntrega    `json:"estado"`
	TipoEntrega    models.TipoEntrega      `json:"tipo_entrega"`
	EntregadoPor   models.EntregadoPor     `json:"entregado_por"`
	SessionID      *string                 `json:"session_id,omitempty"`
	FechaAccion    string                  `json:"fecha_accion"`
}

func ToDeliveryResponse(delivery *models.Delivery) DeliveryResponse {
	itemDispensers := make([]ItemDispenserResponse, len(delivery.ItemDispensers))
	for i, item := range delivery.ItemDispensers {
		itemDispensers[i] = ItemDispenserResponse{
			Tipo:     item.Tipo,
			Cantidad: item.Cantidad,
		}
	}

	return DeliveryResponse{
		ID:             delivery.ID,
		NroCta:         delivery.NroCta,
		NroRto:         delivery.NroRto,
		ItemDispensers: itemDispensers,
		Cantidad:       delivery.Cantidad,
		Token:          delivery.Token,
		Estado:         delivery.Estado,
		TipoEntrega:    delivery.TipoEntrega,
		EntregadoPor:   delivery.EntregadoPor,
		SessionID:      delivery.SessionID,
		FechaAccion:    delivery.FechaAccion.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToDeliveryResponseList(deliveries []models.Delivery) []DeliveryResponse {
	responses := make([]DeliveryResponse, len(deliveries))
	for i, delivery := range deliveries {
		responses[i] = ToDeliveryResponse(&delivery)
	}
	return responses
}
