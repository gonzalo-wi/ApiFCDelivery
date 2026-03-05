package dto

// ValidateTokenRequest - Solicitud del repartidor para validar el token del cliente
// Requiere token + nro_cta + fecha para mayor seguridad
type ValidateTokenRequest struct {
	Token       string `json:"token" binding:"required,min=4"`
	NroCta      string `json:"nro_cta" binding:"required"`
	FechaAccion string `json:"fecha_accion" binding:"required"`
}

// ValidateTokenResponse - Respuesta de validación de token
type ValidateTokenResponse struct {
	Valid          bool                   `json:"valid"`
	Message        string                 `json:"message"`
	Delivery       *DeliveryInfoDTO       `json:"delivery,omitempty"`
	ItemDispensers []ItemDispenserInfoDTO `json:"item_dispensers,omitempty"`
}

// DeliveryInfoDTO - Información básica del delivery para el repartidor
type DeliveryInfoDTO struct {
	ID          int    `json:"id"`
	NroCta      string `json:"nro_cta"`
	NroRto      string `json:"nro_rto"`
	Cantidad    uint   `json:"cantidad"`
	TipoEntrega string `json:"tipo_entrega"`
	FechaAccion string `json:"fecha_accion"`
}

// ItemDispenserInfoDTO - Información de items de dispensers en la entrega
type ItemDispenserInfoDTO struct {
	Tipo     string `json:"tipo"`
	Cantidad uint   `json:"cantidad"`
}

// Eliminadas las validaciones individuales de dispenser
// Ya no se escanean dispensers individuales, solo se registran cantidades por tipo

// MobileCompleteDeliveryRequest - Completar la entrega desde app móvil
type MobileCompleteDeliveryRequest struct {
	DeliveryID          int      `json:"delivery_id" binding:"required"`
	Name                string   `json:"name"`
	Email               string   `json:"email"`
	Address             string   `json:"address"`
	Locality            string   `json:"locality"`
	Token               string   `json:"token" binding:"required"`
	ValidatedDispensers []string `json:"validated_dispensers" binding:"required,min=1"`
}

// ItemDispenserDelivered - Items de dispensers efectivamente entregados
type ItemDispenserDelivered struct {
	Tipo     string `json:"tipo" binding:"required,oneof=P M"`
	Cantidad uint   `json:"cantidad" binding:"required,min=1"`
}

// ItemDispenserCompletedDTO - Información de items de dispensers en respuesta de completar entrega
type ItemDispenserCompletedDTO struct {
	Tipo     string `json:"tipo"`
	Cantidad uint   `json:"cantidad"`
}

// MobileCompleteDeliveryResponse - Respuesta al completar entrega desde app móvil
type MobileCompleteDeliveryResponse struct {
	NroCta              string                      `json:"nroCta"`
	Name                string                      `json:"name"`
	Email               string                      `json:"email"`
	Address             string                      `json:"address"`
	Locality            string                      `json:"locality"`
	NroRto              string                      `json:"nroRto"`
	CreatedAt           string                      `json:"createdAt"`
	TipoAccion          string                      `json:"tipoAccion"`
	Token               string                      `json:"token"`
	ItemDispensers      []ItemDispenserCompletedDTO `json:"item_dispensers"`
	ValidatedDispensers []string                    `json:"validated_dispensers,omitempty"`
	WorkOrderQueued     bool                        `json:"work_order_queued"`
}

// MobileDeliverySearchResponse - Respuesta simplificada para búsqueda de deliveries (mobile)
type MobileDeliverySearchResponse struct {
	FechaAccion string `json:"fecha_accion"`
	NroCta      string `json:"nro_cta"`
	Token       string `json:"token"`
}
