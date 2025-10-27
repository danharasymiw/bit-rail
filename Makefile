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


