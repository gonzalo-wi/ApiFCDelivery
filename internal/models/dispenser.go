package models

import "time"

type Dispenser struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	Marca      string    `gorm:"not null" json:"marca"`
	NroSerie   string    `gorm:"not null" json:"nro_serie"`
	DeliveryID int       `gorm:"not null" json:"delivery_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
