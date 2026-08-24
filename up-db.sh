#!/usr/bin/env bash
set -euo pipefail

CONTAINER_NAME="freek-test-db"
PORT=5432
DB_NAME="freek_test"
IMAGE="postgres:18"

if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "Container ${CONTAINER_NAME} is already running"
elif docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "Starting existing container ${CONTAINER_NAME}"
    docker start "${CONTAINER_NAME}"
else
    echo "Creating container ${CONTAINER_NAME} (${IMAGE})"
    docker run -d \
        --name "${CONTAINER_NAME}" \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_PASSWORD=postgres \
        -e POSTGRES_DB="${DB_NAME}" \
        -p "${PORT}:5432" \
        "${IMAGE}"
fi

echo "Waiting for PostgreSQL to be ready..."
until docker exec "${CONTAINER_NAME}" pg_isready -U postgres -d "${DB_NAME}" >/dev/null 2>&1; do
    sleep 0.5
done

echo "Applying migrations..."
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
for file in "${SCRIPT_DIR}"/migrations/*.up.sql; do
    echo "  $(basename "$file")"
    docker exec -i "${CONTAINER_NAME}" psql -U postgres -d "${DB_NAME}" -v ON_ERROR_STOP=1 < "$file" > /dev/null
done

echo "PostgreSQL 18 is up:"
echo "  URL: postgres://localhost:${PORT}/${DB_NAME}?sslmode=disable"
echo "Run tests with: go test ./..."
