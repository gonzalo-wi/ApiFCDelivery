package models

import "time"

type TipoEntrega string

const (
	Instalacion TipoEntrega = "Instalacion"
	Retiro      TipoEntrega = "Retiro"
	Recambio    TipoEntrega = "Recambio"
)

type EstadoEntrega string

const (
	Pendiente  EstadoEntrega = "Pendiente"
	Completado EstadoEntrega = "Completado"
	Cancelado  EstadoEntrega = "Cancelado"
)

type Delivery struct {
	ID          int           `gorm:"primaryKey" json:"id"`
	NroCta      string        `gorm:"not null" json:"nro_cta" binding:"required,min=1,max=50"`
	NroRto      string        `gorm:"not null" json:"nro_rto" binding:"required,min=1,max=50"`
	Dispensers  []Dispenser   `gorm:"foreignKey:DeliveryID" json:"dispensers"`
	Token       string        `gorm:"not null" json:"token"`
	Estado      EstadoEntrega `gorm:"not null" json:"estado" binding:"required,oneof=Pendiente Completado Cancelado"`
	TipoEntrega TipoEntrega   `gorm:"not null" json:"tipo_entrega" binding:"required,oneof=Instalacion Retiro Recambio"`
	FechaAccion CustomDate    `json:"fecha_accion"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
