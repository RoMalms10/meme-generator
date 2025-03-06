package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/RoMalms10/grpc/proto/meme"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Default port for gRPC server
const defaultPort = "50051"

// Server implements the MemeService gRPC service
type memeServer struct {
	pb.UnimplementedMemeServiceServer
	templates map[string]*templateInfo
}

type templateInfo struct {
	name           string
	textFieldCount int32
	category       string
	filename       string
}

func newMemeServer() *memeServer {
	// Initialize with some default templates
	templates := map[string]*templateInfo{
		"drake": {
			name:           "Drake Hotline Bling",
			textFieldCount: 2,
			category:       "classic",
			filename:       "drake.jpg",
		},
		"distracted-boyfriend": {
			name:           "Distracted Boyfriend",
			textFieldCount: 3,
			category:       "classic",
			filename:       "distracted-boyfriend.jpg",
		},
		"two-buttons": {
			name:           "Two Buttons",
			textFieldCount: 3,
			category:       "classic",
			filename:       "two-buttons.jpg",
		},
		"change-my-mind": {
			name:           "Change My Mind",
			textFieldCount: 1,
			category:       "debate",
			filename:       "change-my-mind.jpg",
		},
	}

	return &memeServer{
		templates: templates,
	}
}

// GenerateMeme implements the GenerateMeme RPC method
func (s *memeServer) GenerateMeme(ctx context.Context, req *pb.GenerateMemeRequest) (*pb.GenerateMemeResponse, error) {
	log.Printf("Received meme generation request for template: %s", req.TemplateId)

	// Check if the template exists
	_, exists := s.templates[req.TemplateId]
	if !exists {
		return &pb.GenerateMemeResponse{
			Error: fmt.Sprintf("Template '%s' not found", req.TemplateId),
		}, nil
	}

	// TODO: Implement actual meme generation logic
	// For now, we'll just return a placeholder response

	// In a real implementation:
	// 1. Load the template image
	// 2. Add the text from the request
	// 3. Return the resulting image

	return &pb.GenerateMemeResponse{
		ImageData: "base64_encoded_image_would_go_here",
		MimeType:  "image/jpeg",
		Error:     "",
	}, nil
}

// ListTemplates implements the ListTemplates RPC method
func (s *memeServer) ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error) {
	log.Printf("Received request to list templates, category filter: %s", req.Category)

	var templates []*pb.Template

	for id, info := range s.templates {
		// Apply category filter if specified
		if req.Category != "" && info.category != req.Category {
			continue
		}

		templates = append(templates, &pb.Template{
			Id:             id,
			Name:           info.name,
			TextFieldCount: info.textFieldCount,
			Category:       info.category,
			PreviewUrl:     fmt.Sprintf("/templates/%s", info.filename),
		})
	}

	return &pb.ListTemplatesResponse{
		Templates: templates,
	}, nil
}

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMemeServiceServer(s, newMemeServer())

	// Register reflection service for easier debugging
	reflection.Register(s)

	log.Printf("Meme Generator service started, listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
