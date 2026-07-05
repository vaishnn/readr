# .github/workflows

GitHub Actions CI workflows.

## Files

- **ci.yml** — Single CI pipeline with three parallel jobs:
  - `backend`: Go vet + build + integration tests. Spins up MongoDB 7, Redis 7, and MinIO (bitnami) as service containers. Uses `go test -tags integration`. Env vars mirror what `tests/integration/` expects (`TEST_MONGO_URI`, `TEST_REDIS_ADDR`, `TEST_MINIO_*`).
  - `frontend`: `npm ci` → `tsc --noEmit` (type check) → `npm run build:prod`. No browser tests because Angular's Karma requires a display; type checking + prod build catches most regressions.
  - `docker`: Builds both `backend/Dockerfile` and `frontend/Dockerfile` (no push) after the other two jobs pass. Uses GHA layer cache.

## Constraints
- Integration tests need the `integration` build tag (`-tags integration`) or they are skipped.
- MinIO health check uses `curl` on `/minio/health/live`; the bitnami image exposes this.
- Frontend prod build target is `build:prod` (defined in `package.json` as `ng build --configuration production`).
