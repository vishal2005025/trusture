#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}=== TrustuRe NGO Blockchain Platform ===${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Function to run with error checking
run_cmd() {
    echo -e "${YELLOW}Running: $1${NC}"
    if eval "$1"; then
        echo -e "${GREEN}✓ Success${NC}"
        echo
    else
        echo -e "${RED}✗ Failed${NC}"
        exit 1
    fi
}

# Ensure we're in the project root
cd "$(dirname "$0")"

echo "Checking Go version..."
go version

echo -e "\n${YELLOW}=== Installing Dependencies ===${NC}"
run_cmd "go mod tidy"
run_cmd "go mod download"

echo -e "${YELLOW}=== Running Tests ===${NC}"
run_cmd "go test -v ./pkg/blockchain/..."
run_cmd "go test -v ./pkg/..."

echo -e "${YELLOW}=== Building Backend ===${NC}"
run_cmd "go build -v -o bin/trusture ./cmd/api"

echo -e "${YELLOW}=== Setting Up Frontend ===${NC}"
if [ -d "Frontend" ]; then
    cd Frontend
    if ! command -v npm &> /dev/null; then
        echo -e "${RED}Error: npm is not installed${NC}"
        echo "Please install Node.js and npm from https://nodejs.org/"
        exit 1
    fi
    
    echo "Installing frontend dependencies..."
    run_cmd "npm install"
    
    echo "Building frontend..."
    run_cmd "npm run build"
    cd ..
fi

echo -e "${GREEN}=== Setup Complete! ===${NC}"
echo -e "To run the application:"
echo -e "1. Start the backend:"
echo -e "   ${YELLOW}./bin/trusture${NC}"
echo -e "2. Start the frontend (in another terminal):"
echo -e "   ${YELLOW}cd Frontend && npm run dev${NC}"
echo
echo -e "API documentation will be available at: http://localhost:8080/swagger/index.html"
echo -e "Frontend will be available at: http://localhost:3000"