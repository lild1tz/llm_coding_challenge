#!/bin/sh

echo "Running migrations..."
migrate -path /app/migrations -database "$DATABASE_URL" up

echo "Starting application..."
exec /app/hermes