# Readr — Claude Context

## Project Overview
Readr is a self-hosted digital library app. Users upload PDFs (and eventually EPUBs/CBZ), read them in-browser, and manage collections, bookmarks, highlights, and notes.

## Stack
- **Frontend**: Angular 17 (standalone components, signals), Tailwind CSS v3, pdfjs-dist v3
- **Backend**: Go, Chi router, MongoDB, MinIO (object storage), Redis, bcrypt/JWT auth
- **Infrastructure**: Docker Compose; MinIO bucket `readr`

## Key Architecture Decisions
- **PDF streaming**: Backend proxies the file directly (no redirect to MinIO). Handler reads from MinIO via `BookService.Stream()` and writes to `http.ResponseWriter`.
- **MinIO presigned URLs**: Internal hostname `minio:9000` is replaced with `localhost:9000` (`MINIO_PUBLIC_URL` env var) so browsers can reach cover images.
- **Auto-thumbnail**: On PDF upload, the frontend renders page 1 with pdfjs into a canvas (no auth needed — uses local ArrayBuffer) and sends the JPEG blob as `cover` in the multipart upload.
- **Theme**: Dark/light toggle via `html.light` CSS class overrides (no Tailwind dark mode refactor). Persisted in localStorage and synced to `PATCH /users/me/settings`.
- **Angular templates**: No arrow functions allowed in templates — always use component methods.
- **pdfjs worker**: Served from `/assets/pdf.worker.min.js` (copied from `node_modules/pdfjs-dist/build/` via `angular.json` assets).

## Directory Structure
```
backend/
  cmd/server/main.go       — entry point, wires services/handlers
  internal/
    config/                — env-based config (MinioPublicURL, etc.)
    database/              — MongoDB connection + collection helpers
    handlers/              — HTTP handlers (books, auth, users, …)
    middleware/            — JWT auth, CORS
    models/                — BSON/JSON structs
    services/              — business logic (BookService, UserService, …)
    storage/               — MinioClient (Upload, Stream, PresignedURL, Delete)
frontend/
  src/app/
    core/
      models/              — TypeScript interfaces (Book, User, …)
      services/            — Angular services (BookService, ThemeService, …)
    features/
      library/             — Book grid, upload modal (auto-thumbnail)
      reader/              — PDF viewer (continuous scroll), notes panel
      settings/            — Password change, theme picker
    shared/components/     — Spinner, Toast, Navbar
  proxy.conf.json          — /api → http://api:8080
  angular.json             — proxyConfig, pdfjs worker asset
```

## Folder-level Context Files

Every directory worked on must have a `claude.md` file describing the files and responsibilities of that folder. Rules:
- **Create** `claude.md` in any folder when first working in it.
- **Update** `claude.md` immediately whenever a file in that folder is created, deleted, or significantly changed — keep it current, not historical.
- **Content**: list each file, its purpose, key exports/types/routes, and any non-obvious constraints. One entry per file, concise.
- **Goal**: Claude should be able to understand a folder's contents from `claude.md` alone, without re-reading every source file.

## Common Gotchas
- Go nil slices marshal to JSON `null` — always use `make([]T, 0)`.
- Docker named volume `frontend_node_modules` prevents stale packages after rebuild.
- `docker-compose build --no-cache api` required when Go binary changes don't take effect.
- pdfjs-dist v4 requires TypeScript 5.7; Angular 17 only supports TS 5.4 — stay on v3.
