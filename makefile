build: main.go
	@echo "Building binary..."
	go build .

clean:
	@echo "Cleaning up..."
	rm grep
	go clean