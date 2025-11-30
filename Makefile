# Install dependencies
.PHONY: dependencies
dependencies:
	@echo "Installing dependencies..."
	@go install github.com/a-h/templ/cmd/templ@latest
	@make dependency-tailwind
	@echo "Dependencies installed successfully"

.PHONY: dev-dependencies
dev-dependencies: dependencies
	@echo "Installing development dependencies..."
	@go install github.com/air-verse/air@latest
	@echo "Development dependencies installed successfully"

# Install Tailwind CSS based on machine architecture and OS
.PHONY: dependency-tailwind
dependency-tailwind:
	@echo "Installing Tailwind CSS..."
	@mkdir -p ./bin
	@OS=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	ARCH=$$(uname -m); \
	if [ "$$ARCH" = "x86_64" ]; then \
		ARCH="x64"; \
	elif [ "$$ARCH" = "aarch64" ] || [ "$$ARCH" = "arm64" ]; then \
		ARCH="arm64"; \
	fi; \
	if [ "$$OS" = "darwin" ]; then \
		OS="macos"; \
	elif [ "$$OS" = "linux" ]; then \
		OS="linux"; \
	elif [ "$$OS" = "windows_nt" ] || [[ "$$OS" == "msys"* ]] || [[ "$$OS" == "mingw"* ]]; then \
		OS="windows"; \
		EXT=".exe"; \
	fi; \
	echo "Detected OS: $$OS, Architecture: $$ARCH"; \
	URL="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$$OS-$$ARCH$$EXT"; \
	echo "Downloading Tailwind CSS from: $$URL"; \
	curl -sLO "$$URL"; \
	chmod +x "tailwindcss-$$OS-$$ARCH$$EXT"; \
	mv "tailwindcss-$$OS-$$ARCH$$EXT" ./bin/tailwindcss$$EXT; \
	echo "Tailwind CSS installed successfully at ./bin/tailwindcss$$EXT"

.PHONY: tailwind-clean
tailwind-clean:
	@./bin/tailwindcss -i ./internal/frontend/static/css/input.css -o ./internal/frontend/static/css/output.css --clean

# Run the application with hot reload using air
.PHONY: air
air:
	air

# Run the application with hot reload using air
.PHONY: dev
dev: tailwind-clean
	make -j3 templ-watch air tailwind-watch

.PHONY: tailwind-watch
tailwind-watch:
	@./bin/tailwindcss -i ./internal/frontend/static/css/input.css -o ./internal/frontend/static/css/output.css --watch

.PHONY: tailwind-build
tailwind-build:
	@./bin/tailwindcss -i ./internal/frontend/static/css/input.css -o ./internal/frontend/static/css/output.css

.PHONY: templ-watch
templ-watch:
	@templ generate --watch

.PHONY: templ-generate
templ-generate:
	@templ generate
	
# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t raspberry-bookshelf:latest .

.PHONY: docker-run
docker-run:
	@echo "Running Docker image..."
	@docker run --rm -it --name raspberry-bookshelf -p 8080:8080 raspberry-bookshelf:latest 

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ./bin

.PHONY: build
build:
	@echo "Building application..."
	go build -o ./bin/app ./cmd/main.go

.PHONY: build-all
build-all:
	@echo "Tailwind CSS build..."
	@make tailwind-build
	@echo "Compiling templ templates..."
	@make templ-generate
	@echo "Building application on all platforms..."
	GOOS=windows GOARCH=amd64 go build -o ./bin/bookshelf_win_x64.exe ./cmd/main.go
	GOOS=windows GOARCH=386 go build -o ./bin/bookshelf_win_x86.exe ./cmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./bin/bookshelf_mac_intel ./cmd/main.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/bookshelf_mac_silicon ./cmd/main.go
	GOOS=linux GOARCH=386 go build -o ./bin/bookshelf_linux_386 ./cmd/main.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/bookshelf_linux_amd64 ./cmd/main.go
	GOOS=linux GOARCH=arm64 go build -o ./bin/bookshelf_linux_arm64 ./cmd/main.go

