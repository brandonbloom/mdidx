#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

TESTS_RUN=0
TESTS_FAILED=0

test_passed() {
    echo -e "${GREEN}✓ $1${NC}"
}

test_failed() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

run_test() {
    ((TESTS_RUN++))
    echo -e "\n${YELLOW}Test $TESTS_RUN: $1${NC}"
}

cleanup() {
    rm -f test_*.md test_*.mdidx temp_output.mdidx
}

trap cleanup EXIT

echo "Building mdidx..."
go build -o mdidx

echo "Running integration tests..."

# Test 1: Help flag
run_test "Help flag works"
if ./mdidx --help | grep -q "mdidx - Build an index from a markdown file"; then
    test_passed "Help message displayed correctly"
else
    test_failed "Help message not displayed"
fi

# Test 2: Simple file processing
run_test "Simple file processing"
cat > test_simple.md << 'EOF'
# Title

Some content here

## Subtitle

Some other content here.

# Footer

The end
EOF

./mdidx test_simple.md

if [[ -f test_simple.mdidx ]]; then
    test_passed "Output file created"
    
    if grep -q "source: test_simple.md" test_simple.mdidx; then
        test_passed "Source field present"
    else
        test_failed "Source field missing"
    fi
    
    if grep -q "1-.*: Title" test_simple.mdidx; then
        test_passed "Title header found"
    else
        test_failed "Title header missing"
    fi
    
    if grep -q "  5-.*: Subtitle" test_simple.mdidx; then
        test_passed "Subtitle indentation correct"
    else
        test_failed "Subtitle indentation incorrect"
    fi
else
    test_failed "Output file test_simple.mdidx not created"
fi

# Test 3: Custom output file
run_test "Custom output file (-o flag)"
./mdidx -o temp_output.mdidx test_simple.md

if [[ -f temp_output.mdidx ]]; then
    test_passed "Custom output file created"
else
    test_failed "Custom output file not created"
fi

# Test 4: Stdin/stdout processing
run_test "Stdin/stdout processing"
output=$(echo -e "# Test\n\nContent\n\n## Sub\n\nMore content" | ./mdidx)

if echo "$output" | grep -q "1-.*: Test"; then
    test_passed "Test header found in stdin output"
else
    test_failed "Stdin processing failed"
fi

if echo "$output" | grep -q "  5-.*: Sub"; then
    test_passed "Sub header indentation correct in stdin output"
else
    test_failed "Sub header indentation incorrect in stdin output"
fi

if ! echo "$output" | grep -q "source:"; then
    test_passed "No source field in stdin output (correct)"
else
    test_failed "Source field present in stdin output (incorrect)"
fi

# Test 5: Error handling - nonexistent file
run_test "Error handling for nonexistent file"
if ./mdidx nonexistent.md 2>/dev/null; then
    test_failed "Should have failed for nonexistent file"
else
    test_passed "Correctly handles nonexistent file"
fi

# Test 6: Error handling - too many arguments
run_test "Error handling for too many arguments"
if ./mdidx file1.md file2.md 2>/dev/null; then
    test_failed "Should have failed for too many arguments"
else
    test_passed "Correctly handles too many arguments"
fi

# Summary
echo -e "\n${YELLOW}=== Test Summary ===${NC}"
echo "Tests run: $TESTS_RUN"
echo "Test failures: $TESTS_FAILED"

if [[ $TESTS_FAILED -eq 0 ]]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}$TESTS_FAILED test(s) failed.${NC}"
    exit 1
fi