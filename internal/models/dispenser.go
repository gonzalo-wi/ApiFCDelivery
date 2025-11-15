package models

import "time"

type Dispenser struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	Marca      string    `gorm:"not null" json:"marca" binding:"required,min=2,max=50"`
	NroSerie   string    `gorm:"not null" json:"nro_serie" binding:"required,min=3,max=100"`
	DeliveryID int       `gorm:"not null" json:"delivery_id" binding:"required,gt=0"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
