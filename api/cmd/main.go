package main

import (
	"GoFrioCalor/config"
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/routes"
	"GoFrioCalor/internal/service"
	"GoFrioCalor/internal/store"
	"GoFrioCalor/internal/transport"
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

	// Create SQLX connection for audit store
	sqlxDB, err := config.NewSQLXDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database with sqlx")
	}

	// Stores
	deliveryStore := store.NewDeliveryStore(db)
	workOrderStore := store.NewWorkOrderStore(db)
	termsSessionStore := store.NewTermsSessionStore(db)
	auditEventStore := store.NewAuditEventStore(sqlxDB)

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

	// Audit Service
	auditService := service.NewAuditService(auditEventStore)
	auditHandler := transport.NewAuditHandler(auditService)

	// Términos y Condiciones con Infobip
	infobipClient := service.NewInfobipClient(cfg.InfobipBaseURL, cfg.InfobipAPIKey)
	termsSessionService := service.NewTermsSessionService(termsSessionStore, infobipClient)
	termsSessionHandler := transport.NewTermsSessionHandler(termsSessionService, cfg.AppBaseURL, cfg.TermsTTLHours)

	// Inicializar email service real o mock según configuración
	var emailService service.EmailService
	emailService, err = service.NewSMTPEmailService(
		cfg.EmailHost,
		cfg.EmailPort,
		cfg.EmailFrom,
		cfg.EmailPassword,
		cfg.EmailTo,
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to initialize SMTP email service, using mock instead")
		emailService = service.NewMockEmailService()
	} else {
		log.Info().
			Str("host", cfg.EmailHost).
			Str("port", cfg.EmailPort).
			Str("from", cfg.EmailFrom).
			Str("to", cfg.EmailTo).
			Msg("SMTP Email Service initialized successfully")
	}

	// Services
	deliveryService := service.NewDeliveryServiceWithEmail(deliveryStore, emailService)
	deliveryHandler := transport.NewDeliveryHandler(deliveryService, auditService)

	pdfService := service.NewPDFService(workOrderStore)
	workOrderHandler := transport.NewWorkOrderHandler(pdfService)

	// Flujo integrado: Entregas con Términos y Condiciones
	deliveryWithTermsService := service.NewDeliveryWithTermsService(deliveryStore, termsSessionStore, termsSessionService)
	deliveryWithTermsHandler := transport.NewDeliveryWithTermsHandler(deliveryWithTermsService, cfg.AppBaseURL, cfg.TermsTTLHours)

	// Mobile Delivery - Validación y Completar Entregas
	var mobileDeliveryHandler *transport.MobileDeliveryHandler
	if rabbitPublisher != nil {
		clientLookupService := service.NewClientLookupService(cfg.ClientLookupBaseURL, cfg.ClientLookupAPIKey, cfg.ClientLookupDefaultEmail)
		mobileDeliveryService := service.NewMobileDeliveryServiceWithServices(deliveryStore, termsSessionStore, rabbitPublisher, pdfService, emailService, clientLookupService)
		mobileDeliveryHandler = transport.NewMobileDeliveryHandler(mobileDeliveryService, auditService)
		log.Info().Msg("Mobile Delivery Service initialized with PDF and Email services")
	}

	// RabbitMQ Consumer para Work Orders
	if rabbitPublisher != nil {
		realPDFGenerator := service.NewRealWorkOrderPDFGenerator()

		consumer, err := service.NewWorkOrderConsumer(
			rabbitConfig,
			workOrderStore,
			deliveryStore,
			realPDFGenerator,
			emailService,
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

	router := routes.SetupRouter(deliveryHandler, workOrderHandler, termsSessionHandler, deliveryWithTermsHandler, mobileDeliveryHandler, auditHandler, cfg)

	// Scheduler: cancelar deliveries pendientes cuya fecha_accion ya pasó (se ejecuta a medianoche)
	scheduler := service.NewScheduler(deliveryStore)
	scheduler.Start()
	defer scheduler.Stop()

	log.Info().Str("port", cfg.Port).Msgf(constants.MsgServerRunning, cfg.Port)

	if err := router.Run("0.0.0.0:" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg(constants.MsgErrorStartingServer)
	}
}
