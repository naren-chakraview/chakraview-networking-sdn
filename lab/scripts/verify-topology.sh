#!/bin/bash

CONTROLLER_URL="${CONTROLLER_URL:-http://localhost:8080}"

echo "Topology Status:"
curl -s -X GET "$CONTROLLER_URL/api/v1/topology" | python -m json.tool

echo ""
echo "Registered Devices:"
curl -s -X GET "$CONTROLLER_URL/api/v1/topology/devices" | python -m json.tool
