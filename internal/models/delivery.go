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

type EntregadoPor string

const (
	Repartidor EntregadoPor = "Repartidor"
	Tecnico    EntregadoPor = "Tecnico"
)

type Delivery struct {
	ID             int           `gorm:"primaryKey" json:"id"`
	NroCta         string        `gorm:"not null" json:"nro_cta" binding:"required,min=1,max=50"`
	Name           string        `gorm:"type:varchar(200)" json:"name"`
	Address        string        `gorm:"type:varchar(300)" json:"address"`
	Locality       string        `gorm:"type:varchar(100)" json:"locality"`
	NroRto         string        `gorm:"not null" json:"nro_rto" binding:"required,min=1,max=50"`
	Dispensers     []Dispenser   `gorm:"foreignKey:DeliveryID" json:"dispensers"`
	Cantidad       uint          `gorm:"not null" json:"cantidad" binding:"required,min=1,max=3"`
	Token          string        `gorm:"not null" json:"token"`
	Estado         EstadoEntrega `gorm:"not null" json:"estado" binding:"required,oneof=Pendiente Completado Cancelado"`
	TipoEntrega    TipoEntrega   `gorm:"not null" json:"tipo_entrega" binding:"required,oneof=Instalacion Retiro Recambio"`
	EntregadoPor   EntregadoPor  `gorm:"not null" json:"entregado_por" binding:"required,oneof=Repartidor Tecnico"`
	SessionID      string        `gorm:"index:idx_session_id,unique" json:"session_id,omitempty"`
	TermsSessionID *int64        `gorm:"index" json:"terms_session_id,omitempty"`
	FechaAccion    CustomDate    `json:"fecha_accion"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
