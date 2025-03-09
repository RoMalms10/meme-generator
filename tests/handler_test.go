package tests

import (
	"context"
	"testing"

	pb "github.com/RoMalms10/grpc/meme"
	"github.com/RoMalms10/meme-generator/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Define an interface that matches the methods we need from service.MemeService
type MemeServiceInterface interface {
	GenerateMeme(ctx context.Context, req *pb.GenerateMemeRequest) (*pb.GenerateMemeResponse, error)
	ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error)
}

// MockMemeService implements a mock of the MemeService for testing the handler
type MockMemeService struct {
	mock.Mock
}

// GenerateMeme mocks the service.GenerateMeme method
func (m *MockMemeService) GenerateMeme(ctx context.Context, req *pb.GenerateMemeRequest) (*pb.GenerateMemeResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.GenerateMemeResponse), args.Error(1)
}

// ListTemplates mocks the service.ListTemplates method
func (m *MockMemeService) ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.ListTemplatesResponse), args.Error(1)
}

// Modify the handler package to accept an interface instead of concrete type
func TestHandler_GenerateMeme(t *testing.T) {
	mockService := new(MockMemeService)
	// Create a handler that accepts our mock service
	h := handler.NewMemeHandler(mockService)

	tests := []struct {
		name           string
		templateID     string
		topText        string
		bottomText     string
		additionalText []string
		useAICaption   bool
		mockResponse   *pb.GenerateMemeResponse
		expectError    bool
	}{
		{
			name:           "Valid request",
			templateID:     "drake",
			topText:        "Top Text",
			bottomText:     "Bottom Text",
			additionalText: []string{},
			useAICaption:   false,
			mockResponse: &pb.GenerateMemeResponse{
				ImageData:         "base64encodedimage",
				MimeType:          "image/jpeg",
				GeneratedCaptions: []string{},
				Error:             "",
			},
			expectError: false,
		},
		{
			name:           "Missing template ID",
			templateID:     "",
			topText:        "Top Text",
			bottomText:     "Bottom Text",
			additionalText: []string{},
			useAICaption:   false,
			mockResponse:   nil, // Not used as handler will return early
			expectError:    false,
		},
		{
			name:           "Service error",
			templateID:     "error-template",
			topText:        "Top Text",
			bottomText:     "Bottom Text",
			additionalText: []string{},
			useAICaption:   false,
			mockResponse: &pb.GenerateMemeResponse{
				ImageData:         "",
				MimeType:          "",
				GeneratedCaptions: []string{},
				Error:             "Service error",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := &pb.GenerateMemeRequest{
				TemplateId:     tt.templateID,
				TopText:        tt.topText,
				BottomText:     tt.bottomText,
				AdditionalText: tt.additionalText,
				UseAiCaption:   tt.useAICaption,
			}

			// Set up mock expectations
			if tt.templateID != "" {
				mockService.On("GenerateMeme", mock.Anything, req).Return(tt.mockResponse, nil).Once()
			}

			// Call the handler
			resp, err := h.GenerateMeme(context.Background(), req)

			// Verify results
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.templateID == "" {
					// Should return error response for empty template ID
					assert.Equal(t, "template_id is required", resp.Error)
				} else {
					// Should return the mock response
					assert.Equal(t, tt.mockResponse, resp)
				}
			}

			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_ListTemplates(t *testing.T) {
	mockService := new(MockMemeService)
	h := handler.NewMemeHandler(mockService)

	tests := []struct {
		name       string
		category   string
		mockResult *pb.ListTemplatesResponse
	}{
		{
			name:     "All templates",
			category: "",
			mockResult: &pb.ListTemplatesResponse{
				Templates: []*pb.Template{
					{
						Id:             "drake",
						Name:           "Drake Hotline Bling",
						TextFieldCount: 2,
						Category:       "classic",
						PreviewUrl:     "/templates/drake.jpg",
					},
				},
			},
		},
		{
			name:     "Filtered by category",
			category: "classic",
			mockResult: &pb.ListTemplatesResponse{
				Templates: []*pb.Template{
					{
						Id:             "drake",
						Name:           "Drake Hotline Bling",
						TextFieldCount: 2,
						Category:       "classic",
						PreviewUrl:     "/templates/drake.jpg",
					},
				},
			},
		},
		{
			name:     "Empty result",
			category: "nonexistent",
			mockResult: &pb.ListTemplatesResponse{
				Templates: []*pb.Template{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := &pb.ListTemplatesRequest{
				Category: tt.category,
			}

			// Set up mock expectations
			mockService.On("ListTemplates", mock.Anything, req).Return(tt.mockResult, nil).Once()

			// Call the handler
			resp, err := h.ListTemplates(context.Background(), req)

			// Verify results
			require.NoError(t, err)
			assert.Equal(t, tt.mockResult, resp)

			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
