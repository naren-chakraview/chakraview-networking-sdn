#!/bin/bash

CONTROLLER_URL="${CONTROLLER_URL:-http://localhost:8080}"

# Register leaf switches
for i in 1 2; do
    curl -X POST "$CONTROLLER_URL/api/v1/devices/register" \
        -H "Content-Type: application/json" \
        -d "{\"device_id\":\"leaf$i\",\"address\":\"10.0.$i.1\",\"role\":\"leaf\"}" \
        2>/dev/null || true
done

# Register spine switches
for i in 1 2; do
    curl -X POST "$CONTROLLER_URL/api/v1/devices/register" \
        -H "Content-Type: application/json" \
        -d "{\"device_id\":\"spine$i\",\"address\":\"10.1.$i.1\",\"role\":\"spine\"}" \
        2>/dev/null || true
done

echo "Device registration complete"
