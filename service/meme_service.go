package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/RoMalms10/grpc/meme"
	"github.com/RoMalms10/meme-generator/config"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// MemeService configuration values will be provided via config package

// TemplateInfo represents information about a meme template
type TemplateInfo struct {
	Name           string
	TextFieldCount int32
	Category       string
	Filename       string
}

// MemeService handles the business logic for meme generation
type MemeService struct {
	Templates map[string]*TemplateInfo
	Config    *config.Config
}

// NewMemeService creates a new instance of the meme service
func NewMemeService() *MemeService {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize with some default templates
	templates := map[string]*TemplateInfo{
		"drake": {
			Name:           "Drake Hotline Bling",
			TextFieldCount: 2,
			Category:       "classic",
			Filename:       "drake.jpg",
		},
		"distracted-boyfriend": {
			Name:           "Distracted Boyfriend",
			TextFieldCount: 3,
			Category:       "classic",
			Filename:       "distracted-boyfriend.jpg",
		},
		"two-buttons": {
			Name:           "Two Buttons",
			TextFieldCount: 3,
			Category:       "classic",
			Filename:       "two-buttons.jpg",
		},
		"change-my-mind": {
			Name:           "Change My Mind",
			TextFieldCount: 1,
			Category:       "debate",
			Filename:       "change-my-mind.jpg",
		},
	}

	return &MemeService{
		Templates: templates,
		Config:    cfg,
	}
}

// GenerateMeme creates a meme with the given parameters
func (s *MemeService) GenerateMeme(ctx context.Context, req *pb.GenerateMemeRequest) (*pb.GenerateMemeResponse, error) {
	log.Printf("Service: Processing meme generation for template: %s", req.TemplateId)

	// Check if the template exists
	_, exists := s.Templates[req.TemplateId]
	if !exists {
		return &pb.GenerateMemeResponse{
			Error: fmt.Sprintf("Template '%s' not found", req.TemplateId),
		}, nil
	}

	// Handle AI caption generation if requested
	var generatedCaptions []string
	if req.UseAiCaption {
		// In a real implementation, you'd call an AI service here
		// For now, we'll just generate something simple based on the template
		template := s.Templates[req.TemplateId]
		generatedCaptions = []string{
			fmt.Sprintf("AI generated caption for %s meme", template.Name),
		}

		// Use the generated caption if no text was provided
		if req.TopText == "" && len(generatedCaptions) > 0 {
			req.TopText = generatedCaptions[0]
		}
	}

	// Generate the meme image
	imageData, mimeType, err := s.generateMemeImage(
		req.TemplateId,
		req.TopText,
		req.BottomText,
		req.AdditionalText,
	)

	if err != nil {
		log.Printf("Error generating meme: %v", err)
		return &pb.GenerateMemeResponse{
			Error: fmt.Sprintf("Failed to generate meme: %v", err),
		}, nil
	}

	return &pb.GenerateMemeResponse{
		ImageData:         imageData,
		MimeType:          mimeType,
		GeneratedCaptions: generatedCaptions,
		Error:             "",
	}, nil
}

// ListTemplates returns a list of available meme templates
func (s *MemeService) ListTemplates(ctx context.Context, req *pb.ListTemplatesRequest) (*pb.ListTemplatesResponse, error) {
	log.Printf("Service: Listing templates with category filter: %s", req.Category)

	var templates []*pb.Template

	for id, info := range s.Templates {
		// Apply category filter if specified
		if req.Category != "" && info.Category != req.Category {
			continue
		}

		templates = append(templates, &pb.Template{
			Id:             id,
			Name:           info.Name,
			TextFieldCount: info.TextFieldCount,
			Category:       info.Category,
			PreviewUrl:     fmt.Sprintf("/templates/%s", info.Filename),
		})
	}

	return &pb.ListTemplatesResponse{
		Templates: templates,
	}, nil
}

// generateMemeImage creates a meme image with the given template and text
func (s *MemeService) generateMemeImage(templateID, topText, bottomText string, additionalText []string) (string, string, error) {
	template, exists := s.Templates[templateID]
	if !exists {
		return "", "", fmt.Errorf("template '%s' not found", templateID)
	}

	// Load the template image
	imgPath := filepath.Join(s.Config.TemplateDir, template.Filename)
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open template image: %v", err)
	}
	defer imgFile.Close()

	// Decode the image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode template image: %v", err)
	}

	// Create a new RGBA image
	bounds := img.Bounds()
	memeImg := image.NewRGBA(bounds)

	// Draw the template onto the new image
	draw.Draw(memeImg, bounds, img, bounds.Min, draw.Src)

	// Load the font
	fontData, err := ioutil.ReadFile(s.Config.FontFile)
	if err != nil {
		// Fallback to default font if custom font can't be loaded
		// In a real implementation, you'd embed the font or handle this better
		return "", "", fmt.Errorf("failed to load font: %v", err)
	}

	f, err := freetype.ParseFont(fontData)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse font: %v", err)
	}

	// Set up the context for drawing text
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetClip(bounds)
	c.SetDst(memeImg)
	c.SetHinting(font.HintingFull)
	c.SetSrc(image.NewUniform(color.White))

	// Get image dimensions
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// Adjust font size based on image dimensions
	dynamicFontSize := float64(imgWidth) / 12
	if dynamicFontSize < 16 {
		dynamicFontSize = 16
	} else if dynamicFontSize > 48 {
		dynamicFontSize = 48
	}
	c.SetFontSize(dynamicFontSize)

	// Draw top text
	if topText != "" {
		s.drawTextWithStroke(c, topText, imgWidth/2, int(dynamicFontSize*1.5), imgWidth, f)
	}

	// Draw bottom text
	if bottomText != "" {
		s.drawTextWithStroke(c, bottomText, imgWidth/2, imgHeight-int(dynamicFontSize*1.5), imgWidth, f)
	}

	// Handle additional text for multi-panel memes
	for i, text := range additionalText {
		if i >= int(template.TextFieldCount)-2 {
			break // Only use as many text fields as the template supports
		}

		// Position additional text fields based on template and panel count
		// This is a simplified approach; real implementation would be template-specific
		x := imgWidth / 2
		y := imgHeight/2 + (i-1)*int(dynamicFontSize*2)

		s.drawTextWithStroke(c, text, x, y, imgWidth, f)
	}

	// Encode the image to JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, memeImg, &jpeg.Options{Quality: s.Config.ImageQuality})
	if err != nil {
		return "", "", fmt.Errorf("failed to encode image: %v", err)
	}

	// Return base64 encoded image
	return base64.StdEncoding.EncodeToString(buf.Bytes()), "image/jpeg", nil
}

// drawTextWithStroke draws text with a black outline
func (s *MemeService) drawTextWithStroke(c *freetype.Context, text string, x, y, maxWidth int, f *truetype.Font) {
	// Get the current font size from the context
	fontSize := s.Config.FontSize
	lineSpacing := s.Config.LineSpacing

	// Split the text into lines if it's too long
	lines := []string{text}
	if len(text)*int(c.PointToFixed(fontSize)>>6)/2 > maxWidth {
		// Split text into multiple lines
		words := strings.Split(text, " ")
		currentLine := ""
		lines = []string{}

		for _, word := range words {
			testLine := currentLine
			if testLine != "" {
				testLine += " "
			}
			testLine += word

			// Measure the width of the test line
			opts := truetype.Options{
				Size: fontSize,
			}
			face := truetype.NewFace(f, &opts)
			width := font.MeasureString(face, testLine)

			if width.Ceil() < maxWidth-20 {
				currentLine = testLine
			} else {
				if currentLine != "" {
					lines = append(lines, currentLine)
				}
				currentLine = word
			}
		}
		if currentLine != "" {
			lines = append(lines, currentLine)
		}
	}

	// Calculate y position for all lines
	lineHeight := int(c.PointToFixed(fontSize*lineSpacing) >> 6)
	startY := y - (len(lines)-1)*lineHeight/2

	// Draw each line
	for i, line := range lines {
		// Set up for measuring text width
		opts := truetype.Options{
			Size: fontSize,
		}
		face := truetype.NewFace(f, &opts)
		width := font.MeasureString(face, line)

		// Center text horizontally
		textX := x - width.Ceil()/2
		textY := startY + i*lineHeight

		// Draw text outline (stroke)
		strokeSize := int(fontSize / 6)
		if strokeSize < 1 {
			strokeSize = 1
		}

		// Save current color (need to recreate it since we can't access it directly)
		origColor := image.NewUniform(color.White)

		// Set color for outline
		blackColor := image.NewUniform(color.Black)
		c.SetSrc(blackColor)

		// Draw outline by shifting text position slightly in all directions
		for dy := -strokeSize; dy <= strokeSize; dy++ {
			for dx := -strokeSize; dx <= strokeSize; dx++ {
				if dx == 0 && dy == 0 {
					continue // Skip the center position
				}
				pt := freetype.Pt(textX+dx, textY+dy)
				_, err := c.DrawString(line, pt)
				if err != nil {
					continue
				}
			}
		}

		// Restore original color for main text
		c.SetSrc(origColor)

		// Draw the main text
		pt := freetype.Pt(textX, textY)
		_, err := c.DrawString(line, pt)
		if err != nil {
			continue
		}
	}
}
