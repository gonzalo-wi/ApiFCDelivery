package models

import "time"

type WorkOrder struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	OrderNumber string    `gorm:"unique;not null" json:"order_number"`
	NroCta      string    `gorm:"not null" json:"nro_cta"`
	NroRto      string    `gorm:"not null" json:"nro_rto"`
	Name        string    `gorm:"not null" json:"name"`
	Address     string    `gorm:"not null" json:"address"`
	Localidad   string    `gorm:"not null" json:"localidad"`
	TipoAccion  string    `gorm:"not null" json:"tipo_accion"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
