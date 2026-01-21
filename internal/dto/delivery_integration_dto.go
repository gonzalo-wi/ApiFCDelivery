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
	Tipo     models.TipoDispenser `json:"tipo" binding:"required,oneof=A B C HELADERA"`
}
