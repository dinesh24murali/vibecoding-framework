#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

load_env_file() {
  local file="$1"
  if [[ -f "$file" ]]; then
    set -a
    # shellcheck disable=SC1090
    source "$file"
    set +a
    return 0
  fi
  return 1
}

load_env_file "${REPO_ROOT}/.env" || load_env_file "${REPO_ROOT}/backend/.env" || true

DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-dth_app}"
DB_USER="${DB_USER:-dth_user}"
DB_PASSWORD="${DB_PASSWORD:-}"

OUT_DIR="${1:-${REPO_ROOT}/infra/backups}"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"
OUT_FILE="${OUT_DIR}/dth_backup_${DB_NAME}_${TIMESTAMP}.sql.gz"
CHECKSUM_FILE="${OUT_FILE}.sha256"

if ! command -v pg_dump >/dev/null 2>&1; then
  echo "error: pg_dump not found in PATH" >&2
  exit 1
fi

if ! command -v gzip >/dev/null 2>&1; then
  echo "error: gzip not found in PATH" >&2
  exit 1
fi

mkdir -p "${OUT_DIR}"

export PGPASSWORD="${DB_PASSWORD}"

echo "Starting backup for database '${DB_NAME}' at ${DB_HOST}:${DB_PORT}"
pg_dump \
  --host "${DB_HOST}" \
  --port "${DB_PORT}" \
  --username "${DB_USER}" \
  --dbname "${DB_NAME}" \
  --no-owner \
  --no-privileges \
  | gzip -c > "${OUT_FILE}"

if command -v sha256sum >/dev/null 2>&1; then
  sha256sum "${OUT_FILE}" > "${CHECKSUM_FILE}"
fi

echo "Backup completed: ${OUT_FILE}"
if [[ -f "${CHECKSUM_FILE}" ]]; then
  echo "Checksum file: ${CHECKSUM_FILE}"
fi
