package service

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

// EmailService maneja el envío de emails
type EmailService interface {
	SendWorkOrderEmail(ctx context.Context, workOrder *models.WorkOrder, pdfPath string) error
	SendHTMLEmail(ctx context.Context, to string, subject string, htmlBody string) error
	SendHTMLEmailWithPDFBytes(ctx context.Context, to string, subject string, htmlBody string, pdfBytes []byte, pdfFilename string) error
}

// SMTPEmailConfig contiene la configuración del servicio SMTP
type SMTPEmailConfig struct {
	Host     string
	Port     int
	From     string
	Password string
	To       string
}

// SMTPEmailService implementa el envío real de emails mediante SMTP
type SMTPEmailService struct {
	config SMTPEmailConfig
}

func NewSMTPEmailService(host string, port string, from string, password string, to string) (EmailService, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid email port: %w", err)
	}

	if host == "" || from == "" || password == "" || to == "" {
		return nil, fmt.Errorf("email configuration incomplete: host=%s, from=%s, to=%s", host, from, to)
	}

	return &SMTPEmailService{
		config: SMTPEmailConfig{
			Host:     host,
			Port:     portInt,
			From:     from,
			Password: password,
			To:       to,
		},
	}, nil
}

func (s *SMTPEmailService) SendWorkOrderEmail(ctx context.Context, workOrder *models.WorkOrder, pdfPath string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", s.config.To)
	m.SetHeader("Subject", fmt.Sprintf("Nueva Orden de Trabajo - %s", workOrder.OrderNumber))

	// Cuerpo del email
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Nueva Orden de Trabajo Generada</h2>
			<p><strong>Número de Orden:</strong> %s</p>
			<p><strong>Cliente:</strong> %s</p>
			<p><strong>Cuenta:</strong> %s</p>
			<p><strong>Ruta:</strong> %s</p>
			<p><strong>Dirección:</strong> %s</p>
			<p><strong>Localidad:</strong> %s</p>
			<p><strong>Tipo de Acción:</strong> %s</p>
			<p><strong>Fecha de Creación:</strong> %s</p>
			<br>
			<p>Adjunto encontrará el PDF con los detalles completos de la orden de trabajo.</p>
		</body>
		</html>
	`,
		workOrder.OrderNumber,
		workOrder.Name,
		workOrder.NroCta,
		workOrder.NroRto,
		workOrder.Address,
		workOrder.Localidad,
		workOrder.TipoAccion,
		workOrder.CreatedAt.Format("02/01/2006 15:04"),
	)

	m.SetBody("text/html", body)

	// Adjuntar PDF si existe
	if pdfPath != "" {
		m.Attach(pdfPath)
	}

	// Configurar dialer SMTP
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.From, s.config.Password)

	// Enviar email
	if err := d.DialAndSend(m); err != nil {
		log.Error().
			Err(err).
			Str("order_number", workOrder.OrderNumber).
			Str("to", s.config.To).
			Msg("Failed to send work order email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("order_number", workOrder.OrderNumber).
		Str("to", s.config.To).
		Str("pdf_attachment", pdfPath).
		Msg("📧 Work order email sent successfully")

	return nil
}

// SendHTMLEmail envía un email HTML genérico
func (s *SMTPEmailService) SendHTMLEmail(ctx context.Context, to string, subject string, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// Configurar dialer SMTP
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.From, s.config.Password)

	// Enviar email
	if err := d.DialAndSend(m); err != nil {
		log.Error().
			Err(err).
			Str("to", to).
			Str("subject", subject).
			Msg("Failed to send HTML email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("to", to).
		Str("subject", subject).
		Msg("📧 HTML email sent successfully")

	return nil
}

// SendHTMLEmailWithPDFBytes envía un email HTML con un PDF adjunto desde bytes
func (s *SMTPEmailService) SendHTMLEmailWithPDFBytes(ctx context.Context, to string, subject string, htmlBody string, pdfBytes []byte, pdfFilename string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// Adjuntar PDF desde bytes
	if len(pdfBytes) > 0 && pdfFilename != "" {
		m.Attach(pdfFilename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(pdfBytes)
			return err
		}))
	}

	// Configurar dialer SMTP
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.From, s.config.Password)

	// Enviar email
	if err := d.DialAndSend(m); err != nil {
		log.Error().
			Err(err).
			Str("to", to).
			Str("subject", subject).
			Str("pdf_attachment", pdfFilename).
			Msg("Failed to send HTML email with PDF")
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("to", to).
		Str("subject", subject).
		Str("pdf_attachment", pdfFilename).
		Int("pdf_size_bytes", len(pdfBytes)).
		Msg("📧 HTML email with PDF attachment sent successfully")

	return nil
}

// MockEmailService es una implementación de prueba que solo logea
type MockEmailService struct{}

func NewMockEmailService() EmailService {
	return &MockEmailService{}
}

func (s *MockEmailService) SendWorkOrderEmail(ctx context.Context, workOrder *models.WorkOrder, pdfPath string) error {
	log.Info().
		Str("order_number", workOrder.OrderNumber).
		Str("nro_cta", workOrder.NroCta).
		Str("name", workOrder.Name).
		Str("email_to", "TODO: agregar email del cliente").
		Str("pdf_attachment", pdfPath).
		Msg("📧 [MOCK] Email sent (not really, this is a mock)")

	// TODO: Implementar envío real con SMTP o servicio de email
	// Por ahora solo logueamos
	return nil
}

func (s *MockEmailService) SendHTMLEmail(ctx context.Context, to string, subject string, htmlBody string) error {
	log.Info().
		Str("to", to).
		Str("subject", subject).
		Msg("📧 [MOCK] HTML email sent (not really, this is a mock)")
	return nil
}

func (s *MockEmailService) SendHTMLEmailWithPDFBytes(ctx context.Context, to string, subject string, htmlBody string, pdfBytes []byte, pdfFilename string) error {
	log.Info().
		Str("to", to).
		Str("subject", subject).
		Str("pdf_attachment", pdfFilename).
		Int("pdf_size_bytes", len(pdfBytes)).
		Msg("📧 [MOCK] HTML email with PDF sent (not really, this is a mock)")
	return nil
}

// WorkOrderPDFGenerator maneja la generación de PDFs para órdenes de trabajo
type WorkOrderPDFGenerator interface {
	GenerateWorkOrderPDF(ctx context.Context, workOrder *models.WorkOrder, operations []dto.OperationMessage) (string, error)
}

// MockWorkOrderPDFGenerator es una implementación de prueba
type MockWorkOrderPDFGenerator struct{}

func NewMockPDFService() WorkOrderPDFGenerator {
	return &MockWorkOrderPDFGenerator{}
}

func (s *MockWorkOrderPDFGenerator) GenerateWorkOrderPDF(ctx context.Context, workOrder *models.WorkOrder, operations []dto.OperationMessage) (string, error) {
	pdfPath := fmt.Sprintf("/tmp/work_order_%s.pdf", workOrder.OrderNumber)

	log.Info().
		Str("order_number", workOrder.OrderNumber).
		Str("pdf_path", pdfPath).
		Int("operations_count", len(operations)).
		Msg("📄 [MOCK] PDF generated (not really, this is a mock)")

	// TODO: Implementar generación real de PDF con gofpdf o similar
	// Por ahora retornamos una ruta ficticia
	return pdfPath, nil
}
