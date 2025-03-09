package tests

import (
	"os"
	"testing"
	"time"

	"github.com/RoMalms10/meme-generator/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Test with default values
	t.Run("Default values", func(t *testing.T) {
		// Clear any existing env vars that might interfere
		os.Unsetenv("PORT")
		os.Unsetenv("TEMPLATE_DIR")
		os.Unsetenv("FONT_FILE")
		os.Unsetenv("IMAGE_QUALITY")
		os.Unsetenv("FONT_SIZE")
		os.Unsetenv("LINE_SPACING")
		os.Unsetenv("ENABLE_AI_CAPTION")

		cfg := config.LoadConfig()

		// Check default values
		assert.Equal(t, "50051", cfg.Port)
		assert.Equal(t, "./templates", cfg.TemplateDir)
		assert.Equal(t, "./fonts/impact.ttf", cfg.FontFile)
		assert.Equal(t, 90, cfg.ImageQuality)
		assert.Equal(t, 36.0, cfg.FontSize)
		assert.Equal(t, 1.5, cfg.LineSpacing)
		assert.False(t, cfg.EnableAICaption)
	})

	// Test with custom values
	t.Run("Custom values", func(t *testing.T) {
		// Set environment variables
		os.Setenv("PORT", "8080")
		os.Setenv("TEMPLATE_DIR", "/custom/templates")
		os.Setenv("FONT_FILE", "/custom/fonts/comic.ttf")
		os.Setenv("IMAGE_QUALITY", "75")
		os.Setenv("FONT_SIZE", "42")
		os.Setenv("LINE_SPACING", "2.0")
		os.Setenv("ENABLE_AI_CAPTION", "true")
		os.Setenv("GRPC_MAX_CONNECTION_AGE", "1h")

		cfg := config.LoadConfig()

		// Check custom values
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "/custom/templates", cfg.TemplateDir)
		assert.Equal(t, "/custom/fonts/comic.ttf", cfg.FontFile)
		assert.Equal(t, 75, cfg.ImageQuality)
		assert.Equal(t, 42.0, cfg.FontSize)
		assert.Equal(t, 2.0, cfg.LineSpacing)
		assert.True(t, cfg.EnableAICaption)
		assert.Equal(t, 1*time.Hour, cfg.GRPCMaxConnectionAge)

		// Clean up environment
		os.Unsetenv("PORT")
		os.Unsetenv("TEMPLATE_DIR")
		os.Unsetenv("FONT_FILE")
		os.Unsetenv("IMAGE_QUALITY")
		os.Unsetenv("FONT_SIZE")
		os.Unsetenv("LINE_SPACING")
		os.Unsetenv("ENABLE_AI_CAPTION")
		os.Unsetenv("GRPC_MAX_CONNECTION_AGE")
	})

	// Test with invalid values (should use defaults)
	t.Run("Invalid values", func(t *testing.T) {
		os.Setenv("IMAGE_QUALITY", "invalid")
		os.Setenv("FONT_SIZE", "invalid")
		os.Setenv("LINE_SPACING", "invalid")
		os.Setenv("ENABLE_AI_CAPTION", "invalid")
		os.Setenv("GRPC_MAX_CONNECTION_AGE", "invalid")

		cfg := config.LoadConfig()

		// Should use default values for invalid inputs
		assert.Equal(t, 90, cfg.ImageQuality)
		assert.Equal(t, 36.0, cfg.FontSize)
		assert.Equal(t, 1.5, cfg.LineSpacing)
		assert.False(t, cfg.EnableAICaption)
		assert.Equal(t, 30*time.Minute, cfg.GRPCMaxConnectionAge)

		// Clean up environment
		os.Unsetenv("IMAGE_QUALITY")
		os.Unsetenv("FONT_SIZE")
		os.Unsetenv("LINE_SPACING")
		os.Unsetenv("ENABLE_AI_CAPTION")
		os.Unsetenv("GRPC_MAX_CONNECTION_AGE")
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("getEnv", func(t *testing.T) {
		os.Setenv("TEST_ENV_VAR", "test-value")
		defer os.Unsetenv("TEST_ENV_VAR")

		// Test existing env var
		assert.Equal(t, "test-value", config.GetEnv("TEST_ENV_VAR", "default"))

		// Test non-existent env var
		assert.Equal(t, "default", config.GetEnv("NON_EXISTENT_VAR", "default"))
	})

	t.Run("getIntEnv", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		os.Setenv("TEST_INVALID_INT", "not-an-int")
		defer func() {
			os.Unsetenv("TEST_INT")
			os.Unsetenv("TEST_INVALID_INT")
		}()

		// Test valid int
		assert.Equal(t, 42, config.GetIntEnv("TEST_INT", 10))

		// Test invalid int
		assert.Equal(t, 10, config.GetIntEnv("TEST_INVALID_INT", 10))

		// Test non-existent env var
		assert.Equal(t, 10, config.GetIntEnv("NON_EXISTENT_VAR", 10))
	})

	// Add more tests for other helper functions as needed
}
