package dto

import "GoFrioCalor/internal/models"

// DispenserTypesCount cuenta la cantidad de cada tipo de dispenser
type DispenserTypesCount struct {
	P uint `json:"P,omitempty"`
	M uint `json:"M,omitempty"`
}

type DeliveryResponse struct {
	ID              int                  `json:"id,omitempty"`
	NroCta          string               `json:"nro_cta"`
	NroRto          string               `json:"nro_rto"`
	Dispensers      []DispenserResponse  `json:"dispensers"`
	TiposDispensers DispenserTypesCount  `json:"tipos_dispensers"`
	Cantidad        uint                 `json:"cantidad"`
	Token           string               `json:"token"`
	Estado          models.EstadoEntrega `json:"estado"`
	TipoEntrega     models.TipoEntrega   `json:"tipo_entrega"`
	EntregadoPor    models.EntregadoPor  `json:"entregado_por"`
	SessionID       string               `json:"session_id,omitempty"`
	FechaAccion     string               `json:"fecha_accion"`
}

func ToDeliveryResponse(delivery *models.Delivery) DeliveryResponse {
	dispensers := make([]DispenserResponse, len(delivery.Dispensers))
	// Contar tipos de dispensers
	tiposCount := DispenserTypesCount{}
	for i, d := range delivery.Dispensers {
		dispensers[i] = DispenserResponse{
			Marca:    d.Marca,
			NroSerie: d.NroSerie,
			Tipo:     d.Tipo,
		}
		switch d.Tipo {
		case models.TipoDispenserPie:
			tiposCount.P++
		case models.TipoDispenserMesada:
			tiposCount.M++
		}
	}

	return DeliveryResponse{
		ID:              delivery.ID,
		NroCta:          delivery.NroCta,
		NroRto:          delivery.NroRto,
		Dispensers:      dispensers,
		TiposDispensers: tiposCount,
		Cantidad:        delivery.Cantidad,
		Token:           delivery.Token,
		Estado:          delivery.Estado,
		TipoEntrega:     delivery.TipoEntrega,
		EntregadoPor:    delivery.EntregadoPor,
		SessionID:       delivery.SessionID,
		FechaAccion:     delivery.FechaAccion.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToDeliveryResponseList(deliveries []models.Delivery) []DeliveryResponse {
	responses := make([]DeliveryResponse, len(deliveries))
	for i, delivery := range deliveries {
		responses[i] = ToDeliveryResponse(&delivery)
	}
	return responses
}
