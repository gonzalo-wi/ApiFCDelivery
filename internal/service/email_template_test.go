package service

import (
	"context"
	"os"
	"testing"
	"time"

	"GoFrioCalor/internal/models"

	"github.com/joho/godotenv"
)

// TestSendTemplateEmail envía un email de prueba para validar el diseño.
// Uso: go test ./internal/service/ -run TestSendTemplateEmail -v
// Requiere EMAIL_HOST, EMAIL_PORT, EMAIL_FROM, EMAIL_PASSWORD en .env
func TestSendTemplateEmail(t *testing.T) {
	_ = godotenv.Load("../../.env")

	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	if host == "" || from == "" || password == "" {
		t.Skip("EMAIL_HOST / EMAIL_FROM / EMAIL_PASSWORD no configurados en .env")
	}

	delivery := &models.Delivery{
		Name:        "g.e.v.n. s.a.",
		Address:     "Igral.J de la Cruz 1855",
		Locality:    "CABA",
		NroCta:      "12345",
		FechaAccion: models.CustomDate{Time: time.Now()},
		TipoEntrega: models.Instalacion,
	}

	svc := &mobileDeliveryService{}
	html := svc.buildCompletionEmailHTML(delivery, []string{"17B0200", "17B0201"}, "OT-14")

	emailSvc, err := NewSMTPEmailService(host, port, from, password, "gwinazki@el-jumillano.com.ar")
	if err != nil {
		t.Fatalf("init email service: %v", err)
	}

	err = emailSvc.SendHTMLEmailWithPDFBytesAndLogo(
		context.Background(),
		"gwinazki@el-jumillano.com.ar",
		"TEST - Instalación Completada (IVESS)",
		html,
		nil,
		"",
		"../../assets/images/blanco.png",
	)
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	t.Log("Email enviado correctamente")
}
