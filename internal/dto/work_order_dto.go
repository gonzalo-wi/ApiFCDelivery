package dto

type WorkOrderRequest struct {
	NroCta     string                      `json:"nroCta" binding:"required,min=1,max=50"`
	Name       string                      `json:"name" binding:"required,min=3,max=200"`
	Address    string                      `json:"address" binding:"required,min=5,max=300"`
	Locality   string                      `json:"locality" binding:"required,min=2,max=100"`
	NroRto     string                      `json:"nroRto" binding:"required,min=1,max=50"`
	CreatedAt  string                      `json:"createdAt" binding:"required"`
	Dispensers []WorkOrderDispenserRequest `json:"dispensers"`
	TipoAccion string                      `json:"tipoAccion" binding:"required,oneof=Instalacion Retiro Recambio"`
	Token      string                      `json:"token" binding:"omitempty,len=4,numeric"`
}

type WorkOrderDispenserRequest struct {
	Marca    string `json:"marca" binding:"required,min=2,max=50"`
	NroSerie string `json:"nro_serie" binding:"required,min=3,max=100"`
}
