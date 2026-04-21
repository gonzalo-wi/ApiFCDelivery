package dto

// OperationMessage - Operación individual incluida en el mensaje de OT
type OperationMessage struct {
	Type                   string `json:"type"`
	InstalledDispenserCode string `json:"installed_dispenser_code,omitempty"`
	RetiredDispenserCode   string `json:"retired_dispenser_code,omitempty"`
}

// WorkOrderMessageDTO - Mensaje que se publica en RabbitMQ para crear orden de trabajo
type WorkOrderMessageDTO struct {
	OrderNumber string             `json:"order_number"` // Número de orden proporcionado por app móvil
	NroCta      string             `json:"nroCta"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Address     string             `json:"address"`
	Locality    string             `json:"locality"`
	NroRto      string             `json:"nroRto"`
	CreatedAt   string             `json:"createdAt"`
	TipoAccion  string             `json:"tipoAccion"`
	Token       string             `json:"token"`
	Operations  []OperationMessage `json:"operations"`
	DeliveryID  int                `json:"deliveryId"` // Para actualizar después
}
