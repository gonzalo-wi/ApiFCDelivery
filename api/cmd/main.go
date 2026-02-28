package main

import (
	"GoFrioCalor/config"
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/routes"
	"GoFrioCalor/internal/service"
	"GoFrioCalor/internal/store"
	"GoFrioCalor/internal/transport"
	"GoFrioCalor/migrations"
	"context"

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

	// Run database migrations
	migrationService := service.NewMigrationService(db)
	if err := migrationService.RunMigrations(migrations.FS); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	log.Info().Msg("Database migrations completed successfully")

	// Stores
	deliveryStore := store.NewDeliveryStore(db)
	dispenserStore := store.NewDispenserStore(db)
	workOrderStore := store.NewWorkOrderStore(db)
	termsSessionStore := store.NewTermsSessionStore(db)

	// RabbitMQ Configuration
	rabbitConfig := config.LoadRabbitMQConfig()

	// RabbitMQ Publisher
	rabbitPublisher, err := service.NewRabbitMQPublisher(rabbitConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize RabbitMQ Publisher, continuing without it")
		rabbitPublisher = nil
	} else {
		defer rabbitPublisher.Close()
		log.Info().Msg("RabbitMQ Publisher initialized successfully")
	}

	// Services
	deliveryService := service.NewDeliveryService(deliveryStore)
	deliveryHandler := transport.NewDeliveryHandler(deliveryService)

	dispenserService := service.NewDispenserService(dispenserStore)
	dispenserHandler := transport.NewDispenserHandler(dispenserService)

	pdfService := service.NewPDFService(workOrderStore)
	workOrderHandler := transport.NewWorkOrderHandler(pdfService)

	// Términos y Condiciones con Infobip
	infobipClient := service.NewInfobipClient(cfg.InfobipBaseURL, cfg.InfobipAPIKey)
	termsSessionService := service.NewTermsSessionService(termsSessionStore, infobipClient)
	termsSessionHandler := transport.NewTermsSessionHandler(termsSessionService, cfg.AppBaseURL, cfg.TermsTTLHours)

	// Flujo integrado: Entregas con Términos y Condiciones
	deliveryWithTermsService := service.NewDeliveryWithTermsService(deliveryStore, termsSessionStore, termsSessionService)
	deliveryWithTermsHandler := transport.NewDeliveryWithTermsHandler(deliveryWithTermsService, cfg.AppBaseURL, cfg.TermsTTLHours)

	// Mobile Delivery - Validación y Completar Entregas
	var mobileDeliveryHandler *transport.MobileDeliveryHandler
	if rabbitPublisher != nil {
		mobileDeliveryService := service.NewMobileDeliveryService(deliveryStore, rabbitPublisher)
		mobileDeliveryHandler = transport.NewMobileDeliveryHandler(mobileDeliveryService)
		log.Info().Msg("Mobile Delivery Service initialized")
	}

	// RabbitMQ Consumer para Work Orders
	if rabbitPublisher != nil {
		mockPDFService := service.NewMockPDFService()
		mockEmailService := service.NewMockEmailService()

		consumer, err := service.NewWorkOrderConsumer(
			rabbitConfig,
			workOrderStore,
			deliveryStore,
			mockPDFService,
			mockEmailService,
		)

		if err != nil {
			log.Error().Err(err).Msg("Failed to initialize Work Order Consumer")
		} else {
			ctx := context.Background()
			if err := consumer.Start(ctx); err != nil {
				log.Error().Err(err).Msg("Failed to start Work Order Consumer")
			} else {
				defer consumer.Stop()
				log.Info().Msg("Work Order Consumer started successfully")
			}
		}
	}

	router := routes.SetupRouter(deliveryHandler, dispenserHandler, workOrderHandler, termsSessionHandler, deliveryWithTermsHandler, mobileDeliveryHandler, cfg)

	log.Info().Str("port", cfg.Port).Msgf(constants.MsgServerRunning, cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg(constants.MsgErrorStartingServer)
	}
}
