#!/bin/bash
set -e

CONTROLLER_URL="${CONTROLLER_URL:-http://localhost:8080}"
GRPC_URL="${GRPC_URL:-localhost:9090}"

echo "Initializing SDN Lab..."
echo "Controller URL: $CONTROLLER_URL"
echo "gRPC URL: $GRPC_URL"

# Wait for controller to be ready
echo "Waiting for controller to be ready..."
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
    if curl -s -f "$CONTROLLER_URL/api/v1/health" > /dev/null; then
        echo "Controller is ready!"
        break
    fi
    attempt=$((attempt + 1))
    sleep 1
done

if [ $attempt -eq $max_attempts ]; then
    echo "ERROR: Controller did not become ready"
    exit 1
fi

# Register devices
echo "Registering network devices..."
bash /app/scripts/register-devices.sh

# Verify topology
echo "Verifying topology..."
bash /app/scripts/verify-topology.sh

echo "Lab initialization complete!"
