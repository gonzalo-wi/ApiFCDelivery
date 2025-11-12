package config

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(constants.MsgDatabaseConnectionError, err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(constants.MAX_IDLE_CONNS)
	sqlDB.SetMaxOpenConns(constants.MAX_OPEN_CONNS)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&models.Delivery{}, &models.Dispenser{}, &models.WorkOrder{}); err != nil {
		log.Println(constants.MsgInternalServerError, err)
		return nil, err
	}

	log.Println(constants.MsgDatabaseConnected)
	return db, nil
}
