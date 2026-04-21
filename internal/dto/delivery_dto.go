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
	ConversationID *string                 `json:"conversation_id,omitempty"`
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
		ConversationID: delivery.ConversationID,
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

// TallerPrepDeliveryItem representa un delivery con los datos relevantes para preparación de taller
type TallerPrepDeliveryItem struct {
	ID     int    `json:"id"`
	NroRto string `json:"nro_rto"`

	ItemDispensers []ItemDispenserResponse `json:"item_dispensers"`
}

// TallerPrepResponse es el resumen de dispensers a preparar para una fecha dada
type TallerPrepResponse struct {
	Fecha            string                   `json:"fecha"`
	TotalDispenserP  uint                     `json:"total_dispensers_P"`
	TotalDispenserM  uint                     `json:"total_dispensers_M"`
	TotalDispensers  uint                     `json:"total_dispensers"`
	CantidadEntregar int                      `json:"cantidad_deliveries"`
	Deliveries       []TallerPrepDeliveryItem `json:"deliveries"`
}

func ToTallerPrepResponse(fecha string, deliveries []models.Delivery) TallerPrepResponse {
	items := make([]TallerPrepDeliveryItem, 0, len(deliveries))
	var totalP, totalM uint

	for _, d := range deliveries {
		dispensers := make([]ItemDispenserResponse, len(d.ItemDispensers))
		for i, item := range d.ItemDispensers {
			dispensers[i] = ItemDispenserResponse{
				Tipo:     item.Tipo,
				Cantidad: item.Cantidad,
			}
			switch item.Tipo {
			case models.TipoDispenserPie:
				totalP += item.Cantidad
			case models.TipoDispenserMesada:
				totalM += item.Cantidad
			}
		}
		items = append(items, TallerPrepDeliveryItem{
			ID:             d.ID,
			NroRto:         d.NroRto,
			ItemDispensers: dispensers,
		})
	}

	return TallerPrepResponse{
		Fecha:            fecha,
		TotalDispenserP:  totalP,
		TotalDispenserM:  totalM,
		TotalDispensers:  totalP + totalM,
		CantidadEntregar: len(deliveries),
		Deliveries:       items,
	}
}
