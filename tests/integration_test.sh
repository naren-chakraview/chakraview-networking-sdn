#!/bin/bash

set -e

echo "Running End-to-End Integration Tests..."

CONTROLLER_URL="http://localhost:8080"

# Test 1: Health check
echo "Test 1: Controller health check"
curl -f "$CONTROLLER_URL/api/v1/health" > /dev/null && echo "  PASS" || echo "  FAIL"

# Test 2: Topology endpoint
echo "Test 2: Topology endpoint"
curl -f "$CONTROLLER_URL/api/v1/topology" > /dev/null && echo "  PASS" || echo "  FAIL"

# Test 3: Devices endpoint
echo "Test 3: Devices endpoint"
curl -f "$CONTROLLER_URL/api/v1/topology/devices" > /dev/null && echo "  PASS" || echo "  FAIL"

echo "Integration tests complete"
