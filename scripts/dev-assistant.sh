#!/bin/bash

# Promptext Dev Assistant
# Interactive tool to check if staged changes need documentation updates

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🤖 Promptext Dev Assistant${NC}"
echo ""

# Check if there are staged changes
if ! git diff --cached --quiet; then
    # Run the Go assistant
    go run cmd/dev-assistant/main.go "$@"
else
    echo "⚠️  No staged changes detected."
    echo ""
    echo "Stage some changes first:"
    echo "  git add <files>"
    echo ""
    exit 0
fi
