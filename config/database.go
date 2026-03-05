package config

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

	if err := db.AutoMigrate(&models.Delivery{}, &models.ItemDispenser{}, &models.WorkOrder{}, &models.TermsSession{}); err != nil {
		log.Println(constants.MsgInternalServerError, err)
		return nil, err
	}

	log.Println(constants.MsgDatabaseConnected)
	return db, nil
}

// NewSQLXDatabase crea una conexión de base de datos usando sqlx (para audit store)
func NewSQLXDatabase(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Println(constants.MsgDatabaseConnectionError, err)
		return nil, err
	}

	db.SetMaxIdleConns(constants.MAX_IDLE_CONNS)
	db.SetMaxOpenConns(constants.MAX_OPEN_CONNS)
	db.SetConnMaxLifetime(time.Hour)

	log.Println("SQLX database connection successful")
	return db, nil
}
