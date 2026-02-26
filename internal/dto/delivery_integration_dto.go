package dto

import "GoFrioCalor/internal/models"

type InitiateDeliveryRequest struct {
	NroCta      string             `json:"nro_cta" binding:"required,min=1,max=50"`
	NroRto      string             `json:"nro_rto" binding:"required,min=1,max=50"`
	Dispensers  []DispenserRequest `json:"dispensers" binding:"required,dive"`
	Cantidad    uint               `json:"cantidad" binding:"required,min=1,max=3"`
	TipoEntrega models.TipoEntrega `json:"tipo_entrega" binding:"required,oneof=Instalacion Retiro Recambio"`
	FechaAccion string             `json:"fecha_accion,omitempty"`
}

type InitiateDeliveryResponse struct {
	Token     string `json:"token"`
	TermsURL  string `json:"terms_url"`
	ExpiresAt string `json:"expires_at"`
	Message   string `json:"message"`
}

type CompleteDeliveryRequest struct {
}

type CompleteDeliveryResponse struct {
	Success  bool              `json:"success"`
	Message  string            `json:"message"`
	Delivery *DeliveryResponse `json:"delivery,omitempty"`
}

type DispenserRequest struct {
	Marca    string               `json:"marca" binding:"required,min=1,max=50"`
	NroSerie string               `json:"nro_serie" binding:"required,min=1,max=50"`
	Tipo     models.TipoDispenser `json:"tipo" binding:"required,oneof=A B C HELADERA P M"`
}

// InfobipDeliveryRequest es el request que envía el chatbot de Infobip para crear una entrega
type InfobipDeliveryRequest struct {
	NroCta       string                 `json:"nro_cta" binding:"required,min=1,max=50"`
	NroRto       string                 `json:"nro_rto" binding:"required,min=1,max=50"`
	Tipos        DispenserTypesQuantity `json:"tipos" binding:"required"`
	TipoEntrega  models.TipoEntrega     `json:"tipo_entrega" binding:"required,oneof=Instalacion Retiro Recambio"`
	EntregadoPor models.EntregadoPor    `json:"entregado_por" binding:"required,oneof=Repartidor Tecnico"`
	SessionID    string                 `json:"session_id" binding:"required,min=1"`
	FechaAccion  string                 `json:"fecha_accion,omitempty"`
}

// DispenserTypesQuantity especifica la cantidad de dispensers por tipo
type DispenserTypesQuantity struct {
	P uint `json:"P"` // Dispensers de Pie
	M uint `json:"M"` // Dispensers de Mesada
}

// InfobipDeliveryResponse es la respuesta que se envía al chatbot de Infobip
type InfobipDeliveryResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
