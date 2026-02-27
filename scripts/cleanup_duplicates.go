package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Delivery struct {
	ID        int    `gorm:"primaryKey"`
	SessionID string `gorm:"index"`
}

func main() {
	// Cargar .env
	_ = godotenv.Load()

	// Construir DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Conectar a BD
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a BD:", err)
	}

	fmt.Println("=== Limpieza de Session IDs Duplicados ===")

	// 1. Encontrar duplicados
	var duplicates []struct {
		SessionID string
		Count     int64
		MinID     int
	}

	err = db.Table("deliveries").
		Select("session_id, COUNT(*) as count, MIN(id) as min_id").
		Where("session_id IS NOT NULL AND session_id != ''").
		Group("session_id").
		Having("COUNT(*) > 1").
		Scan(&duplicates).Error

	if err != nil {
		log.Fatal("Error buscando duplicados:", err)
	}

	if len(duplicates) == 0 {
		fmt.Println("✓ No se encontraron session_ids duplicados")
		fmt.Println("La base de datos está lista para el índice único")
		return
	}

	fmt.Printf("\n⚠️  Encontrados %d session_ids duplicados\n\n", len(duplicates))

	// 2. Limpiar duplicados (setear session_id = NULL excepto el más antiguo)
	for _, dup := range duplicates {
		fmt.Printf("Limpiando duplicados de session_id: %s (conservando ID: %d)\n", dup.SessionID, dup.MinID)

		result := db.Table("deliveries").
			Where("session_id = ? AND id > ?", dup.SessionID, dup.MinID).
			Update("session_id", nil)

		if result.Error != nil {
			log.Printf("Error limpiando session_id %s: %v\n", dup.SessionID, result.Error)
		} else {
			fmt.Printf("  → %d registros actualizados\n", result.RowsAffected)
		}
	}

	fmt.Println("\n✓ Limpieza completada")
	fmt.Println("Ahora puedes reiniciar el servidor para aplicar el índice único")
}
