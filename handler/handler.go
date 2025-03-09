package handler

import (
	"context"
	"log"

	pb "github.com/RoMalms10/grpc/meme"
)

// MemeServiceInterface defines the interface that the handler needs
type MemeServiceInterface interface {
	GenerateMeme(ctx context.Context, req *pb.GenerateMemeRequest) (*pb.GenerateMemeResponse, error)
	ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error)
}

// MemeHandler handles gRPC requests for the meme service
type MemeHandler struct {
	memeService MemeServiceInterface
	pb.UnimplementedMemeServiceServer
}

// NewMemeHandler creates a new handler with the given service
func NewMemeHandler(memeService MemeServiceInterface) *MemeHandler {
	return &MemeHandler{
		memeService: memeService,
	}
}

// GenerateMeme handles requests to generate a meme
func (h *MemeHandler) GenerateMeme(ctx context.Context, req *pb.GenerateMemeRequest) (*pb.GenerateMemeResponse, error) {
	log.Printf("Handler: Received meme generation request for template: %s", req.TemplateId)

	// Validate the request
	if req.TemplateId == "" {
		return &pb.GenerateMemeResponse{
			Error: "template_id is required",
		}, nil
	}

	// Call the service layer
	return h.memeService.GenerateMeme(ctx, req)
}

// ListTemplates handles requests to list available meme templates
func (h *MemeHandler) ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error) {
	log.Printf("Handler: Received request to list templates, category filter: %s", req.Category)

	// Call the service layer
	return h.memeService.ListTemplates(ctx, req)
}
