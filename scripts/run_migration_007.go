package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// Construir DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", "dispensers_db"),
	)

	// Conectar a la base de datos
	fmt.Println("Conectando a la base de datos...")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Printf("Error conectando a la base de datos: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("✓ Conexión exitosa")

	// Leer archivo de migración
	migrationPath := filepath.Join("migrations", "007_create_audit_events.sql")
	fmt.Printf("Leyendo migración: %s\n", migrationPath)

	sqlContent, err := os.ReadFile(migrationPath)
	if err != nil {
		fmt.Printf("Error leyendo archivo de migración: %v\n", err)
		os.Exit(1)
	}

	// Ejecutar migración
	fmt.Println("Ejecutando migración...")
	_, err = db.Exec(string(sqlContent))
	if err != nil {
		fmt.Printf("Error ejecutando migración: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ Migración 007 aplicada exitosamente!")
	fmt.Println("  - Tabla 'audit_events' creada")
	fmt.Println("  - 7 índices creados")
	fmt.Println("  - Función cleanup_old_audit_events() creada")
	fmt.Println("\nEl servidor está listo para usar el sistema de auditoría.")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
