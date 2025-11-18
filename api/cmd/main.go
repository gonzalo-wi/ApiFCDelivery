package main

import (
	"GoFrioCalor/config"
	"GoFrioCalor/internal/routes"
	"GoFrioCalor/internal/service"
	"GoFrioCalor/internal/store"
	"GoFrioCalor/internal/transport"
	"log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error cargando configuración:", err)
	}

	db, err := config.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	deliveryStore := store.NewDeliveryStore(db)
	deliveryService := service.NewDeliveryService(deliveryStore)
	deliveryHandler := transport.NewDeliveryHandler(deliveryService)

	// Inicializar Dispenser Store, Service y Handler
	dispenserStore := store.NewDispenserStore(db)
	dispenserService := service.NewDispenserService(dispenserStore)
	dispenserHandler := transport.NewDispenserHandler(dispenserService)

	// Inicializar WorkOrder Store, PDF Service y Handler
	workOrderStore := store.NewWorkOrderStore(db)
	pdfService := service.NewPDFService(workOrderStore)
	workOrderHandler := transport.NewWorkOrderHandler(pdfService)

	router := routes.SetupRouter(deliveryHandler, dispenserHandler, workOrderHandler, cfg)

	log.Println("Base de datos conectada y tablas migradas correctamente")
	log.Printf("Servidor corriendo en http://localhost:%s\n", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Error iniciando el servidor:", err)
	}
}
