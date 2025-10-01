
OUTPUT=ldapPubKeyReader
# Build flags for optimization and security
LDFLAGS=-ldflags '-w -s -extldflags "-static"'
BUILDFLAGS=-a -installsuffix cgo -trimpath
CGO_ENV=CGO_ENABLED=0

all: build-all

# Build for all platforms
build-all: linux-amd64 darwin-amd64 darwin-arm64

# Linux amd64
linux-amd64:
	$(CGO_ENV) GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(OUTPUT)-linux-amd64

# macOS amd64 (Intel)  
darwin-amd64:
	$(CGO_ENV) GOOS=darwin GOARCH=amd64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(OUTPUT)-darwin-amd64

# macOS arm64 (Apple Silicon)
darwin-arm64:
	$(CGO_ENV) GOOS=darwin GOARCH=arm64 go build $(BUILDFLAGS) $(LDFLAGS) -o $(OUTPUT)-darwin-arm64

# Development builds (with debug info and race detection)
dev-all: dev-linux-amd64 dev-darwin-amd64 dev-darwin-arm64

dev-linux-amd64:
	$(CGO_ENV) GOOS=linux GOARCH=amd64 go build -race -o $(OUTPUT)-dev-linux-amd64

dev-darwin-amd64:
	$(CGO_ENV) GOOS=darwin GOARCH=amd64 go build -race -o $(OUTPUT)-dev-darwin-amd64

dev-darwin-arm64:
	$(CGO_ENV) GOOS=darwin GOARCH=arm64 go build -race -o $(OUTPUT)-dev-darwin-arm64

# Legacy build target (defaults to linux amd64)
build: linux-amd64

# Legacy static target (now same as regular build since CGO is disabled)
static: linux-amd64

# Clean built binaries
clean:
	rm -f $(OUTPUT)-*

# Show build information
info:
	@echo "Build Configuration:"
	@echo "  OUTPUT: $(OUTPUT)"
	@echo "  LDFLAGS: $(LDFLAGS)"
	@echo "  BUILDFLAGS: $(BUILDFLAGS)"
	@echo "  CGO_ENV: $(CGO_ENV)"
	@echo ""
	@echo "Available targets:"
	@echo "  build-all    - Build optimized binaries for all platforms"
	@echo "  dev-all      - Build development binaries with race detection"
	@echo "  linux-amd64  - Build for Linux x86_64"
	@echo "  darwin-amd64 - Build for macOS Intel"
	@echo "  darwin-arm64 - Build for macOS Apple Silicon"
	@echo "  clean        - Remove all built binaries"
	@echo "  info         - Show this information"

.PHONY: all build-all linux-amd64 darwin-amd64 darwin-arm64 dev-all dev-linux-amd64 dev-darwin-amd64 dev-darwin-arm64 build static clean info

