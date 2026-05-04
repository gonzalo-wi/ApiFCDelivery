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
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
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

// DispenserOperation - Una operación sobre un dispenser (instalación, retiro, recambio o servicio técnico)
type DispenserOperation struct {
	Type                   string `json:"type" binding:"required,oneof=installation retirement replacement service"`
	InstalledDispenserCode string `json:"installed_dispenser_code,omitempty"`
	RetiredDispenserCode   string `json:"retired_dispenser_code,omitempty"`
	ServiceDispenserCode   string `json:"service_dispenser_code,omitempty"`
}

// MobileCompleteDeliveryRequest - Completar la entrega desde app móvil
// delivery_id y token son opcionales: solo requeridos para instalaciones pre-coordinadas (Infobip)
// Para retiros y recambios se crea un delivery nuevo en el momento
type MobileCompleteDeliveryRequest struct {
	DeliveryID  int                  `json:"delivery_id"`
	OrderNumber string               `json:"order_number" binding:"required"`
	NroCta      string               `json:"nro_cta" binding:"required"`
	NroRto      string               `json:"nro_rto"`
	Name        string               `json:"name"`
	Email       string               `json:"email"`
	Address     string               `json:"address"`
	Locality    string               `json:"locality"`
	Token       string               `json:"token"`
	Operations  []DispenserOperation `json:"operations" binding:"required,min=1,dive"`
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

// OperationCompletedDTO - Operación completada devuelta en la respuesta
type OperationCompletedDTO struct {
	Type                   string `json:"type"`
	InstalledDispenserCode string `json:"installed_dispenser_code,omitempty"`
	RetiredDispenserCode   string `json:"retired_dispenser_code,omitempty"`
	ServiceDispenserCode   string `json:"service_dispenser_code,omitempty"`
}

// MobileCompleteDeliveryResponse - Respuesta al completar entrega desde app móvil
type MobileCompleteDeliveryResponse struct {
	DeliveryID      int                     `json:"delivery_id"`
	NroCta          string                  `json:"nro_cta"`
	Name            string                  `json:"name"`
	Email           string                  `json:"email"`
	Address         string                  `json:"address"`
	Locality        string                  `json:"locality"`
	NroRto          string                  `json:"nro_rto"`
	TipoAccion      string                  `json:"tipo_accion"`
	OrderNumber     string                  `json:"order_number"`
	Operations      []OperationCompletedDTO `json:"operations"`
	WorkOrderQueued bool                    `json:"work_order_queued"`
}

// MobileDeliverySearchResponse - Respuesta simplificada para búsqueda de deliveries (mobile)
type MobileDeliverySearchResponse struct {
	ID          int    `json:"id"`
	FechaAccion string `json:"fecha_accion"`
	NroCta      string `json:"nro_cta"`
	Token       string `json:"token"`
}
