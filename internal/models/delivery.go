package models

import "time"

type Delivery struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	NroCta      string     `gorm:"not null" json:"nro_cta"`
	NroRto      string     `gorm:"not null" json:"nro_rto"`
	NroSerie    string     `json:"nro_serie"`
	Token       string     `gorm:"not null" json:"token"`
	Estado      string     `gorm:"not null" json:"estado"`
	TipoEntrega string     `gorm:"not null" json:"tipo_entrega"`
	FechaAccion CustomDate `json:"fecha_accion"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
