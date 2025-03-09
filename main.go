package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/RoMalms10/grpc/meme"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Default address for the meme generator service
const defaultMemeServiceAddr = "meme-generator:50051"

func main() {
	// Set up gRPC connection to the meme service
	memeServiceAddr := os.Getenv("MEME_SERVICE_ADDR")
	if memeServiceAddr == "" {
		memeServiceAddr = defaultMemeServiceAddr
	}

	// Set up the router
	r := gin.Default()

	// Connect to the meme service
	conn, err := grpc.Dial(memeServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to meme service: %v", err)
	}
	defer conn.Close()

	memeClient := pb.NewMemeServiceClient(conn)

	// Routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// List available meme templates
	r.GET("/templates", func(c *gin.Context) {
		category := c.Query("category")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := memeClient.ListTemplates(ctx, &pb.ListTemplatesRequest{
			Category: category,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	// Generate a meme
	r.POST("/generate", func(c *gin.Context) {
		var request struct {
			TemplateID     string   `json:"template_id" binding:"required"`
			TopText        string   `json:"top_text"`
			BottomText     string   `json:"bottom_text"`
			AdditionalText []string `json:"additional_text"`
			UseAICaption   bool     `json:"use_ai_caption"`
			CaptionPrompt  string   `json:"caption_prompt"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := memeClient.GenerateMeme(ctx, &pb.GenerateMemeRequest{
			TemplateId:     request.TemplateID,
			TopText:        request.TopText,
			BottomText:     request.BottomText,
			AdditionalText: request.AdditionalText,
			UseAiCaption:   request.UseAICaption,
			CaptionPrompt:  request.CaptionPrompt,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if resp.Error != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": resp.Error,
			})
			return
		}

		// If requested as JSON, return the full response
		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusOK, resp)
			return
		}

		// Otherwise, return the image directly
		imageData, err := base64.StdEncoding.DecodeString(resp.ImageData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to decode image data",
			})
			return
		}

		c.Data(http.StatusOK, resp.MimeType, imageData)
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Edge API server started, listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
