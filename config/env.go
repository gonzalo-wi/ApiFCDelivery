package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	Port           string
	CORSOrigins    string
	Environment    string
	InfobipBaseURL string
	InfobipAPIKey  string
	AppBaseURL     string
	TermsTTLHours  int
	AuthServiceURL string
	EmailHost      string
	EmailPort      string
	EmailFrom      string
	EmailPassword  string
	EmailTo        string
}

func LoadConfig() (*Config, error) {
	// Intentar cargar .env pero no fallar si no existe (Docker inyecta variables)
	_ = godotenv.Load()

	config := &Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		Port:           os.Getenv("PORT"),
		CORSOrigins:    os.Getenv("CORS_ORIGINS"),
		Environment:    getEnvOrDefault("ENVIRONMENT", "development"),
		InfobipBaseURL: getEnvOrDefault("INFOBIP_BASE_URL", "https://api2.infobip.com"),
		InfobipAPIKey:  os.Getenv("INFOBIP_API_KEY"),
		AppBaseURL:     getEnvOrDefault("APP_BASE_URL", "http://localhost:5173"),
		TermsTTLHours:  getEnvAsInt("TERMS_TTL_HOURS", 48),
		AuthServiceURL: getEnvOrDefault("AUTH_SERVICE_URL", "http://192.168.0.55:8087"),
		EmailHost:      getEnvOrDefault("EMAIL_HOST", "smtp.gmail.com"),
		EmailPort:      getEnvOrDefault("EMAIL_PORT", "587"),
		EmailFrom:      os.Getenv("EMAIL_FROM"),
		EmailPassword:  os.Getenv("EMAIL_PASSWORD"),
		EmailTo:        os.Getenv("EMAIL_TO"),
	}

	return config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}

func (c *Config) GetCORSOrigins() []string {
	if c.CORSOrigins == "" {
		return []string{"http://localhost:5173"}
	}
	origins := []string{}
	for i := 0; i < len(c.CORSOrigins); i++ {
		start := i
		for i < len(c.CORSOrigins) && c.CORSOrigins[i] != ',' {
			i++
		}
		origin := c.CORSOrigins[start:i]
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

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	var value int
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return defaultValue
	}
	return value
}
