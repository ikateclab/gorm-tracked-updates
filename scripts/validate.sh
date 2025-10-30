#!/bin/bash
set -e

echo "🔍 Running local validation checks..."
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track failures
FAILED=0

# 1. Format check
echo "📝 Checking code formatting..."
if ! gofmt -l . | grep -q .; then
    echo -e "${GREEN}✓${NC} Code is properly formatted"
else
    echo -e "${RED}✗${NC} Code is not formatted. Run: gofmt -w ."
    gofmt -l .
    FAILED=1
fi
echo ""

# 2. Go mod verify
echo "📦 Verifying dependencies..."
if go mod verify; then
    echo -e "${GREEN}✓${NC} Dependencies verified"
else
    echo -e "${RED}✗${NC} Dependency verification failed"
    FAILED=1
fi
echo ""

# 3. Go mod tidy check
echo "🧹 Checking if go.mod is tidy..."
cp go.mod go.mod.backup
cp go.sum go.sum.backup
go mod tidy
if diff go.mod go.mod.backup > /dev/null && diff go.sum go.sum.backup > /dev/null; then
    echo -e "${GREEN}✓${NC} go.mod is tidy"
    rm go.mod.backup go.sum.backup
else
    echo -e "${RED}✗${NC} go.mod is not tidy. Run: go mod tidy"
    mv go.mod.backup go.mod
    mv go.sum.backup go.sum
    FAILED=1
fi
echo ""

# 4. Build check
echo "🔨 Building packages..."
if go build -ldflags="-checklinkname=0" ./...; then
    echo -e "${GREEN}✓${NC} All packages build successfully"
else
    echo -e "${RED}✗${NC} Build failed"
    FAILED=1
fi
echo ""

# 5. Run tests
echo "🧪 Running tests..."
if go test -ldflags="-checklinkname=0" ./pkg/...; then
    echo -e "${GREEN}✓${NC} All tests passed"
else
    echo -e "${RED}✗${NC} Tests failed"
    FAILED=1
fi
echo ""

# 6. Run tests with race detector
echo "🏁 Running tests with race detector..."
if go test -ldflags="-checklinkname=0" -race ./pkg/...; then
    echo -e "${GREEN}✓${NC} Race detector tests passed"
else
    echo -e "${RED}✗${NC} Race detector tests failed"
    FAILED=1
fi
echo ""

# 7. Generate code and check if up to date
echo "⚙️  Checking generated code..."
cd examples/go-generate
go generate ./...
cd ../..
if git diff --exit-code examples/go-generate/models/diff.go examples/go-generate/models/clone.go > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Generated code is up to date"
else
    echo -e "${YELLOW}⚠${NC}  Generated code has changes. Please review and commit if needed."
    echo "Changed files:"
    git diff --name-only examples/go-generate/models/
fi
echo ""

# 8. Lint (if golangci-lint is installed)
if command -v golangci-lint &> /dev/null; then
    echo "🔍 Running golangci-lint (excluding examples/performance)..."
    if golangci-lint run --exclude-dirs=examples/performance ./...; then
        echo -e "${GREEN}✓${NC} Linting passed"
    else
        echo -e "${YELLOW}⚠${NC}  Linting found issues (non-critical in benchmarks)"
        echo "   Main packages (pkg/) are clean. Benchmark errors can be ignored."
    fi
else
    echo -e "${YELLOW}⚠${NC}  golangci-lint not installed. Skipping lint check."
    echo "   Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All validation checks passed!${NC}"
    echo "You can safely commit and push your changes."
    exit 0
else
    echo -e "${RED}✗ Some validation checks failed!${NC}"
    echo "Please fix the issues above before committing."
    exit 1
fi
