#!/bin/bash

# Script to update module path for release
# Usage: ./scripts/update-module-path.sh YOUR_GITHUB_USERNAME

if [ $# -eq 0 ]; then
    echo "Usage: $0 YOUR_GITHUB_USERNAME"
    echo "Example: $0 johndoe"
    exit 1
fi

USERNAME=$1
OLD_PATH="github.com/yourusername/gorm-tracked-updates"
NEW_PATH="github.com/$USERNAME/gorm-tracked-updates"

echo "🔄 Updating module path from 'yourusername' to '$USERNAME'..."

# Update go.mod
sed -i.bak "s|yourusername|$USERNAME|g" go.mod

# Update all Go files
find . -name "*.go" -type f -exec sed -i.bak "s|yourusername|$USERNAME|g" {} \;

# Update documentation
find . -name "*.md" -type f -exec sed -i.bak "s|yourusername|$USERNAME|g" {} \;

# Clean up backup files
find . -name "*.bak" -type f -delete

echo "✅ Module path updated successfully!"
echo "📝 Updated files:"
echo "   - go.mod"
echo "   - All *.go files"
echo "   - All *.md files"
echo ""
echo "🚀 Next steps:"
echo "   1. Review the changes: git diff"
echo "   2. Test the build: go build ./..."
echo "   3. Commit changes: git add . && git commit -m 'Update module path'"
echo "   4. Create GitHub repository: $NEW_PATH"
echo "   5. Push to GitHub: git push origin main"
echo "   6. Create release tag: git tag v1.0.0 && git push origin v1.0.0"
