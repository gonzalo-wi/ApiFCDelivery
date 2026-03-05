package service

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"testing"
	"time"
)

func TestParseFechaAccion(t *testing.T) {
	tests := []struct {
		name      string
		fechaStr  string
		wantError bool
	}{
		{
			name:      "Fecha vacía retorna fecha actual",
			fechaStr:  "",
			wantError: false,
		},
		{
			name:      "Formato YYYY-MM-DD válido",
			fechaStr:  "2026-02-25",
			wantError: false,
		},
		{
			name:      "Formato ISO 8601 válido",
			fechaStr:  "2026-02-25T10:30:00Z",
			wantError: false,
		},
		{
			name:      "Formato inválido",
			fechaStr:  "25/02/2026",
			wantError: true,
		},
		{
			name:      "Fecha inválida",
			fechaStr:  "2026-13-45",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFechaAccion(tt.fechaStr)

			if tt.wantError {
				if err == nil {
					t.Errorf("parseFechaAccion() esperaba error pero no lo obtuvo")
				}
			} else {
				if err != nil {
					t.Errorf("parseFechaAccion() error = %v, wantError %v", err, tt.wantError)
					return
				}

				if tt.fechaStr == "" {
					// Debe ser fecha actual (con margen de 1 segundo)
					now := time.Now()
					diff := now.Sub(result.Time)
					if diff < 0 {
						diff = -diff
					}
					if diff > time.Second {
						t.Errorf("parseFechaAccion() con fecha vacía debería retornar fecha actual")
					}
				} else {
					if result.Time.IsZero() {
						t.Errorf("parseFechaAccion() retornó fecha zero")
					}
				}
			}
		})
	}
}

func TestValidateDispenserQuantity(t *testing.T) {
	tests := []struct {
		name      string
		cantidad  uint
		wantError bool
	}{
		{
			name:      "Cantidad cero es inválida",
			cantidad:  0,
			wantError: true,
		},
		{
			name:      "Cantidad mínima válida (1)",
			cantidad:  constants.MIN_DISPENSERS,
			wantError: false,
		},
		{
			name:      "Cantidad media válida (5)",
			cantidad:  5,
			wantError: false,
		},
		{
			name:      "Cantidad máxima válida (10)",
			cantidad:  constants.MAX_DISPENSERS,
			wantError: false,
		},
		{
			name:      "Cantidad sobre el máximo (11)",
			cantidad:  11,
			wantError: true,
		},
		{
			name:      "Cantidad muy alta (100)",
			cantidad:  100,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDispenserQuantity(tt.cantidad)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateDispenserQuantity() esperaba error pero no lo obtuvo")
				}
			} else {
				if err != nil {
					t.Errorf("validateDispenserQuantity() error = %v, wantError %v", err, tt.wantError)
				}
			}
		})
	}
}

func TestCreateItemDispensers(t *testing.T) {
	tests := []struct {
		name           string
		cantidadPie    uint
		cantidadMesada uint
		wantTotal      int
		wantPie        uint
		wantMesada     uint
	}{
		{
			name:           "Solo dispensers de pie",
			cantidadPie:    2,
			cantidadMesada: 0,
			wantTotal:      1,
			wantPie:        2,
			wantMesada:     0,
		},
		{
			name:           "Solo dispensers de mesada",
			cantidadPie:    0,
			cantidadMesada: 3,
			wantTotal:      1,
			wantPie:        0,
			wantMesada:     3,
		},
		{
			name:           "Mixto: pie y mesada",
			cantidadPie:    2,
			cantidadMesada: 1,
			wantTotal:      2,
			wantPie:        2,
			wantMesada:     1,
		},
		{
			name:           "Sin dispensers",
			cantidadPie:    0,
			cantidadMesada: 0,
			wantTotal:      0,
			wantPie:        0,
			wantMesada:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := createItemDispensers(tt.cantidadPie, tt.cantidadMesada)

			// Verificar cantidad total de items
			if len(items) != tt.wantTotal {
				t.Errorf("createItemDispensers() cantidad items = %v, want %v", len(items), tt.wantTotal)
			}

			// Verificar cantidades por tipo
			foundPie := uint(0)
			foundMesada := uint(0)

			for _, item := range items {
				if item.Tipo == models.TipoDispenserPie {
					foundPie = item.Cantidad
				} else if item.Tipo == models.TipoDispenserMesada {
					foundMesada = item.Cantidad
				}
			}

			if foundPie != tt.wantPie {
				t.Errorf("createItemDispensers() cantidad Pie = %v, want %v", foundPie, tt.wantPie)
			}
			if foundMesada != tt.wantMesada {
				t.Errorf("createItemDispensers() cantidad Mesada = %v, want %v", foundMesada, tt.wantMesada)
			}
		})
	}
}
