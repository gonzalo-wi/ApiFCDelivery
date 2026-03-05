package models

import "time"

type TipoDispenser string

const (
	TipoDispenserPie    TipoDispenser = "P"
	TipoDispenserMesada TipoDispenser = "M"
)

type ItemDispenser struct {
	ID         int           `gorm:"primaryKey" json:"id,omitempty"`
	Tipo       TipoDispenser `gorm:"not null" json:"tipo" binding:"required,oneof=P M"`
	Cantidad   uint          `gorm:"not null" json:"cantidad" binding:"required,min=1"`
	DeliveryID int           `gorm:"not null" json:"delivery_id"`
	CreatedAt  time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
