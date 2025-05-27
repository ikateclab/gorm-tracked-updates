# GORM Code Generation Makefile

.PHONY: help install generate clean demo test

# Default target
help:
	@echo "GORM Code Generation Commands:"
	@echo ""
	@echo "  install   - Install the gorm-gen tool"
	@echo "  generate  - Run go generate to create clone and diff methods"
	@echo "  clean     - Remove generated files"
	@echo "  demo      - Run the demo showing generated methods"
	@echo "  test      - Test the generated code"
	@echo ""
	@echo "Usage:"
	@echo "  make install && make generate && make demo"

# Install the gorm-gen tool
install:
	@echo "🔧 Installing gorm-gen tool..."
	@cd ../../ && go install ./cmd/gorm-gen
	@echo "✅ gorm-gen installed successfully"

# Run go generate
generate:
	@echo "🚀 Running go generate..."
	@go generate ./models/...
	@echo "✅ Code generation completed"

# Clean generated files
clean:
	@echo "🧹 Cleaning generated files..."
	@rm -f models/clone.go models/diff.go
	@echo "✅ Generated files removed"

# Run the demo
demo: generate
	@echo "🎯 Running demo..."
	@go run demo.go

# Test the generated code
test: generate
	@echo "🧪 Testing generated code..."
	@go build ./models/...
	@echo "✅ Generated code compiles successfully"

# Full workflow
all: install generate demo
	@echo "🎉 Full workflow completed successfully!"
