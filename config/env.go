package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port                  string
	GRPCMaxConnectionAge  time.Duration
	GRPCMaxConnectionIdle time.Duration
	GRPCKeepaliveTime     time.Duration
	GRPCKeepaliveTimeout  time.Duration

	// Service configuration
	TemplateDir  string
	FontFile     string
	ImageQuality int
	FontSize     float64
	LineSpacing  float64

	// Feature flags
	EnableAICaption bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		// Server defaults
		Port:                  GetEnv("PORT", "50051"),
		GRPCMaxConnectionAge:  GetDurationEnv("GRPC_MAX_CONNECTION_AGE", 30*time.Minute),
		GRPCMaxConnectionIdle: GetDurationEnv("GRPC_MAX_CONNECTION_IDLE", 15*time.Minute),
		GRPCKeepaliveTime:     GetDurationEnv("GRPC_KEEPALIVE_TIME", 5*time.Minute),
		GRPCKeepaliveTimeout:  GetDurationEnv("GRPC_KEEPALIVE_TIMEOUT", 20*time.Second),

		// Service defaults
		TemplateDir:  GetEnv("TEMPLATE_DIR", "./templates"),
		FontFile:     GetEnv("FONT_FILE", "./fonts/impact.ttf"),
		ImageQuality: GetIntEnv("IMAGE_QUALITY", 90),
		FontSize:     GetFloatEnv("FONT_SIZE", 36),
		LineSpacing:  GetFloatEnv("LINE_SPACING", 1.5),

		// Feature flags
		EnableAICaption: GetBoolEnv("ENABLE_AI_CAPTION", false),
	}

	// Log configuration
	log.Println("Configuration loaded:")
	log.Printf("- Server port: %s", cfg.Port)
	log.Printf("- Template directory: %s", cfg.TemplateDir)
	log.Printf("- AI caption enabled: %v", cfg.EnableAICaption)

	return cfg
}

// Helper functions to get environment variables with defaults - exported for testing

// GetEnv retrieves an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetIntEnv retrieves an environment variable as an integer or returns a default value
func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return intValue
}

// GetFloatEnv retrieves an environment variable as a float or returns a default value
func GetFloatEnv(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Warning: Invalid float value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return floatValue
}

// GetBoolEnv retrieves an environment variable as a boolean or returns a default value
func GetBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: Invalid boolean value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return boolValue
}

// GetDurationEnv retrieves an environment variable as a duration or returns a default value
func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	durationValue, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Warning: Invalid duration value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return durationValue
}
