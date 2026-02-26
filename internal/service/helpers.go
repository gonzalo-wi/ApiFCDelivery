package service

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/models"
	"fmt"
	"time"
)

// parseFechaAccion parsea una fecha en formato YYYY-MM-DD o ISO 8601
// Retorna la fecha actual si no se proporciona una fecha
func parseFechaAccion(fechaStr string) (models.CustomDate, error) {
	if fechaStr == "" {
		return models.CustomDate{Time: time.Now()}, nil
	}

	// Intentar formato simple YYYY-MM-DD
	parsed, err := time.Parse("2006-01-02", fechaStr)
	if err == nil {
		return models.CustomDate{Time: parsed}, nil
	}

	// Intentar formato ISO 8601
	parsed, err = time.Parse(time.RFC3339, fechaStr)
	if err == nil {
		return models.CustomDate{Time: parsed}, nil
	}

	return models.CustomDate{}, fmt.Errorf("formato de fecha inválido, use YYYY-MM-DD o ISO 8601")
}

// validateDispenserQuantity valida que la cantidad esté dentro de los límites permitidos
func validateDispenserQuantity(cantidad uint) error {
	if cantidad < constants.MIN_DISPENSERS {
		return fmt.Errorf("debe especificar al menos %d dispenser", constants.MIN_DISPENSERS)
	}
	if cantidad > constants.MAX_DISPENSERS {
		return fmt.Errorf("la cantidad total de dispensers no puede superar %d", constants.MAX_DISPENSERS)
	}
	return nil
}

// createPlaceholderDispensers crea dispensers placeholder que serán completados posteriormente
func createPlaceholderDispensers(nroRto string, cantidadPie, cantidadMesada uint) []models.Dispenser {
	cantidadTotal := cantidadPie + cantidadMesada
	dispensers := make([]models.Dispenser, 0, cantidadTotal)

	// Crear dispensers de Pie
	for i := uint(0); i < cantidadPie; i++ {
		dispensers = append(dispensers, models.Dispenser{
			Marca:    constants.DISPENSER_MARCA_PENDIENTE,
			NroSerie: fmt.Sprintf("P-%s-%d", nroRto, i+1),
			Tipo:     models.TipoDispenserPie,
		})
	}

	// Crear dispensers de Mesada
	for i := uint(0); i < cantidadMesada; i++ {
		dispensers = append(dispensers, models.Dispenser{
			Marca:    constants.DISPENSER_MARCA_PENDIENTE,
			NroSerie: fmt.Sprintf("M-%s-%d", nroRto, i+1),
			Tipo:     models.TipoDispenserMesada,
		})
	}

	return dispensers
}
