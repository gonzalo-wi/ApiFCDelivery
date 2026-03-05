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
	deliveryStore store.DeliveryStore
	publisher     *RabbitMQPublisher
	pdfService    PDFService
	emailService  EmailService
}

func NewMobileDeliveryService(deliveryStore store.DeliveryStore, publisher *RabbitMQPublisher) MobileDeliveryService {
	return &mobileDeliveryService{
		deliveryStore: deliveryStore,
		publisher:     publisher,
		pdfService:    nil,
		emailService:  nil,
	}
}

func NewMobileDeliveryServiceWithServices(deliveryStore store.DeliveryStore, publisher *RabbitMQPublisher, pdfService PDFService, emailService EmailService) MobileDeliveryService {
	return &mobileDeliveryService{
		deliveryStore: deliveryStore,
		publisher:     publisher,
		pdfService:    pdfService,
		emailService:  emailService,
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

	// Construir respuesta con información de items de dispensers
	deliveryInfo := &dto.DeliveryInfoDTO{
		ID:          foundDelivery.ID,
		NroCta:      foundDelivery.NroCta,
		NroRto:      foundDelivery.NroRto,
		Cantidad:    foundDelivery.Cantidad,
		TipoEntrega: string(foundDelivery.TipoEntrega),
		FechaAccion: foundDelivery.FechaAccion.String(),
	}

	itemDispensers := make([]dto.ItemDispenserInfoDTO, 0, len(foundDelivery.ItemDispensers))
	for _, item := range foundDelivery.ItemDispensers {
		itemDispensers = append(itemDispensers, dto.ItemDispenserInfoDTO{
			Tipo:     string(item.Tipo),
			Cantidad: item.Cantidad,
		})
	}

	log.Info().
		Int("delivery_id", foundDelivery.ID).
		Str("token", req.Token).
		Str("nro_cta", req.NroCta).
		Msg("Token validated successfully")

	return &dto.ValidateTokenResponse{
		Valid:          true,
		Message:        "Token válido",
		Delivery:       deliveryInfo,
		ItemDispensers: itemDispensers,
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
		NroCta:     delivery.NroCta,
		Name:       delivery.Name,
		Email:      delivery.Email,
		Address:    delivery.Address,
		Locality:   delivery.Locality,
		NroRto:     delivery.NroRto,
		CreatedAt:  delivery.CreatedAt.Format("2006-01-02"),
		TipoAccion: string(delivery.TipoEntrega),
		Token:      delivery.Token,
		Dispensers: dispensersMsg,
		DeliveryID: delivery.ID,
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

	// 1. Construir WorkOrderRequest para generar el PDF
	dispensers := make([]dto.WorkOrderDispenserRequest, 0, len(validatedDispensers))
	for _, code := range validatedDispensers {
		dispensers = append(dispensers, dto.WorkOrderDispenserRequest{
			NroSerie: code,
		})
	}

	workOrderReq := &dto.WorkOrderRequest{
		NroCta:     delivery.NroCta,
		Name:       delivery.Name,
		Address:    delivery.Address,
		Locality:   delivery.Locality,
		NroRto:     delivery.NroRto,
		CreatedAt:  delivery.FechaAccion.Format("2006-01-02"),
		Dispensers: dispensers,
		TipoAccion: string(delivery.TipoEntrega),
		Token:      delivery.Token,
	}

	// 2. Generar PDF
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
	dispensersList := ""
	for i, code := range dispensers {
		dispensersList += fmt.Sprintf("<li>Dispenser #%d: %s</li>", i+1, code)
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2c3e50; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f9f9f9; padding: 20px; border: 1px solid #ddd; }
        .info-box { background-color: #e8f5e9; padding: 15px; margin: 15px 0; border-left: 4px solid #4caf50; }
        .dispensers-list { background-color: white; padding: 15px; margin: 10px 0; }
        .footer { text-align: center; padding: 20px; color: #777; font-size: 12px; }
        h2 { color: #2c3e50; margin-top: 0; }
        ul { padding-left: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✓ Instalación Completada</h1>
        </div>
        <div class="content">
            <h2>¡Instalación Exitosa!</h2>
            <p>Estimado/a <strong>%s</strong>,</p>
            <p>Le informamos que la instalación de sus dispensers se ha completado exitosamente en <strong>%s</strong>.</p>
            
            <div class="info-box">
                <strong>📋 Orden de Trabajo: %s</strong><br>
                <strong>📍 Dirección:</strong> %s<br>
                <strong>📅 Fecha:</strong> %s<br>
                <strong>📦 Cantidad de dispensers:</strong> %d
            </div>

            <div class="dispensers-list">
                <h3>Dispensers Instalados:</h3>
                <ul>
                    %s
                </ul>
            </div>

            <p>Encontrará adjunto el documento PDF con los detalles completos de la orden de trabajo y los números de serie de los dispensers instalados.</p>
            
            <p>Ante cualquier consulta, no dude en contactarnos.</p>
            
            <p><strong>Gracias por confiar en El Jumillano!</strong></p>
        </div>
        <div class="footer">
            <p>Este es un email automático, por favor no responder.<br>
            El Jumillano - Sistema de Gestión de Entregas</p>
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
		dispensersList,
	)
}
