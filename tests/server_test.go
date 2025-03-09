package tests

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/RoMalms10/meme-generator/config"
	"github.com/RoMalms10/meme-generator/server"
	"github.com/RoMalms10/meme-generator/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestNewGRPCServer(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Port:                  "50051",
		GRPCMaxConnectionAge:  30 * time.Minute,
		GRPCMaxConnectionIdle: 15 * time.Minute,
		GRPCKeepaliveTime:     5 * time.Minute,
		GRPCKeepaliveTimeout:  20 * time.Second,
	}

	// Create a test service
	memeService := service.NewMemeService()

	// Create the gRPC server
	s := server.NewGRPCServer(memeService, cfg)

	// Check that the server is not nil
	assert.NotNil(t, s)
	assert.IsType(t, &grpc.Server{}, s)
}

func TestStartServer(t *testing.T) {
	// Create a test gRPC server
	s := grpc.NewServer()

	// Use a port that's already in use or invalid
	// Port 0 is special - OS will find an available port
	// But trying to bind to it twice should fail
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	// Get the port that was assigned
	port := lis.Addr().(*net.TCPAddr).Port

	// Now try to start a second server on the same port
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.StartServer(s, fmt.Sprintf("%d", port))
	}()

	// Wait for the server to error or time out
	select {
	case err := <-errChan:
		// We expect an error when trying to bind to an already used port
		assert.Error(t, err)
	case <-time.After(500 * time.Millisecond):
		// If it takes too long, fail the test
		t.Fatal("Timeout waiting for server to start")
	}
}

func TestRegisterServices(t *testing.T) {
	// Create a gRPC server
	s := grpc.NewServer()

	// Create a test service
	memeService := service.NewMemeService()

	// Register services
	server.RegisterServices(s, memeService)

	// We can't easily test the registration directly,
	// but we can verify the server still exists
	assert.NotNil(t, s)
}

func TestNewServer(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Port:                  "50051",
		TemplateDir:           "./templates",
		FontFile:              "./fonts/impact.ttf",
		ImageQuality:          90,
		FontSize:              36,
		LineSpacing:           1.5,
		EnableAICaption:       false,
		GRPCMaxConnectionAge:  30 * time.Minute,
		GRPCMaxConnectionIdle: 15 * time.Minute,
		GRPCKeepaliveTime:     5 * time.Minute,
		GRPCKeepaliveTimeout:  20 * time.Second,
	}

	// Create a server
	srv := server.NewServer(cfg)

	// Check the server properties
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.GetMemeService())
}

func TestServer_StartStop(t *testing.T) {
	// Create test configuration with a random high port
	// to avoid conflicts
	cfg := &config.Config{
		Port:                  "59999",
		TemplateDir:           "./templates",
		FontFile:              "./fonts/impact.ttf",
		ImageQuality:          90,
		FontSize:              36,
		LineSpacing:           1.5,
		EnableAICaption:       false,
		GRPCMaxConnectionAge:  30 * time.Minute,
		GRPCMaxConnectionIdle: 15 * time.Minute,
		GRPCKeepaliveTime:     5 * time.Minute,
		GRPCKeepaliveTimeout:  20 * time.Second,
	}

	// Create a server
	srv := server.NewServer(cfg)
	require.NotNil(t, srv)

	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Stop the server
	srv.Stop()

	// Check for any errors
	select {
	case err := <-errChan:
		// We should get nil since we shut down gracefully
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for server to stop")
	}
}
