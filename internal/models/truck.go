package models

import "time"

type Truck struct {
	ID        int        `gorm:"primaryKey" json:"id,omitempty"`
	NroTruck  string     `gorm:"not null;unique" json:"nro_truck" binding:"required,min=1,max=20"`
	Deliverys []Delivery `gorm:"foreignKey:TruckID" json:"deliverys"`
	Capacity  int        `gorm:"not null" json:"capacity" binding:"required,gt=0"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
