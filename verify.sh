#!/bin/bash

echo "=== TrustuRe NGO Blockchain Platform Verification ==="
echo "Running complete test suite..."

# Function to check command success
check_step() {
    if [ $? -eq 0 ]; then
        echo "✓ $1"
    else
        echo "✗ $1 failed"
        exit 1
    fi
}

# 1. Run all tests with coverage
echo "\nRunning tests with coverage..."
go test ./... -coverprofile=coverage.out
check_step "All tests"

# 2. Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
check_step "Coverage report generated"

# 3. Run blockchain specific tests
echo "\nRunning blockchain package tests..."
cd pkg/blockchain
go test -v
check_step "Blockchain tests"

# 4. Run main application with test data
echo "\nTesting main application..."
cd ../../
go run cmd/main.go &
APP_PID=$!
sleep 5
kill $APP_PID
check_step "Application startup"

# 5. Print results
echo "\n=== Test Results ==="
echo "- Coverage report available in: coverage.html"
echo "- All test cases executed"
echo "- Blockchain functionality verified"
echo "- Application startup confirmed"

# Cleanup
rm -f coverage.out

echo "\nVerification complete! ✨"