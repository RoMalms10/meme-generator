package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RoMalms10/meme-generator/config"
	"github.com/RoMalms10/meme-generator/service"
	"google.golang.org/grpc"
)

// Server encapsulates the gRPC server and services
type Server struct {
	grpcServer  *grpc.Server
	memeService *service.MemeService
	config      *config.Config
}

// NewServer initializes a new server instance
func NewServer(cfg *config.Config) *Server {
	// Initialize the meme service with configuration
	memeService := service.NewMemeService()

	// Create and configure the gRPC server
	grpcServer := NewGRPCServer(memeService, cfg)

	return &Server{
		grpcServer:  grpcServer,
		memeService: memeService,
		config:      cfg,
	}
}

// Start begins the server and blocks until shutdown
func (s *Server) Start() error {
	// Create a context that will be canceled on termination signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- StartServer(s.grpcServer, s.config.Port)
	}()

	// Server is already configured with all services in NewGRPCServer

	// Wait for termination signal or server error
	select {
	case err := <-errChan:
		return err
	case <-signalChan:
		log.Println("Received termination signal, shutting down...")
		s.Stop()
		return nil
	case <-ctx.Done():
		log.Println("Context canceled, shutting down...")
		s.Stop()
		return nil
	}
}

// Stop gracefully shuts down the server
func (s *Server) Stop() {
	if s.grpcServer != nil {
		log.Println("Stopping gRPC server gracefully...")
		s.grpcServer.GracefulStop()
		log.Println("Server stopped")
	}
}

// GetMemeService returns the meme service instance for testing or configuration
func (s *Server) GetMemeService() *service.MemeService {
	return s.memeService
}
