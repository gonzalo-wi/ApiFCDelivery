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
	parsed, err := time.Parse("2006-01-02", fechaStr)
	if err == nil {
		return models.CustomDate{Time: parsed}, nil
	}
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

// createItemDispensers crea items de dispensers basados en cantidades por tipo
func createItemDispensers(cantidadPie, cantidadMesada uint) []models.ItemDispenser {
	items := make([]models.ItemDispenser, 0, 2)
	if cantidadPie > 0 {
		items = append(items, models.ItemDispenser{
			Tipo:     models.TipoDispenserPie,
			Cantidad: cantidadPie,
		})
	}
	if cantidadMesada > 0 {
		items = append(items, models.ItemDispenser{
			Tipo:     models.TipoDispenserMesada,
			Cantidad: cantidadMesada,
		})
	}
	return items
}
