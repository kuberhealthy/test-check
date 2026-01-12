IMAGE := "kuberhealthy/test-check"
TAG := "latest"

# Build the test check container locally.
build:
	podman build -f Containerfile -t {{IMAGE}}:{{TAG}} .

# Run the unit tests for the test check.
test:
	go test ./...

# Build the test check binary locally.
binary:
	go build -o bin/test-check ./cmd/test-check
