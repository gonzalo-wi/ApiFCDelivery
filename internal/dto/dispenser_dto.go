package dto

import "GoFrioCalor/internal/models"

type DispenserResponse struct {
	Marca    string               `json:"marca"`
	NroSerie string               `json:"nro_serie"`
	Tipo     models.TipoDispenser `json:"tipo"`
}
