package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	pb "github.com/RoMalms10/grpc/meme"
	"github.com/RoMalms10/meme-generator/config"
	"github.com/RoMalms10/meme-generator/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestService creates a new MemeService with test configuration
func setupTestService(t *testing.T) *service.MemeService {
	// Create a temporary directory for test templates
	tempDir := t.TempDir()

	// Copy test template files to the temp directory
	setupTestTemplates(t, tempDir)

	// Create test configuration
	cfg := &config.Config{
		TemplateDir:     tempDir,
		FontFile:        "../testdata/impact.ttf", // You'll need to provide a test font file
		ImageQuality:    90,
		FontSize:        36,
		LineSpacing:     1.5,
		EnableAICaption: true,
	}

	// Create service with test templates
	s := &service.MemeService{
		Config: cfg,
		Templates: map[string]*service.TemplateInfo{
			"test-template": {
				Name:           "Test Template",
				TextFieldCount: 2,
				Category:       "test",
				Filename:       "test-template.jpg",
			},
		},
	}

	return s
}

// setupTestTemplates copies test image files to the temp directory
func setupTestTemplates(t *testing.T, tempDir string) {
	// Create test template file
	// In a real test, you'd have a sample image in testdata directory
	testImagePath := "../testdata/test-template.jpg"

	// Create the test template in temp directory
	err := copyFile(testImagePath, filepath.Join(tempDir, "test-template.jpg"))
	if err != nil {
		t.Skip("Test template file not found, skipping test")
	}
}

// copyFile is a helper function to copy a file
func copyFile(src, dst string) error {
	// For testing purposes, we could either:
	// 1. Actually copy a file from testdata directory
	// 2. Create a minimal valid JPEG/PNG file programmatically

	// This is a simplified implementation - in a real test, use proper file copying:
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

func TestMemeService_ListTemplates(t *testing.T) {
	// Create a test service
	s := setupTestService(t)

	tests := []struct {
		name     string
		category string
		want     int // Expected number of templates
	}{
		{
			name:     "All templates",
			category: "",
			want:     1,
		},
		{
			name:     "Test category",
			category: "test",
			want:     1,
		},
		{
			name:     "Non-existent category",
			category: "nonexistent",
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.ListTemplatesRequest{
				Category: tt.category,
			}

			resp, err := s.ListTemplates(context.Background(), req)
			require.NoError(t, err)

			assert.Equal(t, tt.want, len(resp.Templates))

			if tt.want > 0 {
				assert.Equal(t, "test-template", resp.Templates[0].Id)
				assert.Equal(t, "Test Template", resp.Templates[0].Name)
				assert.Equal(t, "test", resp.Templates[0].Category)
			}
		})
	}
}

func TestMemeService_GenerateMeme(t *testing.T) {
	// Skip this test if the test files don't exist
	if _, err := os.Stat("../testdata/test-template.jpg"); os.IsNotExist(err) {
		t.Skip("Test template file not found, skipping test")
	}

	if _, err := os.Stat("../testdata/impact.ttf"); os.IsNotExist(err) {
		t.Skip("Test font file not found, skipping test")
	}

	// Create a test service
	s := setupTestService(t)

	tests := []struct {
		name           string
		templateID     string
		topText        string
		bottomText     string
		additionalText []string
		useAICaption   bool
		expectError    bool
	}{
		{
			name:           "Valid template with text",
			templateID:     "test-template",
			topText:        "Top Text",
			bottomText:     "Bottom Text",
			additionalText: []string{},
			useAICaption:   false,
			expectError:    false,
		},
		{
			name:           "Invalid template",
			templateID:     "nonexistent",
			topText:        "Test",
			bottomText:     "Test",
			additionalText: []string{},
			useAICaption:   false,
			expectError:    true,
		},
		{
			name:           "AI Caption",
			templateID:     "test-template",
			topText:        "",
			bottomText:     "",
			additionalText: []string{},
			useAICaption:   true,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.GenerateMemeRequest{
				TemplateId:     tt.templateID,
				TopText:        tt.topText,
				BottomText:     tt.bottomText,
				AdditionalText: tt.additionalText,
				UseAiCaption:   tt.useAICaption,
			}

			resp, err := s.GenerateMeme(context.Background(), req)
			require.NoError(t, err)

			if tt.expectError {
				assert.NotEmpty(t, resp.Error)
				assert.Empty(t, resp.ImageData)
			} else {
				assert.Empty(t, resp.Error)
				assert.NotEmpty(t, resp.ImageData)
				assert.Equal(t, "image/jpeg", resp.MimeType)

				if tt.useAICaption {
					assert.NotEmpty(t, resp.GeneratedCaptions)
				}
			}
		})
	}
}
