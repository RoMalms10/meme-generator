# Meme Generator Microservice

A gRPC-based microservice for generating custom memes with text overlays. This service is designed to run as part of a distributed system and communicates with other services using Protocol Buffers.

## Features

- Generate memes with customizable text
- Support for multiple meme templates
- Optional AI caption generation
- High-quality text rendering with outline/stroke
- Clean architecture with separation of concerns
- Comprehensive test suite
- Configurable via environment variables

## Architecture

This microservice is part of a larger system:

```
┌─────────────┐     gRPC     ┌─────────────┐
│   Edge API  │◄────────────►│    Meme     │
│             │     (grpc)   │  Generator  │
└─────────────┘              └─────────────┘
```

### Project Structure

```
├── config/             # Configuration management
│   └── env.go          # Environment variable handling
├── handler/            # Request handlers (gRPC interface)
│   └── handler.go      # Implementation of gRPC endpoints
├── server/             # Server setup and lifecycle
│   ├── router.go       # gRPC service registration
│   └── server.go       # Server lifecycle management
├── service/            # Business logic
│   └── meme_service.go # Meme generation logic
├── templates/          # Meme template images
│   └── *.jpg           # Template images
├── tests/              # Test suite
│   ├── config_test.go  # Config tests
│   ├── handler_test.go # Handler tests
│   ├── service_test.go # Service tests
│   └── ...             # Other tests
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── main.go             # Application entry point
└── Makefile            # Build and test commands
```

## Getting Started

### Prerequisites

- Go 1.19+
- gRPC tools (`protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`)
- Font files (Impact or similar for meme text)
- Templates directory with meme template images

### Installation

1. Clone the repository:
```bash
git clone https://github.com/RoMalms10/meme-generator.git
cd meme-generator
```

2. Install dependencies:
```bash
go mod download
```

3. Set up directories and assets:
```bash
mkdir -p templates fonts
# Add your template images to templates/
# Add your font files to fonts/
```

### Running the Service

```bash
go run main.go
```

The service will start on port 50051 by default. You can customize this and other settings using environment variables.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `50051` |
| `TEMPLATE_DIR` | Directory containing template images | `./templates` |
| `FONT_FILE` | Path to font file for text rendering | `./fonts/impact.ttf` |
| `IMAGE_QUALITY` | JPEG quality (1-100) | `90` |
| `FONT_SIZE` | Base font size for text | `36` |
| `LINE_SPACING` | Line spacing multiplier | `1.5` |
| `ENABLE_AI_CAPTION` | Enable AI caption generation | `false` |
| `GRPC_MAX_CONNECTION_AGE` | Max age of gRPC connections | `30m` |
| `GRPC_MAX_CONNECTION_IDLE` | Max idle time for connections | `15m` |
| `GRPC_KEEPALIVE_TIME` | Keepalive ping time | `5m` |
| `GRPC_KEEPALIVE_TIMEOUT` | Keepalive ping timeout | `20s` |

## API Reference

### gRPC Endpoints

#### GenerateMeme

Generates a meme with the specified template and text.

```protobuf
rpc GenerateMeme(GenerateMemeRequest) returns (GenerateMemeResponse);
```

#### ListTemplates

Lists available meme templates, optionally filtered by category.

```protobuf
rpc ListTemplates(ListTemplatesRequest) returns (ListTemplatesResponse);
```

## Development

### Running Tests

Run all tests:
```bash
make test
```

Run only unit tests:
```bash
make test-unit
```

Run only functional tests:
```bash
make test-functional
```

Generate test coverage:
```bash
make test-coverage
```

### Adding New Templates

1. Add the template image to the `templates` directory
2. Update the templates map in `service/meme_service.go` with the new template details

## Deployment

### Docker

Build the Docker image:
```bash
docker build -t meme-generator:latest .
```

Run with Docker:
```bash
docker run -p 50051:50051 meme-generator:latest
```

### Kubernetes

Apply the Kubernetes manifests:
```bash
kubectl apply -f kubernetes/manifests.yaml
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.