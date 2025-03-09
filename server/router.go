package server

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/RoMalms10/grpc/meme"
	"github.com/RoMalms10/meme-generator/config"
	"github.com/RoMalms10/meme-generator/handler"
	"github.com/RoMalms10/meme-generator/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// NewGRPCServer creates and configures a gRPC server with all services registered
func NewGRPCServer(memeService *service.MemeService, cfg *config.Config) *grpc.Server {
	// Configure server parameters based on configuration
	keepaliveParams := keepalive.ServerParameters{
		MaxConnectionIdle:     cfg.GRPCMaxConnectionIdle,
		MaxConnectionAge:      cfg.GRPCMaxConnectionAge,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  cfg.GRPCKeepaliveTime,
		Timeout:               cfg.GRPCKeepaliveTimeout,
	}

	keepaliveEnforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}

	// Create a new gRPC server with configured options
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepaliveParams),
		grpc.KeepaliveEnforcementPolicy(keepaliveEnforcementPolicy),
	)

	// Register all services
	RegisterServices(grpcServer, memeService)

	// Enable server reflection for tools like grpcurl and debugging
	reflection.Register(grpcServer)

	return grpcServer
}

// StartServer starts the gRPC server on the specified port
func StartServer(server *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", port, err)
	}

	log.Printf("Meme Generator service started, listening on port %s", port)
	return server.Serve(lis)
}

// RegisterServices registers all available services with the gRPC server
func RegisterServices(server *grpc.Server, memeService *service.MemeService) {
	// Create handlers for each service
	memeHandler := handler.NewMemeHandler(memeService)

	// Register the MemeService with its endpoints
	// This makes the following gRPC methods available:
	// - GenerateMeme: Creates a meme with custom text based on a template
	// - ListTemplates: Returns available meme templates with filtering options
	pb.RegisterMemeServiceServer(server, memeHandler)

	// Log available endpoints for debugging purposes
	log.Println("Registered MemeService endpoints:")
	log.Println("- GenerateMeme: Creates a meme with the given parameters")
	log.Println("- ListTemplates: Returns a list of available meme templates")

	// This is where you would register additional services as needed
	// Example: pb.RegisterHealthServer(server, healthService)
}
