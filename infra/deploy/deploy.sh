#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

APP_ROOT="${APP_ROOT:-/opt/dth}"
RELEASE_DIR="${APP_ROOT}/current"
BIN_DIR="${RELEASE_DIR}/bin"

run() {
  echo "+ $*"
  "$@"
}

require() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command not found: $1" >&2
    exit 1
  fi
}

require go
require npm
require rsync
require systemctl
require nginx

echo "Preparing release directories"
run mkdir -p "${BIN_DIR}"
run mkdir -p "${RELEASE_DIR}/backend"
run mkdir -p "${RELEASE_DIR}/frontend/customer"
run mkdir -p "${RELEASE_DIR}/frontend/admin"

echo "Building backend binaries"
run bash -lc "cd '${REPO_ROOT}/backend' && go build -o '${BIN_DIR}/dth-api' ./cmd/api"
run bash -lc "cd '${REPO_ROOT}/backend' && go build -o '${BIN_DIR}/dth-worker' ./cmd/worker"

echo "Syncing backend source"
run rsync -a --delete "${REPO_ROOT}/backend/" "${RELEASE_DIR}/backend/"

echo "Building customer frontend"
run bash -lc "cd '${REPO_ROOT}/frontend/customer' && npm ci && npm run build"
run rsync -a --delete "${REPO_ROOT}/frontend/customer/out/" "${RELEASE_DIR}/frontend/customer/out/"

echo "Building admin frontend"
run bash -lc "cd '${REPO_ROOT}/frontend/admin' && npm ci && npm run build"
run rsync -a --delete "${REPO_ROOT}/frontend/admin/out/" "${RELEASE_DIR}/frontend/admin/out/"

echo "Installing service and nginx config"
run sudo install -m 0644 "${REPO_ROOT}/infra/systemd/dth-api.service" /etc/systemd/system/dth-api.service
run sudo install -m 0644 "${REPO_ROOT}/infra/systemd/dth-worker.service" /etc/systemd/system/dth-worker.service
run sudo install -m 0644 "${REPO_ROOT}/infra/nginx/dth.conf" /etc/nginx/conf.d/dth.conf

echo "Reloading systemd and nginx"
run sudo systemctl daemon-reload
run sudo systemctl enable dth-api dth-worker
run sudo systemctl restart dth-api dth-worker
run sudo nginx -t
run sudo systemctl reload nginx

echo "Deployment complete"
