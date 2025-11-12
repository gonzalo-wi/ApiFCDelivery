package dto

type WorkOrderRequest struct {
	NroCta     string                      `json:"nroCta" binding:"required"`
	Name       string                      `json:"name" binding:"required"`
	Address    string                      `json:"address" binding:"required"`
	Locality   string                      `json:"locality" binding:"required"`
	NroRto     string                      `json:"nroRto" binding:"required"`
	CreatedAt  string                      `json:"createdAt" binding:"required"`
	Dispensers []WorkOrderDispenserRequest `json:"dispensers"`
	TipoAccion string                      `json:"tipoAccion" binding:"required"`
}

type WorkOrderDispenserRequest struct {
	Marca    string `json:"marca"`
	NroSerie string `json:"nro_serie"`
}
