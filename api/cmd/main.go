package main

import (
	"GoFrioCalor/config"
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/routes"
	"GoFrioCalor/internal/service"
	"GoFrioCalor/internal/store"
	"GoFrioCalor/internal/transport"

	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg(constants.MsgErrorLoadingConfig)
	}

	config.InitLogger(cfg.Environment)

	db, err := config.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal().Err(err).Msg(constants.MsgErrorConnectingDB)
	}

	log.Info().Msg(constants.MsgDBConnectedSuccess)

	deliveryStore := store.NewDeliveryStore(db)
	deliveryService := service.NewDeliveryService(deliveryStore)
	deliveryHandler := transport.NewDeliveryHandler(deliveryService)

	dispenserStore := store.NewDispenserStore(db)
	dispenserService := service.NewDispenserService(dispenserStore)
	dispenserHandler := transport.NewDispenserHandler(dispenserService)

	workOrderStore := store.NewWorkOrderStore(db)
	pdfService := service.NewPDFService(workOrderStore)
	workOrderHandler := transport.NewWorkOrderHandler(pdfService)

	router := routes.SetupRouter(deliveryHandler, dispenserHandler, workOrderHandler, cfg)

	log.Info().Str("port", cfg.Port).Msgf(constants.MsgServerRunning, cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg(constants.MsgErrorStartingServer)
	}
}
