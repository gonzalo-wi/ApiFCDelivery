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

func TestCreatePlaceholderDispensers(t *testing.T) {
	tests := []struct {
		name           string
		nroRto         string
		cantidadPie    uint
		cantidadMesada uint
		wantTotal      int
		wantPie        int
		wantMesada     int
	}{
		{
			name:           "Solo dispensers de pie",
			nroRto:         "RTO001",
			cantidadPie:    2,
			cantidadMesada: 0,
			wantTotal:      2,
			wantPie:        2,
			wantMesada:     0,
		},
		{
			name:           "Solo dispensers de mesada",
			nroRto:         "RTO002",
			cantidadPie:    0,
			cantidadMesada: 3,
			wantTotal:      3,
			wantPie:        0,
			wantMesada:     3,
		},
		{
			name:           "Mixto: pie y mesada",
			nroRto:         "RTO003",
			cantidadPie:    2,
			cantidadMesada: 1,
			wantTotal:      3,
			wantPie:        2,
			wantMesada:     1,
		},
		{
			name:           "Sin dispensers",
			nroRto:         "RTO004",
			cantidadPie:    0,
			cantidadMesada: 0,
			wantTotal:      0,
			wantPie:        0,
			wantMesada:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispensers := createPlaceholderDispensers(tt.nroRto, tt.cantidadPie, tt.cantidadMesada)

			// Verificar cantidad total
			if len(dispensers) != tt.wantTotal {
				t.Errorf("createPlaceholderDispensers() cantidad = %v, want %v", len(dispensers), tt.wantTotal)
			}

			// Contar por tipo
			countPie := 0
			countMesada := 0

			for _, d := range dispensers {
				// Verificar marca PENDIENTE
				if d.Marca != constants.DISPENSER_MARCA_PENDIENTE {
					t.Errorf("createPlaceholderDispensers() marca = %v, want %v", d.Marca, constants.DISPENSER_MARCA_PENDIENTE)
				}

				// Verificar que el número de serie contiene el nroRto
				if d.NroSerie == "" {
					t.Errorf("createPlaceholderDispensers() NroSerie vacío")
				}

				// Contar tipos
				if d.Tipo == models.TipoDispenserPie {
					countPie++
				} else if d.Tipo == models.TipoDispenserMesada {
					countMesada++
				}
			}

			// Verificar conteo por tipo
			if countPie != tt.wantPie {
				t.Errorf("createPlaceholderDispensers() cantidad Pie = %v, want %v", countPie, tt.wantPie)
			}
			if countMesada != tt.wantMesada {
				t.Errorf("createPlaceholderDispensers() cantidad Mesada = %v, want %v", countMesada, tt.wantMesada)
			}
		})
	}
}

func TestCreatePlaceholderDispensersNroSerie(t *testing.T) {
	nroRto := "RTO999"
	dispensers := createPlaceholderDispensers(nroRto, 2, 2)

	// Verificar que los números de serie sean únicos y tengan formato correcto
	expectedSeries := []string{
		"P-RTO999-1",
		"P-RTO999-2",
		"M-RTO999-1",
		"M-RTO999-2",
	}

	for i, expected := range expectedSeries {
		if dispensers[i].NroSerie != expected {
			t.Errorf("Dispenser[%d].NroSerie = %v, want %v", i, dispensers[i].NroSerie, expected)
		}
	}
}
