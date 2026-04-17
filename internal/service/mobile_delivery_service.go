package service

import (
	"context"
	"fmt"
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
}

func NewMobileDeliveryService(deliveryStore store.DeliveryStore, publisher *RabbitMQPublisher) MobileDeliveryService {
	return &mobileDeliveryService{
		deliveryStore:     deliveryStore,
		termsSessionStore: nil,
		publisher:         publisher,
		pdfService:        nil,
		emailService:      nil,
	}
}

func NewMobileDeliveryServiceWithServices(deliveryStore store.DeliveryStore, termsSessionStore store.TermsSessionStore, publisher *RabbitMQPublisher, pdfService PDFService, emailService EmailService) MobileDeliveryService {
	return &mobileDeliveryService{
		deliveryStore:     deliveryStore,
		termsSessionStore: termsSessionStore,
		publisher:         publisher,
		pdfService:        pdfService,
		emailService:      emailService,
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

// CompleteDelivery marca el delivery como completado y publica mensaje a RabbitMQ
func (s *mobileDeliveryService) CompleteDelivery(ctx context.Context, req dto.MobileCompleteDeliveryRequest) (*dto.MobileCompleteDeliveryResponse, error) {
	// 1. Obtener el delivery
	delivery, err := s.deliveryStore.FindByID(ctx, req.DeliveryID)
	if err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error fetching delivery")
		return nil, fmt.Errorf("error buscando delivery: %w", err)
	}

	// 2. Validar token
	if delivery.Token != req.Token {
		log.Warn().
			Int("delivery_id", req.DeliveryID).
			Msg("Invalid token for completing delivery")
		return nil, fmt.Errorf("token inválido")
	}

	// 3. Validar estado
	if delivery.Estado != models.Pendiente {
		log.Warn().
			Int("delivery_id", req.DeliveryID).
			Str("estado", string(delivery.Estado)).
			Msg("Delivery already processed")
		return nil, fmt.Errorf("la entrega ya fue procesada (estado: %s)", delivery.Estado)
	}

	// 4. Validar que se recibieron dispensers validados
	totalEntregado := uint(len(req.ValidatedDispensers))
	if totalEntregado == 0 {
		log.Warn().
			Int("delivery_id", req.DeliveryID).
			Msg("No dispensers were delivered")
		return nil, fmt.Errorf("debe entregar al menos un dispenser")
	}

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Uint("expected", delivery.Cantidad).
		Uint("delivered", totalEntregado).
		Int("validated_dispensers", len(req.ValidatedDispensers)).
		Msg("Delivery completed with validated dispensers")

	// 5. Actualizar datos del cliente desde la app móvil
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

	// Guardar los códigos reales de dispensers validados
	delivery.ValidatedDispensers = models.StringArray(req.ValidatedDispensers)

	// Guardar el número de orden de trabajo asignado por la app móvil
	delivery.OrderNumber = req.OrderNumber

	// 6. Crear items de dispensers basándose en los códigos validados
	// Agrupamos por tipo (inferido del original o mantenemos la estructura original)
	newItemDispensers := make([]models.ItemDispenser, 0)

	// Si el delivery original tiene items con tipos, mantenemos la proporción
	// Si no, creamos un solo item con todos los dispensers
	if len(delivery.ItemDispensers) > 0 {
		// Mantener la estructura de tipos original pero actualizar cantidad real entregada
		for _, origItem := range delivery.ItemDispensers {
			newItemDispensers = append(newItemDispensers, models.ItemDispenser{
				Tipo:     origItem.Tipo,
				Cantidad: origItem.Cantidad,
			})
		}
	} else {
		// Crear un item genérico con todos los dispensers
		newItemDispensers = append(newItemDispensers, models.ItemDispenser{
			Tipo:     models.TipoDispenserPie, // Por defecto tipo P
			Cantidad: totalEntregado,
		})
	}

	// 7. Actualizar items, estado y cantidad
	delivery.ItemDispensers = newItemDispensers
	delivery.Estado = models.Completado
	delivery.Cantidad = totalEntregado // Cantidad real entregada
	delivery.UpdatedAt = time.Now()

	if err := s.deliveryStore.Update(ctx, delivery); err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error updating delivery status")
		return nil, fmt.Errorf("error actualizando delivery: %w", err)
	}

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Str("name", delivery.Name).
		Str("locality", delivery.Locality).
		Msg("Client data updated from mobile app")

	log.Info().
		Int("delivery_id", req.DeliveryID).
		Msg("Delivery marked as completed")

	// 8. Construir mensaje para RabbitMQ con los códigos reales de dispensers validados
	dispensersMsg := make([]dto.DispenserMessage, 0, len(req.ValidatedDispensers))
	for _, dispenserCode := range req.ValidatedDispensers {
		dispensersMsg = append(dispensersMsg, dto.DispenserMessage{
			NroSerie: dispenserCode,
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
		CreatedAt:   delivery.CreatedAt.Format("2006-01-02"),
		TipoAccion:  string(delivery.TipoEntrega),
		Token:       delivery.Token,
		Dispensers:  dispensersMsg,
		DeliveryID:  delivery.ID,
	}

	// 9. Publicar a RabbitMQ
	workOrderQueued := false
	err = s.publisher.PublishWorkOrder(ctx, workOrderMsg)
	if err != nil {
		log.Error().Err(err).Int("delivery_id", req.DeliveryID).Msg("Error publishing work order message")
	} else {
		workOrderQueued = true
		log.Info().
			Int("delivery_id", req.DeliveryID).
			Int("dispensers_count", len(req.ValidatedDispensers)).
			Msg("Work order message published successfully with validated dispensers")
	}

	// 10. Generar PDF y enviar email si está configurado
	if s.pdfService != nil && s.emailService != nil && delivery.Email != "" {
		go s.sendCompletionEmailWithPDF(context.Background(), delivery, req.ValidatedDispensers)
	}

	// 11. Construir respuesta con items de dispensers entregados
	itemDispensersResponse := make([]dto.ItemDispenserCompletedDTO, 0, len(delivery.ItemDispensers))
	for _, item := range delivery.ItemDispensers {
		itemDispensersResponse = append(itemDispensersResponse, dto.ItemDispenserCompletedDTO{
			Tipo:     string(item.Tipo),
			Cantidad: item.Cantidad,
		})
	}

	return &dto.MobileCompleteDeliveryResponse{
		NroCta:              delivery.NroCta,
		Name:                delivery.Name,
		Email:               delivery.Email,
		Address:             delivery.Address,
		Locality:            delivery.Locality,
		NroRto:              delivery.NroRto,
		CreatedAt:           delivery.CreatedAt.Format("2006-01-02"),
		TipoAccion:          string(delivery.TipoEntrega),
		Token:               delivery.Token,
		OrderNumber:         delivery.OrderNumber,
		ItemDispensers:      itemDispensersResponse,
		ValidatedDispensers: []string(delivery.ValidatedDispensers),
		WorkOrderQueued:     workOrderQueued,
	}, nil
}

// sendCompletionEmailWithPDF genera el PDF de la orden de trabajo y envía el email de confirmación
func (s *mobileDeliveryService) sendCompletionEmailWithPDF(ctx context.Context, delivery *models.Delivery, validatedDispensers []string) {
	localLog := log.With().
		Int("delivery_id", int(delivery.ID)).
		Str("email", delivery.Email).
		Logger()

	localLog.Info().Msg("Generating PDF and sending completion email")

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
	dispensers := make([]dto.WorkOrderDispenserRequest, 0, len(validatedDispensers))
	for _, code := range validatedDispensers {
		dispensers = append(dispensers, dto.WorkOrderDispenserRequest{
			NroSerie: code,
		})
	}

	workOrderReq := &dto.WorkOrderRequest{
		NroCta:      delivery.NroCta,
		Name:        delivery.Name,
		Address:     delivery.Address,
		Locality:    delivery.Locality,
		NroRto:      delivery.NroRto,
		CreatedAt:   delivery.FechaAccion.Format("2006-01-02"),
		AcceptedAt:  acceptedAtStr,
		Dispensers:  dispensers,
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
	emailHTML := s.buildCompletionEmailHTML(delivery, validatedDispensers, orderNumber)

	// 4. Enviar email con PDF adjunto
	err = s.emailService.SendHTMLEmailWithPDFBytes(
		ctx,
		delivery.Email,
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
