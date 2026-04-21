package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"

	"github.com/rs/zerolog/log"
)

type MobileDeliveryService interface {
	ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error)
	CompleteDelivery(ctx context.Context, req dto.MobileCompleteDeliveryRequest) (*dto.MobileCompleteDeliveryResponse, error)
	SearchDeliveries(ctx context.Context, fechaAccion string, nroRto string) ([]dto.MobileDeliverySearchResponse, error)
}

type mobileDeliveryService struct {
	deliveryStore     store.DeliveryStore
	termsSessionStore store.TermsSessionStore
	publisher         *RabbitMQPublisher
	pdfService        PDFService
	emailService      EmailService
	clientLookup      ClientLookupService
}

func NewMobileDeliveryService(deliveryStore store.DeliveryStore, publisher *RabbitMQPublisher) MobileDeliveryService {
	return &mobileDeliveryService{
		deliveryStore:     deliveryStore,
		termsSessionStore: nil,
		publisher:         publisher,
		pdfService:        nil,
		emailService:      nil,
		clientLookup:      nil,
	}
}

func NewMobileDeliveryServiceWithServices(deliveryStore store.DeliveryStore, termsSessionStore store.TermsSessionStore, publisher *RabbitMQPublisher, pdfService PDFService, emailService EmailService, clientLookup ClientLookupService) MobileDeliveryService {
	return &mobileDeliveryService{
		deliveryStore:     deliveryStore,
		termsSessionStore: termsSessionStore,
		publisher:         publisher,
		pdfService:        pdfService,
		emailService:      emailService,
		clientLookup:      clientLookup,
	}
}

// SearchDeliveries busca deliveries por fecha (obligatorio) y opcionalmente por nro_rto
func (s *mobileDeliveryService) SearchDeliveries(ctx context.Context, fechaAccion string, nroRto string) ([]dto.MobileDeliverySearchResponse, error) {
	// Parsear y validar fecha
	parsedDate, err := time.Parse("2006-01-02", fechaAccion)
	if err != nil {
		return nil, fmt.Errorf("fecha inválida. Formato esperado: YYYY-MM-DD")
	}

	// Buscar deliveries
	deliveries, err := s.deliveryStore.FindByRto(ctx, nroRto, &parsedDate)
	if err != nil {
		return nil, fmt.Errorf("error buscando deliveries: %w", err)
	}

	// Mapear a respuesta simplificada
	results := make([]dto.MobileDeliverySearchResponse, 0, len(deliveries))
	for _, d := range deliveries {
		results = append(results, dto.MobileDeliverySearchResponse{
			ID:          d.ID,
			FechaAccion: d.FechaAccion.Format("2006-01-02"),
			NroCta:      d.NroCta,
			Token:       d.Token,
		})
	}

	return results, nil
}

// ValidateToken valida el token del cliente junto con nro_cta y fecha para mayor seguridad
func (s *mobileDeliveryService) ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error) {
	// Buscar delivery directamente en BD con índices optimizados
	foundDelivery, err := s.deliveryStore.FindByTokenAndFilters(ctx, req.Token, req.NroCta, req.FechaAccion, models.Pendiente)
	if err != nil {
		log.Error().Err(err).Msg("Error validating token")
		return nil, fmt.Errorf("error validando token: %w", err)
	}

	if foundDelivery == nil {
		return &dto.ValidateTokenResponse{
			Valid:   false,
			Message: "Datos de validación incorrectos o entrega ya completada",
		}, nil
	}

	log.Info().
		Int("delivery_id", foundDelivery.ID).
		Str("token", req.Token).
		Str("nro_cta", req.NroCta).
		Msg("Token validated successfully")

	return &dto.ValidateTokenResponse{
		Valid:   true,
		Message: "Token válido",
	}, nil
}

// CompleteDelivery marca el delivery como completado y publica mensaje a RabbitMQ.
// Para instalaciones pre-coordinadas (delivery_id > 0) valida token y estado.
// Para retiros y recambios (delivery_id == 0) crea un delivery nuevo en el momento.
func (s *mobileDeliveryService) CompleteDelivery(ctx context.Context, req dto.MobileCompleteDeliveryRequest) (*dto.MobileCompleteDeliveryResponse, error) {
	tipoEntrega := deriveTipoEntrega(req.Operations)

	var delivery *models.Delivery
	var err error

	if req.DeliveryID > 0 {
		// --- Flujo instalación pre-coordinada ---
		delivery, err = s.deliveryStore.FindByID(ctx, req.DeliveryID)
		if err != nil {
			log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error fetching delivery")
			return nil, fmt.Errorf("error buscando delivery: %w", err)
		}

		if req.Token != "" && delivery.Token != req.Token {
			log.Warn().Int("delivery_id", req.DeliveryID).Msg("Invalid token for completing delivery")
			return nil, fmt.Errorf("token inválido")
		}

		if delivery.Estado != models.Pendiente {
			log.Warn().
				Int("delivery_id", req.DeliveryID).
				Str("estado", string(delivery.Estado)).
				Msg("Delivery already processed")
			return nil, fmt.Errorf("la entrega ya fue procesada (estado: %s)", delivery.Estado)
		}
	} else {
		// --- Flujo retiro/recambio no coordinado: crear delivery en el momento ---
		nroRto := req.NroRto
		if nroRto == "" {
			nroRto = req.OrderNumber
		}
		delivery = &models.Delivery{
			NroCta:       req.NroCta,
			Name:         req.Name,
			Email:        req.Email,
			Address:      req.Address,
			Locality:     req.Locality,
			NroRto:       nroRto,
			Estado:       models.Pendiente,
			TipoEntrega:  tipoEntrega,
			EntregadoPor: models.Tecnico,
			Cantidad:     uint(len(req.Operations)),
			FechaAccion:  models.CustomDate{Time: time.Now()},
		}
		if err = s.deliveryStore.Create(ctx, delivery); err != nil {
			log.Error().Err(err).Msg("Error creating on-the-fly delivery")
			return nil, fmt.Errorf("error creando delivery: %w", err)
		}
		log.Info().Int("delivery_id", delivery.ID).Str("tipo", string(tipoEntrega)).Msg("On-the-fly delivery created")
	}

	// Actualizar datos del cliente si vienen en el request
	if req.Name != "" {
		delivery.Name = req.Name
	}
	if req.Email != "" {
		delivery.Email = req.Email
	}
	if req.Address != "" {
		delivery.Address = req.Address
	}
	if req.Locality != "" {
		delivery.Locality = req.Locality
	}

	// Guardar solo los dispensers instalados en ValidatedDispensers
	installed := make([]string, 0, len(req.Operations))
	for _, op := range req.Operations {
		if op.InstalledDispenserCode != "" {
			installed = append(installed, op.InstalledDispenserCode)
		}
	}
	delivery.ValidatedDispensers = models.StringArray(installed)
	delivery.TipoEntrega = tipoEntrega
	delivery.OrderNumber = req.OrderNumber
	delivery.Estado = models.Completado
	delivery.Cantidad = uint(len(req.Operations))
	delivery.UpdatedAt = time.Now()

	if err = s.deliveryStore.Update(ctx, delivery); err != nil {
		log.Error().Err(err).Int("delivery_id", delivery.ID).Msg("Error updating delivery")
		return nil, fmt.Errorf("error actualizando delivery: %w", err)
	}

	log.Info().
		Int("delivery_id", delivery.ID).
		Str("tipo", string(tipoEntrega)).
		Int("operations", len(req.Operations)).
		Msg("Delivery completed")

	// Construir mensaje para RabbitMQ
	opsMsg := make([]dto.OperationMessage, 0, len(req.Operations))
	for _, op := range req.Operations {
		opsMsg = append(opsMsg, dto.OperationMessage{
			Type:                   op.Type,
			InstalledDispenserCode: op.InstalledDispenserCode,
			RetiredDispenserCode:   op.RetiredDispenserCode,
		})
	}

	workOrderMsg := dto.WorkOrderMessageDTO{
		OrderNumber: req.OrderNumber,
		NroCta:      delivery.NroCta,
		Name:        delivery.Name,
		Email:       delivery.Email,
		Address:     delivery.Address,
		Locality:    delivery.Locality,
		NroRto:      delivery.NroRto,
		CreatedAt:   time.Now().Format("2006-01-02"),
		TipoAccion:  string(tipoEntrega),
		Token:       delivery.Token,
		Operations:  opsMsg,
		DeliveryID:  delivery.ID,
	}

	workOrderQueued := false
	if err = s.publisher.PublishWorkOrder(ctx, workOrderMsg); err != nil {
		log.Error().Err(err).Int("delivery_id", delivery.ID).Msg("Error publishing work order")
	} else {
		workOrderQueued = true
		log.Info().Int("delivery_id", delivery.ID).Int("operations", len(opsMsg)).Msg("Work order queued")
	}

	if s.pdfService != nil && s.emailService != nil {
		go s.sendCompletionEmailWithPDF(context.Background(), delivery, req.Operations)
	}

	// Mapear operaciones a response
	opsCompleted := make([]dto.OperationCompletedDTO, 0, len(req.Operations))
	for _, op := range req.Operations {
		opsCompleted = append(opsCompleted, dto.OperationCompletedDTO{
			Type:                   op.Type,
			InstalledDispenserCode: op.InstalledDispenserCode,
			RetiredDispenserCode:   op.RetiredDispenserCode,
		})
	}

	return &dto.MobileCompleteDeliveryResponse{
		DeliveryID:      delivery.ID,
		NroCta:          delivery.NroCta,
		Name:            delivery.Name,
		Email:           delivery.Email,
		Address:         delivery.Address,
		Locality:        delivery.Locality,
		NroRto:          delivery.NroRto,
		TipoAccion:      string(tipoEntrega),
		OrderNumber:     req.OrderNumber,
		Operations:      opsCompleted,
		WorkOrderQueued: workOrderQueued,
	}, nil
}

// deriveTipoEntrega infiere el TipoEntrega a partir del conjunto de operaciones
func deriveTipoEntrega(ops []dto.DispenserOperation) models.TipoEntrega {
	types := make(map[string]bool)
	for _, op := range ops {
		types[op.Type] = true
	}
	if len(types) > 1 {
		return models.Mixto
	}
	if types["installation"] {
		return models.Instalacion
	}
	if types["retirement"] {
		return models.Retiro
	}
	return models.Recambio
}

// sendCompletionEmailWithPDF genera el PDF de la orden de trabajo y envía el email de confirmación
func (s *mobileDeliveryService) sendCompletionEmailWithPDF(ctx context.Context, delivery *models.Delivery, operations []dto.DispenserOperation) {
	localLog := log.With().
		Int("delivery_id", int(delivery.ID)).
		Str("nro_cta", delivery.NroCta).
		Logger()

	localLog.Info().Msg("Generating PDF and sending completion email")

	// Resolver email del destinatario
	// Para retiros/recambios el delivery puede no tener email; se consulta la API externa
	emailTo := strings.TrimSpace(delivery.Email)
	if emailTo == "" && s.clientLookup != nil {
		emailTo = s.clientLookup.GetClientEmail(ctx, delivery.NroCta)
	} else if emailTo == "" {
		emailTo = fallbackClientEmail
	}

	localLog.Info().Str("email_to", emailTo).Msg("Email recipient resolved")

	// 1. Obtener fecha de aceptación de términos si existe
	acceptedAtStr := ""
	if delivery.TermsSessionID != nil && s.termsSessionStore != nil {
		termsSession, err := s.termsSessionStore.GetByID(ctx, *delivery.TermsSessionID)
		if err == nil && termsSession != nil && termsSession.AcceptedAt != nil {
			acceptedAtStr = termsSession.AcceptedAt.Format("02/01/2006 15:04")
			localLog.Info().
				Str("accepted_at", acceptedAtStr).
				Msg("Found terms acceptance date")
		} else {
			localLog.Warn().Msg("Terms session not found or not accepted, will use current time")
		}
	}

	// Si no hay fecha de aceptación, usar la fecha actual
	if acceptedAtStr == "" {
		acceptedAtStr = time.Now().Format("02/01/2006 15:04")
	}

	// 2. Construir WorkOrderRequest para generar el PDF
	workOrderReq := &dto.WorkOrderRequest{
		NroCta:      delivery.NroCta,
		Name:        delivery.Name,
		Address:     delivery.Address,
		Locality:    delivery.Locality,
		NroRto:      delivery.NroRto,
		CreatedAt:   delivery.FechaAccion.Format("2006-01-02"),
		AcceptedAt:  acceptedAtStr,
		Operations:  operations,
		TipoAccion:  string(delivery.TipoEntrega),
		Token:       delivery.Token,
		OrderNumber: delivery.OrderNumber,
	}

	// 3. Generar PDF
	pdfBytes, orderNumber, err := s.pdfService.GenerateWorkOrderPDF(ctx, workOrderReq)
	if err != nil {
		localLog.Error().Err(err).Msg("Error generating PDF for completion email")
		return
	}

	localLog.Info().Str("order_number", orderNumber).Msg("PDF generated successfully")

	// 3. Crear HTML del email de instalación completada
	installedCodes := make([]string, 0)
	for _, op := range operations {
		if op.InstalledDispenserCode != "" {
			installedCodes = append(installedCodes, op.InstalledDispenserCode)
		}
	}
	emailHTML := s.buildCompletionEmailHTML(delivery, installedCodes, orderNumber)

	// 4. Enviar email con PDF adjunto
	err = s.emailService.SendHTMLEmailWithPDFBytes(
		ctx,
		emailTo,
		"Instalación Completada - El Jumillano",
		emailHTML,
		pdfBytes,
		fmt.Sprintf("orden_trabajo_%s.pdf", orderNumber),
	)

	if err != nil {
		localLog.Error().Err(err).Msg("Error sending completion email")
		return
	}

	localLog.Info().Msg("Completion email sent successfully with PDF attachment")
}

// buildCompletionEmailHTML construye el HTML del email de instalación completada
func (s *mobileDeliveryService) buildCompletionEmailHTML(delivery *models.Delivery, dispensers []string, orderNumber string) string {
	dispensersRows := ""
	for i, code := range dispensers {
		bgColor := "#ffffff"
		if i%2 == 1 {
			bgColor = "#f8f9fa"
		}
		dispensersRows += fmt.Sprintf(`
			<tr style="background-color: %s;">
				<td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: center; font-weight: 600; color: #2c3e50;">%d</td>
				<td style="padding: 12px; border-bottom: 1px solid #e0e0e0; font-family: 'Courier New', monospace; color: #1976d2; font-weight: 500;">%s</td>
			</tr>`, bgColor, i+1, code)
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Instalación Completada</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            line-height: 1.6; 
            color: #333;
            background-color: #ffffff;
        }
        .email-wrapper { 
            background-color: #ffffff;
            padding: 40px 20px;
        }
        .container { 
            max-width: 600px; 
            margin: 0 auto; 
            background-color: #ffffff;
            border-radius: 16px;
            overflow: hidden;
            box-shadow: 0 10px 40px rgba(0,0,0,0.15);
        }
        .header { 
            background: linear-gradient(135deg, #1e88e5 0%%, #1565c0 100%%);
            color: white; 
            padding: 40px 30px;
            text-align: center;
            position: relative;
        }
        .header::after {
            content: '';
            position: absolute;
            bottom: -20px;
            left: 0;
            right: 0;
            height: 40px;
            background: #ffffff;
            border-radius: 50%% 50%% 0 0 / 100%% 100%% 0 0;
        }
        .success-icon {
            width: 80px;
            height: 80px;
            background-color: #4caf50;
            border-radius: 50%%;
            display: table-cell;
            vertical-align: middle;
            text-align: center;
            margin: 0 auto 15px;
            box-shadow: 0 4px 15px rgba(76, 175, 80, 0.4);
            font-size: 52px;
            line-height: 80px;
            color: white;
            font-weight: bold;
        }
        .header h1 { 
            font-size: 28px; 
            margin: 0;
            font-weight: 600;
        }
        .header p {
            margin-top: 10px;
            font-size: 16px;
            opacity: 0.95;
        }
        .content { 
            padding: 50px 30px 30px;
        }
        .greeting {
            font-size: 18px;
            color: #2c3e50;
            margin-bottom: 20px;
        }
        .greeting strong {
            color: #1976d2;
            font-weight: 600;
        }
        .message {
            font-size: 16px;
            color: #555;
            margin-bottom: 30px;
            line-height: 1.8;
        }
        .info-card { 
            background: linear-gradient(135deg, #e8f5e9 0%%, #c8e6c9 100%%);
            padding: 25px;
            margin: 25px 0;
            border-radius: 12px;
            border-left: 5px solid #4caf50;
            box-shadow: 0 3px 10px rgba(0,0,0,0.08);
        }
        .info-card h3 {
            color: #2e7d32;
            font-size: 18px;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
        }
        .info-card h3::before {
            content: '📋';
            margin-right: 10px;
            font-size: 24px;
        }
        .info-row {
            display: flex;
            padding: 8px 0;
            border-bottom: 1px solid rgba(46, 125, 50, 0.1);
        }
        .info-row:last-child {
            border-bottom: none;
        }
        .info-label {
            font-weight: 600;
            color: #2e7d32;
            min-width: 140px;
            font-size: 14px;
        }
        .info-value {
            color: #1b5e20;
            font-size: 14px;
        }
        .dispensers-section {
            margin: 30px 0;
        }
        .dispensers-section h3 {
            color: #2c3e50;
            font-size: 20px;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
        }
        .dispensers-section h3::before {
            content: '🔧';
            margin-right: 10px;
            font-size: 26px;
        }
        .dispensers-table {
            width: 100%%;
            border-collapse: collapse;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.08);
        }
        .dispensers-table thead {
            background: linear-gradient(135deg, #1e88e5 0%%, #1565c0 100%%);
            color: white;
        }
        .dispensers-table th {
            padding: 15px;
            text-align: left;
            font-weight: 600;
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .dispensers-table th:first-child {
            text-align: center;
            width: 80px;
        }
        .pdf-notice {
            background: linear-gradient(135deg, #fff3e0 0%%, #ffe0b2 100%%);
            padding: 20px;
            border-radius: 12px;
            border-left: 5px solid #ff9800;
            margin: 25px 0;
            display: flex;
            align-items: center;
        }
        .pdf-notice::before {
            content: '📎';
            font-size: 32px;
            margin-right: 15px;
        }
        .pdf-notice p {
            margin: 0;
            color: #e65100;
            font-size: 15px;
            line-height: 1.6;
        }
        .cta-section {
            background: linear-gradient(135deg, #e3f2fd 0%%, #bbdefb 100%%);
            padding: 25px;
            border-radius: 12px;
            text-align: center;
            margin: 25px 0;
        }
        .cta-section p {
            color: #1565c0;
            font-size: 15px;
            margin-bottom: 0;
        }
        .cta-section strong {
            color: #0d47a1;
            font-size: 16px;
        }
        .footer { 
            background: linear-gradient(135deg, #37474f 0%%, #263238 100%%);
            color: #b0bec5;
            text-align: center; 
            padding: 30px;
        }
        .footer-brand {
            font-size: 20px;
            font-weight: 700;
            color: #ffffff;
            margin-bottom: 10px;
        }
        .footer p { 
            font-size: 13px;
            margin: 5px 0;
            line-height: 1.6;
        }
        .footer-note {
            margin-top: 15px;
            padding-top: 15px;
            border-top: 1px solid rgba(255,255,255,0.1);
            font-size: 12px;
            opacity: 0.8;
        }
        @media only screen and (max-width: 600px) {
            .email-wrapper { padding: 20px 10px; }
            .content { padding: 30px 20px 20px; }
            .header { padding: 30px 20px; }
            .info-row { flex-direction: column; }
            .info-label { min-width: auto; margin-bottom: 5px; }
        }
    </style>
</head>
<body>
    <div class="email-wrapper">
        <div class="container">
            <div class="header">
                <div class="success-icon">✓</div>
                <h1>Instalación Completada</h1>
                <p>Su servicio ya está en funcionamiento</p>
            </div>
            
            <div class="content">
                <div class="greeting">
                    Estimado/a <strong>%s</strong>,
                </div>
                
                <div class="message">
                    Nos complace informarle que la instalación de sus dispensers de agua 
                    ha sido completada exitosamente en <strong>%s</strong>. 
                    Nuestro equipo técnico ha verificado el correcto funcionamiento de todos los equipos.
                </div>
                
                <div class="info-card">
                    <h3>Detalles de su Instalación</h3>
                    <div class="info-row">
                        <span class="info-label">Orden de Trabajo:</span>
                        <span class="info-value"><strong>%s</strong></span>
                    </div>
                    <div class="info-row">
                        <span class="info-label">📍 Dirección:</span>
                        <span class="info-value">%s</span>
                    </div>
                    <div class="info-row">
                        <span class="info-label">📅 Fecha:</span>
                        <span class="info-value">%s</span>
                    </div>
                    <div class="info-row">
                        <span class="info-label">📦 Dispensers:</span>
                        <span class="info-value"><strong>%d unidades</strong></span>
                    </div>
                </div>

                <div class="dispensers-section">
                    <h3>Equipos Instalados</h3>
                    <table class="dispensers-table">
                        <thead>
                            <tr>
                                <th>Item</th>
                                <th>Número de Serie</th>
                            </tr>
                        </thead>
                        <tbody>
                            %s
                        </tbody>
                    </table>
                </div>

                <div class="pdf-notice">
                    <p>
                        <strong>Documento Adjunto:</strong> Encontrará el comprobante de instalación en formato PDF 
                        con todos los detalles técnicos y números de serie de los equipos instalados.
                    </p>
                </div>

                <div class="cta-section">
                    <p><strong>¡Gracias por elegirnos!</strong></p>
                    <p>Ante cualquier consulta o inconveniente, no dude en contactarnos. 
                    Estamos para servirle.</p>
                </div>
            </div>
            
            <div class="footer">
                <div class="footer-brand">El Jumillano</div>
                <p>Sistema de Gestión de Entregas</p>
                <p>Calidad y confianza en cada servicio</p>
                <div class="footer-note">
                    Este es un correo electrónico automático, por favor no responder.<br>
                    Si necesita asistencia, contacte a nuestro servicio al cliente.
                </div>
            </div>
        </div>
    </div>
</body>
</html>
	`,
		delivery.Name,
		delivery.Locality,
		orderNumber,
		delivery.Address,
		delivery.FechaAccion.Format("02/01/2006"),
		len(dispensers),
		dispensersRows,
	)
}
