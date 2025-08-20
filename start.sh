#!/bin/bash

# Start the backend in the background
echo "Starting backend..."
(cd backend && go run .) &

# Start the frontend in the background
echo "Starting frontend..."
(cd frontend && npm install && npm run dev) &

# Wait for both processes to finish
wait
