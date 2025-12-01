package service

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFService interface {
	GenerateWorkOrderPDF(ctx context.Context, workOrder *dto.WorkOrderRequest) ([]byte, string, error)
}

type pdfService struct {
	workOrderStore store.WorkOrderStore
}

func NewPDFService(workOrderStore store.WorkOrderStore) PDFService {
	return &pdfService{workOrderStore: workOrderStore}
}

func (s *pdfService) GenerateWorkOrderPDF(ctx context.Context, workOrder *dto.WorkOrderRequest) ([]byte, string, error) {
	orderNumber, err := s.workOrderStore.GetNextOrderNumber(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get next order number: %w", err)
	}
	woModel := &models.WorkOrder{
		OrderNumber: orderNumber,
		NroCta:      workOrder.NroCta,
		NroRto:      workOrder.NroRto,
		Name:        workOrder.Name,
		Address:     workOrder.Address,
		Localidad:   workOrder.Locality,
		TipoAccion:  workOrder.TipoAccion,
	}
	if err := s.workOrderStore.Create(ctx, woModel); err != nil {
		return nil, "", fmt.Errorf("failed to create work order: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(15, 15, 15)

	colorPrimary := []int{41, 128, 185}
	colorText := []int{52, 73, 94}
	colorAccent := []int{52, 152, 219}

	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(0, 0, 210, 50, "F")
	pdf.SetFillColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.Rect(0, 48, 210, 4, "F")
	logoPath := "assets/images/logoivess.PNG"
	pdf.Image(logoPath, 15, 10, 40, 0, false, "", 0, "")
	pdf.SetTextColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetFont("Arial", "B", 16)
	pdf.SetY(35)
	pdf.SetX(15)
	pdf.CellFormat(60, 8, constants.PDFHeaderTitle, "", 0, "L", false, 0, "")
	pdf.SetDrawColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetLineWidth(0.5)
	pdf.Rect(140, 10, 55, 30, "D")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetY(13)
	pdf.SetX(140)
	pdf.CellFormat(55, 7, constants.PDFOrderTitle, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(colorAccent[0], colorAccent[1], colorAccent[2])
	pdf.SetY(22)
	pdf.SetX(140)
	pdf.CellFormat(55, 10, orderNumber, "", 0, "C", false, 0, "")
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.SetLineWidth(0.2)
	pdf.SetY(58)
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 8, constants.PDFSectionService, "", 1, "L", true, 0, "")
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.Ln(2)
	pdf.SetFont("Arial", "", 10)
	leftCol := 95.0

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(45, 7, constants.PDFLabelDate)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(leftCol-45, 7, time.Now().Format("02/01/2006"))
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(45, 7, constants.PDFLabelActionType)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 7, workOrder.TipoAccion)
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(45, 7, constants.PDFLabelAccountNumber)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(leftCol-45, 7, workOrder.NroCta)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(45, 7, constants.PDFLabelDeliveryNumber)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 7, workOrder.NroRto)
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 8, constants.PDFSectionClient, "", 1, "L", true, 0, "")
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.Ln(2)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(45, 7, constants.PDFLabelName)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 7, workOrder.Name)
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(45, 7, constants.PDFLabelAddress)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 7, workOrder.Address)
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(45, 7, constants.PDFLabelLocality)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 7, workOrder.Locality)
	pdf.Ln(10)

	if len(workOrder.Dispensers) > 0 {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetFillColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(0, 8, constants.PDFSectionEquipment, "", 1, "L", true, 0, "")
		pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
		pdf.Ln(2)
		pdf.SetFillColor(colorAccent[0], colorAccent[1], colorAccent[2])
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(20, 8, constants.PDFTableItem, "1", 0, "C", true, 0, "")
		pdf.CellFormat(60, 8, constants.PDFTableBrand, "1", 0, "C", true, 0, "")
		pdf.CellFormat(0, 8, constants.PDFTableSerialNumber, "1", 1, "C", true, 0, "")

		pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
		pdf.SetFont("Arial", "", 10)
		fill := false
		for i, dispenser := range workOrder.Dispensers {
			if fill {
				pdf.SetFillColor(245, 245, 245)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			pdf.CellFormat(20, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", fill, 0, "")
			pdf.CellFormat(60, 7, dispenser.Marca, "1", 0, "L", fill, 0, "")
			pdf.CellFormat(0, 7, dispenser.NroSerie, "1", 1, "L", fill, 0, "")
			fill = !fill
		}
		pdf.Ln(5)
	}
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 8, constants.PDFSectionTask, "", 1, "L", true, 0, "")
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.Ln(2)

	rectStartY := pdf.GetY()

	var tareaTexto string
	switch workOrder.TipoAccion {
	case "Instalacion":
		tareaTexto = constants.TaskInstallation
	case "Retiro":
		tareaTexto = constants.TaskRemoval
	case "Recambio":
		tareaTexto = constants.TaskReplacement
	default:
		tareaTexto = ""
	}

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.SetX(17)
	pdf.MultiCell(176, 5, tareaTexto, "", "L", false)

	pdf.SetDrawColor(colorAccent[0], colorAccent[1], colorAccent[2])
	pdf.SetLineWidth(0.3)
	pdf.Rect(15, rectStartY, 180, 35, "D")
	pdf.SetY(rectStartY + 35)
	pdf.Ln(2)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 8, constants.PDFSectionTerms, "", 1, "L", true, 0, "")
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(60, 60, 60)
	pdf.MultiCell(0, 3, constants.MsgTextAcepted, "", "J", false)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 8, constants.PDFSectionAcceptance, "", 1, "L", true, 0, "")
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.Ln(3)

	acceptanceStartY := pdf.GetY()
	pdf.SetDrawColor(colorAccent[0], colorAccent[1], colorAccent[2])
	pdf.SetLineWidth(0.5)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetX(17)
	pdf.Cell(50, 6, constants.PDFLabelAccepted)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 150, 0)
	pdf.Cell(0, 6, constants.PDFLabelAcceptedValue)
	pdf.Ln(7)

	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.SetFont("Arial", "B", 10)
	pdf.SetX(17)
	pdf.Cell(50, 6, constants.PDFLabelDateTime)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, time.Now().Format("02/01/2006 15:04"))
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetX(17)
	pdf.Cell(50, 6, constants.PDFLabelToken)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(colorAccent[0], colorAccent[1], colorAccent[2])
	tokenDisplay := workOrder.Token
	if tokenDisplay == "" {
		tokenDisplay = "N/A"
	}
	pdf.Cell(0, 6, tokenDisplay)
	pdf.Ln(2)

	// Dibujar recuadro
	pdf.Rect(15, acceptanceStartY, 180, 22, "D")
	// Ubicar el texto de IMPORTANTE más cerca del recuadro
	pdf.SetY(acceptanceStartY + 21)

	pdf.Ln(0)

	// Generar el PDF en memoria
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), orderNumber, nil
}
