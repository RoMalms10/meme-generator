package tests

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	pb "github.com/RoMalms10/grpc/meme"
	"github.com/RoMalms10/meme-generator/config"
	"github.com/RoMalms10/meme-generator/handler"
	"github.com/RoMalms10/meme-generator/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

// setupFunctionalTest sets up a real gRPC server using bufconn for in-memory communication
func setupFunctionalTest(t *testing.T) (pb.MemeServiceClient, func()) {
	// Create a new buffer listener
	lis = bufconn.Listen(bufSize)

	// Set up test directories
	tempDir := t.TempDir()
	setupTestTemplates(t, tempDir)

	// Create test configuration
	cfg := &config.Config{
		TemplateDir:     tempDir,
		FontFile:        "../testdata/impact.ttf",
		ImageQuality:    90,
		FontSize:        36,
		LineSpacing:     1.5,
		EnableAICaption: true,
		Port:            "0", // Not used with bufconn
	}

	// Create a service with test templates
	memeService := &service.MemeService{
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

	// Create the handler
	memeHandler := handler.NewMemeHandler(memeService)

	// Create a gRPC server
	s := grpc.NewServer()
	pb.RegisterMemeServiceServer(s, memeHandler)

	// Start the server
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

	// Create a client connection
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	// Create a client
	client := pb.NewMemeServiceClient(conn)

	// Return the client and a cleanup function
	return client, func() {
		conn.Close()
		s.Stop()
	}
}

func TestFunctional_ListTemplates(t *testing.T) {
	// Skip if test data is missing
	if _, err := os.Stat("../testdata/test-template.jpg"); os.IsNotExist(err) {
		t.Skip("Test template file not found, skipping test")
	}

	client, cleanup := setupFunctionalTest(t)
	defer cleanup()

	tests := []struct {
		name     string
		category string
		want     int
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &pb.ListTemplatesRequest{
				Category: tt.category,
			}

			resp, err := client.ListTemplates(ctx, req)
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

func TestFunctional_GenerateMeme(t *testing.T) {
	// Skip if test data is missing
	if _, err := os.Stat("../testdata/test-template.jpg"); os.IsNotExist(err) {
		t.Skip("Test template file not found, skipping test")
	}

	if _, err := os.Stat("../testdata/impact.ttf"); os.IsNotExist(err) {
		t.Skip("Test font file not found, skipping test")
	}

	client, cleanup := setupFunctionalTest(t)
	defer cleanup()

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
			name:           "Empty template ID",
			templateID:     "",
			topText:        "Test",
			bottomText:     "Test",
			additionalText: []string{},
			useAICaption:   false,
			expectError:    true, // Handler validates template_id is required
		},
		{
			name:           "Invalid template",
			templateID:     "nonexistent",
			topText:        "Test",
			bottomText:     "Test",
			additionalText: []string{},
			useAICaption:   false,
			expectError:    true, // Expected to return error in response
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &pb.GenerateMemeRequest{
				TemplateId:     tt.templateID,
				TopText:        tt.topText,
				BottomText:     tt.bottomText,
				AdditionalText: tt.additionalText,
				UseAiCaption:   tt.useAICaption,
			}

			resp, err := client.GenerateMeme(ctx, req)

			if tt.expectError && tt.templateID == "" {
				// For empty template ID, we expect handler validation to return an error
				assert.NotEmpty(t, resp.Error)
			} else if tt.expectError {
				// For other expected errors, we expect error in response
				assert.NotEmpty(t, resp.Error)
				assert.Empty(t, resp.ImageData)
			} else {
				// For success cases
				require.NoError(t, err)
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
