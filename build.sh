#!/bin/bash
# Quick build and install script for Unix systems

set -e

echo "Building devgitsecops..."
go build -v -o bin/devgitsecops main.go

echo ""
echo "Build successful! Binary created at: bin/devgitsecops"
echo ""
echo "Quick start:"
echo "  1. Check tool status: ./bin/devgitsecops status"
echo "  2. Install tools:     ./bin/devgitsecops install --all --auto"
echo "  3. Use tools:         ./bin/devgitsecops kubectl get pods"
echo ""
echo "To add to PATH, add this line to your ~/.bashrc or ~/.zshrc:"
echo "  export PATH=\"\$PATH:$(pwd)/bin\""
echo ""
echo "Or install system-wide:"
echo "  sudo cp bin/devgitsecops /usr/local/bin/"
