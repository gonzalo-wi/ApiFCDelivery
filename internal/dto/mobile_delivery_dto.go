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
	Valid      bool               `json:"valid"`
	Message    string             `json:"message"`
	Delivery   *DeliveryInfoDTO   `json:"delivery,omitempty"`
	Dispensers []DispenserInfoDTO `json:"dispensers,omitempty"`
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

// DispenserInfoDTO - Información de dispenser esperado
type DispenserInfoDTO struct {
	ID        int    `json:"id"`
	Marca     string `json:"marca"`
	NroSerie  string `json:"nro_serie"`
	Tipo      string `json:"tipo"`
	Validated bool   `json:"validated"`
}

// ValidateDispenserRequest - Validar código de dispenser escaneado
type ValidateDispenserRequest struct {
	DeliveryID int    `json:"delivery_id" binding:"required"`
	NroSerie   string `json:"nro_serie" binding:"required"`
}

// ValidateDispenserResponse - Respuesta de validación de dispenser
type ValidateDispenserResponse struct {
	Valid     bool              `json:"valid"`
	Message   string            `json:"message"`
	Dispenser *DispenserInfoDTO `json:"dispenser,omitempty"`
}

// MobileCompleteDeliveryRequest - Completar la entrega desde app móvil
type MobileCompleteDeliveryRequest struct {
	DeliveryID int      `json:"delivery_id" binding:"required"`
	Token      string   `json:"token" binding:"required"`
	Validated  []string `json:"validated_dispensers" binding:"required,min=1"`
}

// DispenserCompletedDTO - Información de dispenser en respuesta de completar entrega
type DispenserCompletedDTO struct {
	Marca    string `json:"marca"`
	NroSerie string `json:"nro_serie"`
}

// MobileCompleteDeliveryResponse - Respuesta al completar entrega desde app móvil
type MobileCompleteDeliveryResponse struct {
	NroCta          string                  `json:"nroCta"`
	Name            string                  `json:"name"`
	Address         string                  `json:"address"`
	Locality        string                  `json:"locality"`
	NroRto          string                  `json:"nroRto"`
	CreatedAt       string                  `json:"createdAt"`
	TipoAccion      string                  `json:"tipoAccion"`
	Token           string                  `json:"token"`
	Dispensers      []DispenserCompletedDTO `json:"dispensers"`
	WorkOrderQueued bool                    `json:"work_order_queued"`
}
