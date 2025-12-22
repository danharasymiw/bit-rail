.PHONY: help clean-generated generate build run test

help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

clean:
	@echo "Removing all *_gen.go files"
	@find . -name "*_gen.go" -type f -delete
	@echo "Done! All generated files removed."

generate: clean
	@echo "Running go generate"
	@go generate ./...
	@echo "Done!"

local:
	go run cmd/bit-rail/main.go -local

local-debug:
	go run cmd/bit-rail/main.go -local -debug

server:
	 go run cmd/bit-rail/main.go -server

server-debug:
	go run cmd/bit-rail/main.go -server -debug

client:
	go run cmd/bit-rail/main.go

client-debug:
	go run cmd/bit-rail/main.go -debug

