build:
	@echo "Building binary..."
	go build .

clean:
	@echo "Cleaning up..."
	rm grep
	go clean

test:
	@echo "Running tests..."
	go test -v