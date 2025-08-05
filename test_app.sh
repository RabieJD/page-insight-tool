#!/bin/bash

echo "Testing Page Insight Tool..."

# Build the application
echo "Building application..."
make build

# Start the server in background
echo "Starting server..."
./app/.bin/page-insight-tool --debug &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Test the server
echo "Testing server..."
if curl -s http://localhost:8080 > /dev/null; then
    echo "✅ Server is running and responding"
else
    echo "❌ Server is not responding"
fi

# Kill the server
kill $SERVER_PID 2>/dev/null

echo "Test completed" 