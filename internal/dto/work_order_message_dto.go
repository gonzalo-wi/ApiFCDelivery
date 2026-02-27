package dto

// WorkOrderMessageDTO - Mensaje que se publica en RabbitMQ para crear orden de trabajo
type WorkOrderMessageDTO struct {
	NroCta     string             `json:"nroCta"`
	Name       string             `json:"name"`
	Address    string             `json:"address"`
	Locality   string             `json:"locality"`
	NroRto     string             `json:"nroRto"`
	CreatedAt  string             `json:"createdAt"`
	TipoAccion string             `json:"tipoAccion"`
	Token      string             `json:"token"`
	Dispensers []DispenserMessage `json:"dispensers"`
	DeliveryID int                `json:"deliveryId"` // Para actualizar después
}

// DispenserMessage - Información de dispenser en el mensaje
type DispenserMessage struct {
	Marca    string `json:"marca"`
	NroSerie string `json:"nro_serie"`
}
