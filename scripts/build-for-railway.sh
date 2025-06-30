#!/bin/bash

echo "🚀 Building for Railway deployment..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Make sure you're in the project root."
    exit 1
fi

# Check if Dockerfile exists
if [ ! -f "Dockerfile" ]; then
    echo "❌ Error: Dockerfile not found."
    exit 1
fi

# Validate Go version
echo "📋 Checking Go version..."
go version

# Test the build locally
echo "🔨 Testing build..."
if ! go build -o /tmp/test-build ./src/server/; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build test successful!"

# Clean up test build
rm -f /tmp/test-build

echo "🎉 Ready for Railway deployment!"
echo ""
echo "📝 Next steps:"
echo "1. Push your code to GitHub"
echo "2. Connect your repository to Railway"
echo "3. Configure environment variables in Railway json file"
echo "4. Deploy!"
