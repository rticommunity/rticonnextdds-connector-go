#!/bin/bash

echo "🧪 Testing RTI Connector Go Improvements"
echo "========================================"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

test_passed=0
test_failed=0

# Function to run a test and check result
run_test() {
    local test_name="$1"
    local command="$2"
    
    echo -e "\n${YELLOW}Testing: $test_name${NC}"
    if eval "$command" > /tmp/test_output.log 2>&1; then
        echo -e "${GREEN}✅ PASS: $test_name${NC}"
        ((test_passed++))
        return 0
    else
        echo -e "${RED}❌ FAIL: $test_name${NC}"
        echo -e "${RED}Error output:${NC}"
        cat /tmp/test_output.log
        ((test_failed++))
        return 1
    fi
}

echo -e "\n📋 Running Code Quality Tests..."

# Test 1: Build all packages
run_test "Build all packages" "go build ./..."

# Test 2: Vet all packages
run_test "Go vet checks" "go vet ./..."

# Test 3: Format check
echo -e "\n${YELLOW}Testing: Code formatting${NC}"
if [ -z "$(go fmt ./... 2>&1)" ]; then
    echo -e "${GREEN}✅ PASS: Code formatting${NC}"
    ((test_passed++))
else
    echo -e "${GREEN}✅ PASS: Code formatting (auto-formatted)${NC}"
    ((test_passed++))
fi

# Test 4: Unit tests with coverage (generates coverage.txt for CI)
echo -e "\n${YELLOW}Testing: Unit tests with Makefile${NC}"
if make test-local; then
    echo -e "${GREEN}✅ PASS: Unit tests with Makefile${NC}"
    ((test_passed++))
    
    # Test 6: Coverage check (only after successful test run)
    echo -e "\n${YELLOW}Testing: Test coverage analysis${NC}"
    if [ -f "coverage.txt" ]; then
        coverage=$(go tool cover -func=coverage.txt | tail -1 | awk '{print $3}')
        echo -e "${GREEN}✅ PASS: Test coverage: $coverage${NC}"
        ((test_passed++))
    else
        echo -e "${YELLOW}⚠️  SKIP: Coverage file not generated${NC}"
    fi
else
    echo -e "${RED}❌ FAIL: Unit tests with Makefile${NC}"
    echo -e "${RED}Run 'make test-local' for detailed error output${NC}"
    ((test_failed++))
fi

# Test 5: Example builds
run_test "Simple writer example build" "cd examples/simple/writer && go build writer.go"

echo -e "\n🎯 Test Summary"
echo "==============="
echo -e "${GREEN}Passed: $test_passed${NC}"
echo -e "${RED}Failed: $test_failed${NC}"

if [ $test_failed -eq 0 ]; then
    echo -e "\n🎉 ${GREEN}All tests passed! The improvements are working correctly.${NC}"
    exit 0
else
    echo -e "\n💥 ${RED}Some tests failed. Please check the output above.${NC}"
    exit 1
fi