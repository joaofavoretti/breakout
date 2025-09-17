.PHONY: build run clean test

# Build the game
build:
	go build -o breakout .

# Run the game
run: build
	./breakout

# Clean build artifacts
clean:
	rm -f breakout

# Run tests (when we add them)
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run all checks
check: fmt vet test