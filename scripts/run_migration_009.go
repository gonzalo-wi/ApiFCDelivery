package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Obtener credenciales de base de datos
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("❌ Database environment variables are incomplete")
	}

	// Construir DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Conectar a la base de datos
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Error connecting to database: %v", err)
	}
	defer db.Close()

	fmt.Println("========================================")
	fmt.Println("Applying Migration 009: Validated Dispensers")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Printf("Connected to: %s:%s/%s\n", dbHost, dbPort, dbName)
	fmt.Println()

	ctx := context.Background()

	// Aplicar migración
	fmt.Println("Applying migration...")

	migration := `
-- Add validated_dispensers column to deliveries table
ALTER TABLE deliveries 
ADD COLUMN IF NOT EXISTS validated_dispensers JSONB;

-- Add index for potential JSONB queries
CREATE INDEX IF NOT EXISTS idx_deliveries_validated_dispensers ON deliveries USING GIN (validated_dispensers);
	`

	_, err = db.ExecContext(ctx, migration)
	if err != nil {
		log.Fatalf("❌ Error applying migration: %v", err)
	}

	fmt.Println()
	fmt.Println("✅ Migration applied successfully!")
	fmt.Println()
	fmt.Println("Changes made:")
	fmt.Println("  • Added 'validated_dispensers' JSONB column to 'deliveries' table")
	fmt.Println("  • Created GIN index for JSONB queries on validated_dispensers")
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Migration completed")
	fmt.Println("========================================")
}
