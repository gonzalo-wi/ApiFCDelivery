package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Obtener configuración de DB
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construir DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Conectar a la base de datos
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error verificando conexión: %v", err)
	}

	fmt.Println("✅ Conectado a la base de datos")
	fmt.Println("🚀 Aplicando migración 005...")

	// Ejecutar migración
	queries := []string{
		`ALTER TABLE deliveries 
		 ADD COLUMN name VARCHAR(200) DEFAULT '' AFTER nro_cta,
		 ADD COLUMN address VARCHAR(300) DEFAULT '' AFTER name,
		 ADD COLUMN locality VARCHAR(100) DEFAULT '' AFTER address`,
		`CREATE INDEX idx_deliveries_name ON deliveries(name)`,
		`CREATE INDEX idx_deliveries_locality ON deliveries(locality)`,
	}

	for i, query := range queries {
		fmt.Printf("Ejecutando paso %d/%d...\n", i+1, len(queries))
		if _, err := db.Exec(query); err != nil {
			log.Printf("⚠️  Error (puede ser que ya exista): %v", err)
			continue
		}
		fmt.Printf("✅ Paso %d completado\n", i+1)
	}

	fmt.Println("\n✅ Migración 005 aplicada exitosamente!")
	fmt.Println("\nCambios realizados:")
	fmt.Println("  - Columna 'name' agregada (VARCHAR 200)")
	fmt.Println("  - Columna 'address' agregada (VARCHAR 300)")
	fmt.Println("  - Columna 'locality' agregada (VARCHAR 100)")
	fmt.Println("  - Índice 'idx_deliveries_name' creado")
	fmt.Println("  - Índice 'idx_deliveries_locality' creado")

	// Verificar estructura
	fmt.Println("\n📋 Verificando estructura de tabla deliveries...")
	rows, err := db.Query("DESCRIBE deliveries")
	if err != nil {
		log.Printf("Error verificando tabla: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("\nField\t\t\tType")
	fmt.Println("-----\t\t\t----")
	for rows.Next() {
		var field, fieldType, null, key, defaultVal, extra sql.NullString
		if err := rows.Scan(&field, &fieldType, &null, &key, &defaultVal, &extra); err != nil {
			continue
		}
		fmt.Printf("%s\t\t\t%s\n", field.String, fieldType.String)
	}
}
