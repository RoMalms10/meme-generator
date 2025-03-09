.PHONY: all build test test-unit test-functional clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=meme-generator

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test: test-unit test-functional

test-unit:
	$(GOTEST) -v ./tests -run "^Test[^Functional]"

test-functional:
	$(GOTEST) -v ./tests -run "^TestFunctional"

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./tests
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

deps:
	$(GOMOD) download

tidy:
	$(GOMOD) tidy

# Generate test data directory if it doesn't exist
testdata:
	mkdir -p testdata

# Create a sample test image - uses ImageMagick if available
test-image: testdata
	@if command -v convert >/dev/null 2>&1; then \
		convert -size 400x300 xc:white -fill black -pointsize 20 -gravity center \
		-draw "text 0,0 'Test Template'" ./testdata/test-template.jpg; \
		echo "Created test template image"; \
	else \
		echo "ImageMagick not found, skipping test image creation"; \
	fi

# Initialize project for development
init: deps testdata test-image
	@echo "Project initialized for development"