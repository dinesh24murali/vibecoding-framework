# Backup and Restore Runbook

This runbook covers backup and restore for the Village DTH PostgreSQL database.

## Scope

- Core entities covered: `users`, `providers`, `plans`, `service_requests` and related tables.
- Scripts:
  - `infra/scripts/backup_db.sh`
  - `infra/scripts/restore_db.sh`

## Prerequisites

- PostgreSQL client tools installed: `pg_dump`, `psql`, `gzip`.
- Optional for `--recreate-db`: `dropdb`, `createdb`.
- Database access from execution host.
- Environment variables set in either:
  - repository root `.env`, or
  - `backend/.env`

Expected variables:

- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`

## Backup Procedure

1. Ensure production writes are paused or low-traffic window is active.
2. Run backup:

```bash
bash infra/scripts/backup_db.sh
```

3. Optional custom output directory:

```bash
bash infra/scripts/backup_db.sh /absolute/path/to/backups
```

4. Verify output:
   - backup file exists with timestamped name, example:
     `dth_backup_dth_app_20260301_143000.sql.gz`
   - checksum file exists when `sha256sum` is available (`.sha256`).

## Restore Procedure

### Option A: Restore into existing database

```bash
bash infra/scripts/restore_db.sh infra/backups/dth_backup_dth_app_20260301_143000.sql.gz
```

### Option B: Drop/recreate database, then restore

```bash
bash infra/scripts/restore_db.sh infra/backups/dth_backup_dth_app_20260301_143000.sql.gz --recreate-db
```

## Post-Restore Verification

Run basic checks:

```sql
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM providers;
SELECT COUNT(*) FROM plans;
SELECT COUNT(*) FROM service_requests;
```

Application smoke checks:

1. `POST /api/v1/auth/admin/login`
2. `GET /api/v1/customer/providers`
3. `POST /api/v1/customer/checkout/service-requests`

## Ops Drill (Recommended Monthly)

1. Take fresh backup from active DB.
2. Restore into a staging DB with `--recreate-db`.
3. Run SQL row-count checks and API smoke checks.
4. Record drill date, backup filename, and verification result.

## Failure Handling

- If backup fails, stop and fix credentials/network before retrying.
- If restore fails, keep failed logs and retry from same backup in a clean DB.
- Never run destructive restore on production without a fresh backup from the same day.
