package service

import (
	"context"
	"fmt"

	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"

	"github.com/rs/zerolog/log"
)

// EmailService maneja el envío de emails
type EmailService interface {
	SendWorkOrderEmail(ctx context.Context, workOrder *models.WorkOrder, pdfPath string) error
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

// WorkOrderPDFGenerator maneja la generación de PDFs para órdenes de trabajo
type WorkOrderPDFGenerator interface {
	GenerateWorkOrderPDF(ctx context.Context, workOrder *models.WorkOrder, dispensers []dto.DispenserMessage) (string, error)
}

// MockWorkOrderPDFGenerator es una implementación de prueba
type MockWorkOrderPDFGenerator struct{}

func NewMockPDFService() WorkOrderPDFGenerator {
	return &MockWorkOrderPDFGenerator{}
}

func (s *MockWorkOrderPDFGenerator) GenerateWorkOrderPDF(ctx context.Context, workOrder *models.WorkOrder, dispensers []dto.DispenserMessage) (string, error) {
	pdfPath := fmt.Sprintf("/tmp/work_order_%s.pdf", workOrder.OrderNumber)

	log.Info().
		Str("order_number", workOrder.OrderNumber).
		Str("pdf_path", pdfPath).
		Int("dispensers_count", len(dispensers)).
		Msg("📄 [MOCK] PDF generated (not really, this is a mock)")

	// TODO: Implementar generación real de PDF con gofpdf o similar
	// Por ahora retornamos una ruta ficticia
	return pdfPath, nil
}
