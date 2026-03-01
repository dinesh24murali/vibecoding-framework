# Infrastructure Assets

This folder contains deployment-time assets for Village DTH.

## Included files

- `systemd/dth-api.service`: API service unit (`/opt/dth/current/bin/dth-api`)
- `systemd/dth-worker.service`: worker service unit (`/opt/dth/current/bin/dth-worker`)
- `nginx/dth.conf`: reverse proxy and static frontend routing
- `deploy/deploy.sh`: one-shot deployment helper for a single VM

## Routing model

- `GET /api/v1/*` -> proxied to `127.0.0.1:8080`
- `GET /admin/*` -> serves admin static app from `/opt/dth/current/frontend/admin/out`
- `GET /*` -> serves customer static app from `/opt/dth/current/frontend/customer/out`

## Service environment

Both `dth-api` and `dth-worker` units read runtime environment from:

- `/etc/dth/dth.env`

Ensure this file includes all required backend variables (DB, JWT, Razorpay, reCAPTCHA, CORS).

## Deployment

Run from repository root:

```bash
bash infra/deploy/deploy.sh
```

Optional release root override:

```bash
APP_ROOT=/opt/dth bash infra/deploy/deploy.sh
```

The script builds binaries and static frontends, syncs release files, installs service/nginx configs, and restarts services.
