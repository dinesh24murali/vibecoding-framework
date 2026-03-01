#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

usage() {
  cat <<'EOF'
Usage:
  restore_db.sh <backup-file> [--recreate-db]

Examples:
  restore_db.sh infra/backups/dth_backup_dth_app_20260227_113000.sql.gz
  restore_db.sh infra/backups/dth_backup_dth_app_20260227_113000.sql --recreate-db
EOF
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

BACKUP_FILE="$1"
RECREATE_DB="false"

if [[ $# -ge 2 ]]; then
  if [[ "$2" == "--recreate-db" ]]; then
    RECREATE_DB="true"
  else
    usage
    exit 1
  fi
fi

if [[ ! -f "$BACKUP_FILE" ]]; then
  echo "error: backup file not found: $BACKUP_FILE" >&2
  exit 1
fi

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

if ! command -v psql >/dev/null 2>&1; then
  echo "error: psql not found in PATH" >&2
  exit 1
fi

export PGPASSWORD="${DB_PASSWORD}"

if [[ "${RECREATE_DB}" == "true" ]]; then
  if ! command -v dropdb >/dev/null 2>&1 || ! command -v createdb >/dev/null 2>&1; then
    echo "error: dropdb/createdb are required when using --recreate-db" >&2
    exit 1
  fi

  echo "Recreating database '${DB_NAME}'"
  dropdb --if-exists --host "${DB_HOST}" --port "${DB_PORT}" --username "${DB_USER}" "${DB_NAME}"
  createdb --host "${DB_HOST}" --port "${DB_PORT}" --username "${DB_USER}" "${DB_NAME}"
fi

echo "Restoring backup into database '${DB_NAME}'"
if [[ "${BACKUP_FILE}" == *.gz ]]; then
  gunzip -c "${BACKUP_FILE}" | psql \
    --host "${DB_HOST}" \
    --port "${DB_PORT}" \
    --username "${DB_USER}" \
    --dbname "${DB_NAME}" \
    --set ON_ERROR_STOP=1
else
  psql \
    --host "${DB_HOST}" \
    --port "${DB_PORT}" \
    --username "${DB_USER}" \
    --dbname "${DB_NAME}" \
    --set ON_ERROR_STOP=1 \
    --file "${BACKUP_FILE}"
fi

echo "Restore completed from: ${BACKUP_FILE}"
