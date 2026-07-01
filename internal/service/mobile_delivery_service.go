package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/metrics"
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

func (s *mobileDeliveryService) SearchDeliveries(ctx context.Context, fechaAccion string, nroRto string) ([]dto.MobileDeliverySearchResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", fechaAccion)
	if err != nil {
		return nil, fmt.Errorf("fecha inválida. Formato esperado: YYYY-MM-DD")
	}
	deliveries, err := s.deliveryStore.FindByRto(ctx, nroRto, &parsedDate)
	if err != nil {
		return nil, fmt.Errorf("error buscando deliveries: %w", err)
	}
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

func (s *mobileDeliveryService) ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error) {
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

func (s *mobileDeliveryService) CompleteDelivery(ctx context.Context, req dto.MobileCompleteDeliveryRequest) (*dto.MobileCompleteDeliveryResponse, error) {
	tipoEntrega := deriveTipoEntrega(req.Operations)
	var delivery *models.Delivery
	var err error
	if req.DeliveryID > 0 {
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
	installed := make([]string, 0, len(req.Operations))
	for _, op := range req.Operations {
		if op.InstalledDispenserCode != "" {
			installed = append(installed, op.InstalledDispenserCode)
		}
		if op.ServiceDispenserCode != "" {
			installed = append(installed, op.ServiceDispenserCode)
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
	opsMsg := make([]dto.OperationMessage, 0, len(req.Operations))
	for _, op := range req.Operations {
		opsMsg = append(opsMsg, dto.OperationMessage{
			Type:                   op.Type,
			InstalledDispenserCode: op.InstalledDispenserCode,
			RetiredDispenserCode:   op.RetiredDispenserCode,
			ServiceDispenserCode:   op.ServiceDispenserCode,
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
	opsCompleted := make([]dto.OperationCompletedDTO, 0, len(req.Operations))
	for _, op := range req.Operations {
		opsCompleted = append(opsCompleted, dto.OperationCompletedDTO{
			Type:                   op.Type,
			InstalledDispenserCode: op.InstalledDispenserCode,
			RetiredDispenserCode:   op.RetiredDispenserCode,
			ServiceDispenserCode:   op.ServiceDispenserCode,
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
	if types["service"] {
		return models.Service
	}
	return models.Recambio
}

func (s *mobileDeliveryService) sendCompletionEmailWithPDF(ctx context.Context, delivery *models.Delivery, operations []dto.DispenserOperation) {
	localLog := log.With().
		Int("delivery_id", int(delivery.ID)).
		Str("nro_cta", delivery.NroCta).
		Logger()
	localLog.Info().Msg("Generating PDF and sending completion email")
	emailTo := strings.TrimSpace(delivery.Email)
	if emailTo == "" && s.clientLookup != nil {
		emailTo = s.clientLookup.GetClientEmail(ctx, delivery.NroCta)
	} else if emailTo == "" {
		emailTo = fallbackClientEmail
	}
	localLog.Info().Str("email_to", emailTo).Msg("Email recipient resolved")
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

	if acceptedAtStr == "" {
		acceptedAtStr = time.Now().Format("02/01/2006 15:04")
	}
	workOrderReq := &dto.WorkOrderRequest{
		DeliveryID:  delivery.ID,
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
	pdfBytes, orderNumber, err := s.pdfService.GenerateWorkOrderPDF(ctx, workOrderReq)
	if err != nil {
		localLog.Error().Err(err).Msg("Error generating PDF for completion email")
		return
	}

	localLog.Info().Str("order_number", orderNumber).Msg("PDF generated successfully")
	installedCodes := make([]string, 0)
	for _, op := range operations {
		if op.InstalledDispenserCode != "" {
			installedCodes = append(installedCodes, op.InstalledDispenserCode)
		}
	}
	emailHTML := s.buildCompletionEmailHTML(delivery, installedCodes, orderNumber)
	emailSubject := emailTextsForTipoEntrega(delivery.TipoEntrega).subject
	err = s.emailService.SendHTMLEmailWithPDFBytesAndLogo(
		ctx,
		emailTo,
		emailSubject,
		emailHTML,
		pdfBytes,
		fmt.Sprintf("orden_trabajo_%s.pdf", orderNumber),
		"assets/images/fondo-color.png",
	)

	if err != nil {
		metrics.EmailSent("completion", false)
		localLog.Error().Err(err).Msg("Error sending completion email")
		return
	}

	metrics.EmailSent("completion", true)
	localLog.Info().Msg("Completion email sent successfully with PDF attachment")
}

type completionEmailTexts struct {
	subject        string
	title          string
	subtitle       string
	message        string
	cardTitle      string
	dispenserTitle string
	pdfNotice      string
}

func emailTextsForTipoEntrega(tipo models.TipoEntrega) completionEmailTexts {
	switch tipo {
	case models.Retiro:
		return completionEmailTexts{
			subject:        "Retiro Completado - El Jumillano",
			title:          "Retiro Completado",
			subtitle:       "El retiro fue realizado exitosamente",
			message:        "Te informamos que el retiro del/los dispenser/s de agua fue realizado exitosamente.",
			cardTitle:      "Detalles del Retiro",
			dispenserTitle: "Equipos Retirados",
			pdfNotice:      "Encontrará el comprobante del retiro en formato PDF con todos los detalles técnicos y números de serie de los equipos retirados.",
		}
	case models.Recambio:
		return completionEmailTexts{
			subject:        "Recambio Completado - El Jumillano",
			title:          "Recambio Completado",
			subtitle:       "El recambio fue realizado exitosamente",
			message:        "Te informamos que el recambio del/los dispenser/s de agua fue realizado exitosamente.",
			cardTitle:      "Detalles del Recambio",
			dispenserTitle: "Equipos Instalados",
			pdfNotice:      "Encontrará el comprobante del recambio en formato PDF con todos los detalles técnicos y números de serie de los equipos.",
		}
	case models.Service:
		return completionEmailTexts{
			subject:        "Servicio Técnico Completado - El Jumillano",
			title:          "Servicio Técnico Completado",
			subtitle:       "El servicio técnico fue completado exitosamente",
			message:        "Te informamos que el servicio técnico del/los dispenser/s de agua fue realizado exitosamente.",
			cardTitle:      "Detalles del Servicio Técnico",
			dispenserTitle: "Equipos Atendidos",
			pdfNotice:      "Encontrará el comprobante del servicio técnico en formato PDF con todos los detalles del trabajo realizado.",
		}
	case models.Mixto:
		return completionEmailTexts{
			subject:        "Servicio Completado - El Jumillano",
			title:          "Servicio Completado",
			subtitle:       "El servicio fue completado exitosamente",
			message:        "Te informamos que el servicio en el/los dispenser/s de agua fue realizado exitosamente.",
			cardTitle:      "Detalles del Servicio",
			dispenserTitle: "Equipos Involucrados",
			pdfNotice:      "Encontrará el comprobante del servicio en formato PDF con todos los detalles técnicos correspondientes.",
		}
	default:
		return completionEmailTexts{
			subject:        "Instalación Completada - El Jumillano",
			title:          "Instalación completada",
			subtitle:       "Tu servicio ya está en funcionamiento",
			message:        "Te informamos que la instalación de tu dispenser frío/calor se realizó exitosamente.",
			cardTitle:      "Detalles de su Instalación",
			dispenserTitle: "Equipos Instalados",
			pdfNotice:      "En este correo encontrás el comprobante de instalación en formato PDF con todos los detalles técnicos y los números de serie de los equipos instalados.",
		}
	}
}

func (s *mobileDeliveryService) buildCompletionEmailHTML(delivery *models.Delivery, dispensers []string, orderNumber string) string {
	texts := emailTextsForTipoEntrega(delivery.TipoEntrega)
	message := texts.message

	logoSrc := "cid:blanco.png"

	dispensersRows := ""
	for i, code := range dispensers {
		bgColor := "#ffffff"
		if i%2 == 1 {
			bgColor = "#F5F5F5"
		}
		dispensersRows += fmt.Sprintf(`
				<tr>
					<td bgcolor="%s" style="padding:10px 15px;border-bottom:1px solid #e0e0e0;text-align:center;font-size:13px;color:#1B5EA6;font-weight:bold;font-family:'Montserrat',Arial,sans-serif;">%d</td>
					<td bgcolor="%s" style="padding:10px 15px;border-bottom:1px solid #e0e0e0;font-size:13px;color:#0099CC;font-family:'Montserrat',Arial,sans-serif;">%s</td>
				</tr>`, bgColor, i+1, bgColor, code)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s</title>
</head>
<body style="margin:0;padding:0;background-color:#F5F5F5;font-family:'Montserrat',Arial,Helvetica,sans-serif;">
%s
<div style="display:none;max-height:0;overflow:hidden;mso-hide:all;font-size:1px;color:#F5F5F5;">%s&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;</div>
<table width="100%%" cellpadding="0" cellspacing="0" border="0" bgcolor="#F5F5F5">
<tr><td align="center" style="padding:20px 10px;">

<table width="600" cellpadding="0" cellspacing="0" border="0" style="max-width:600px;background-color:#ffffff;">

<!-- HEADER -->
<tr>
  <td bgcolor="#1B5EA6" style="padding:18px 24px;">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0">
      <tr>
        <td width="65%%" valign="middle">
          <table cellpadding="0" cellspacing="0" border="0">
            <tr>
              <td width="46" height="46" bgcolor="#4CAF50" style="border-radius:23px;text-align:center;vertical-align:middle;width:46px;height:46px;">
                <span style="color:#ffffff;font-size:24px;font-weight:bold;line-height:46px;display:block;font-family:'Montserrat',Arial,sans-serif;">&#10003;</span>
              </td>
              <td style="padding-left:12px;">
                <p style="margin:0;color:#ffffff;font-size:19px;font-weight:bold;line-height:1.2;font-family:'Montserrat',Arial,sans-serif;">%s</p>
                <p style="margin:4px 0 0 0;color:rgba(255,255,255,0.88);font-size:12px;font-family:'Montserrat',Arial,sans-serif;">%s</p>
              </td>
            </tr>
          </table>
        </td>
        <td width="35%%" align="right" valign="middle">
          <img src="%s" alt="IVESS" style="max-height:52px;max-width:135px;display:block;margin-left:auto;" />
        </td>
      </tr>
    </table>
  </td>
</tr>

<!-- BODY -->
<tr>
  <td style="padding:24px 28px 12px 28px;font-family:'Montserrat',Arial,sans-serif;">

    <!-- Greeting -->
    <p style="margin:0 0 4px 0;font-size:15px;color:#333333;font-family:'Montserrat',Arial,sans-serif;">Estimado/a,</p>
    <p style="margin:0 0 20px 0;font-size:15px;color:#555555;line-height:1.6;font-family:'Montserrat',Arial,sans-serif;">%s</p>

    <!-- Description card -->
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:16px;">
      <tr>
        <td bgcolor="#8DC63F" style="padding:10px 16px;border-radius:4px 4px 0 0;">
          <p style="margin:0;color:#ffffff;font-size:13px;font-weight:bold;text-transform:uppercase;letter-spacing:0.5px;font-family:'Montserrat',Arial,sans-serif;">%s</p>
        </td>
      </tr>
      <tr>
        <td bgcolor="#F5F5F5" style="padding:14px 16px;border:1px solid #e0e0e0;border-top:none;border-radius:0 0 4px 4px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0">
            <tr>
              <td style="padding:4px 0;font-size:13px;color:#666666;width:130px;font-family:'Montserrat',Arial,sans-serif;" nowrap="nowrap">Orden de Trabajo:</td>
              <td style="padding:4px 0;font-size:13px;color:#1B5EA6;font-weight:bold;font-family:'Montserrat',Arial,sans-serif;">%s</td>
            </tr>
            <tr>
              <td style="padding:4px 0;font-size:13px;color:#666666;font-family:'Montserrat',Arial,sans-serif;" nowrap="nowrap">Direcci&oacute;n:</td>
              <td style="padding:4px 0;font-size:13px;color:#333333;font-family:'Montserrat',Arial,sans-serif;">%s</td>
            </tr>
            <tr>
              <td style="padding:4px 0;font-size:13px;color:#666666;font-family:'Montserrat',Arial,sans-serif;" nowrap="nowrap">Fecha:</td>
              <td style="padding:4px 0;font-size:13px;color:#333333;font-family:'Montserrat',Arial,sans-serif;">%s</td>
            </tr>
            <tr>
              <td style="padding:4px 0;font-size:13px;color:#666666;font-family:'Montserrat',Arial,sans-serif;" nowrap="nowrap">Dispensers:</td>
              <td style="padding:4px 0;font-size:13px;color:#333333;font-weight:bold;font-family:'Montserrat',Arial,sans-serif;">%d unidades</td>
            </tr>
          </table>
        </td>
      </tr>
    </table>

    <p style="margin:0 0 16px 0;font-size:14px;color:#555555;font-family:'Montserrat',Arial,sans-serif;">Nuestro equipo t&eacute;cnico ya verific&oacute; el correcto funcionamiento de todos los equipos.</p>

    <!-- Dispensers table -->
    <p style="margin:0 0 8px 0;font-size:14px;font-weight:bold;color:#1B5EA6;font-family:'Montserrat',Arial,sans-serif;">%s</p>
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="border-collapse:collapse;margin-bottom:20px;">
      <tr>
        <th bgcolor="#1B5EA6" style="padding:10px 14px;color:#ffffff;font-size:11px;text-transform:uppercase;letter-spacing:0.5px;text-align:center;width:70px;font-family:'Montserrat',Arial,sans-serif;">&#205;TEM</th>
        <th bgcolor="#1B5EA6" style="padding:10px 14px;color:#ffffff;font-size:11px;text-transform:uppercase;letter-spacing:0.5px;text-align:left;font-family:'Montserrat',Arial,sans-serif;">N&Uacute;MERO DE SERIE</th>
      </tr>
      %s
    </table>

    <!-- PDF notice -->
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:20px;">
      <tr>
        <td bgcolor="#E3F2FD" style="padding:14px 16px;border-left:4px solid #0099CC;border-radius:0 4px 4px 0;">
          <p style="margin:0;font-size:14px;color:#1B5EA6;line-height:1.6;font-family:'Montserrat',Arial,sans-serif;">
            <span style="background-color:#FFA500;color:#ffffff;font-weight:bold;padding:2px 7px;border-radius:3px;font-size:13px;font-family:'Montserrat',Arial,sans-serif;">Documento adjunto:</span>&nbsp;%s
          </p>
        </td>
      </tr>
    </table>

    <!-- Thank you + CTA -->
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:20px;">
      <tr>
        <td bgcolor="#E3F2FD" style="padding:20px 24px;border:1px solid #0099CC;border-radius:6px;text-align:center;">
          <p style="margin:0 0 8px 0;font-size:16px;font-weight:bold;color:#1B5EA6;font-family:'Montserrat',Arial,sans-serif;">&#161;Gracias por ser parte de la familia IVESS!</p>
          <p style="margin:0 0 10px 0;font-size:14px;color:#555555;font-family:'Montserrat',Arial,sans-serif;">Agend&aacute; nuestro Whatsapp para realizar tus gestiones</p>
          <a href="https://wa.me/5491122753000" style="color:#0099CC;font-size:15px;font-weight:bold;text-decoration:underline;font-family:'Montserrat',Arial,sans-serif;">112275-3000</a>
        </td>
      </tr>
    </table>

  </td>
</tr>

<!-- FOOTER -->
<tr>
  <td bgcolor="#1B5EA6" style="padding:20px 28px;text-align:center;">
    <p style="margin:0 0 4px 0;color:#ffffff;font-size:15px;font-weight:bold;font-family:'Montserrat',Arial,sans-serif;">El Jumillano</p>
    <p style="margin:0 0 10px 0;color:rgba(255,255,255,0.78);font-size:12px;font-family:'Montserrat',Arial,sans-serif;">Calidad, comodidad y puntualidad en cada entrega</p>
    <p style="margin:0;color:rgba(255,255,255,0.5);font-size:11px;font-family:'Montserrat',Arial,sans-serif;">Este es un correo electr&oacute;nico autom&aacute;tico, por favor no responder.</p>
  </td>
</tr>

</table>
</td></tr>
</table>
</body>
</html>`,
		montserratFontCSS,
		texts.subtitle,
		texts.title,
		texts.title, texts.subtitle,
		logoSrc,
		message,
		texts.cardTitle,
		orderNumber,
		delivery.Address,
		delivery.FechaAccion.Format("02/01/2006"),
		len(dispensers),
		texts.dispenserTitle,
		dispensersRows,
		texts.pdfNotice,
	)
}
