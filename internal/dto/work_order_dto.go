package dto

type WorkOrderRequest struct {
	NroCta      string                      `json:"nroCta" binding:"required,min=1,max=50"`
	Name        string                      `json:"name" binding:"required,min=3,max=200"`
	Address     string                      `json:"address" binding:"required,min=5,max=300"`
	Locality    string                      `json:"locality" binding:"required,min=2,max=100"`
	NroRto      string                      `json:"nroRto" binding:"required,min=1,max=50"`
	CreatedAt   string                      `json:"createdAt" binding:"required"`
	AcceptedAt  string                      `json:"acceptedAt" binding:"omitempty"`
	Dispensers  []WorkOrderDispenserRequest `json:"dispensers"`
	TipoAccion  string                      `json:"tipoAccion" binding:"required,oneof=Instalacion Retiro Recambio"`
	Token       string                      `json:"token" binding:"omitempty,len=4,numeric"`
	OrderNumber string                      `json:"order_number" binding:"omitempty"`
}

type WorkOrderDispenserRequest struct {
	NroSerie string `json:"nro_serie" binding:"required,min=3,max=100"`
}
