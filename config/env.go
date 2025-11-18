package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	Port        string
	CORSOrigins string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error cargando .env: %w", err)
	}

	config := &Config{
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		Port:        os.Getenv("PORT"),
		CORSOrigins: os.Getenv("CORS_ORIGINS"),
	}

	return config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

// GetCORSOrigins devuelve los orígenes permitidos para CORS como slice
// Si no está configurado, retorna localhost:5173 por defecto para desarrollo
func (c *Config) GetCORSOrigins() []string {
	if c.CORSOrigins == "" {
		return []string{"http://localhost:5173"}
	}
	// Separar por comas para múltiples orígenes
	origins := []string{}
	for i := 0; i < len(c.CORSOrigins); i++ {
		start := i
		for i < len(c.CORSOrigins) && c.CORSOrigins[i] != ',' {
			i++
		}
		origin := c.CORSOrigins[start:i]
		// Eliminar espacios
		for len(origin) > 0 && (origin[0] == ' ' || origin[0] == '\t') {
			origin = origin[1:]
		}
		for len(origin) > 0 && (origin[len(origin)-1] == ' ' || origin[len(origin)-1] == '\t') {
			origin = origin[:len(origin)-1]
		}
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}
