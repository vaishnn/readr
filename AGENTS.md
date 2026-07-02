# Readr — Agent Context

## What Each Agent Should Know

### Backend (Go)
- Module path: `github.com/readr/api`
- Entry point: `backend/cmd/server/main.go`
- Services live in `backend/internal/services/` — each file maps to one domain (books, auth, users, highlights, notes, collections, bookmarks, progress)
- Handlers in `backend/internal/handlers/` are thin — validate input, call service, write JSON
- `BookService.Stream()` returns `(io.ReadCloser, int64, string, error)` — the handler proxies bytes directly, never redirects
- `MinioClient.PresignedURL()` replaces internal `minio:9000` with `MINIO_PUBLIC_URL` for browser-accessible cover URLs
- All list queries initialize with `make([]T, 0)` to avoid JSON `null`

### Frontend (Angular 17)
- All components are standalone; imports listed in `@Component({ imports: [...] })`
- Use signals (`signal()`, `computed()`) not RxJS for local state
- No arrow functions in templates — extract to component methods
- API calls go through `core/services/*.service.ts`; proxy forwards `/api` to backend
- pdfjs loaded via dynamic `import('pdfjs-dist')` with worker at `/assets/pdf.worker.min.js`
- Theme is toggled by adding/removing `light` class on `document.documentElement`

### Infrastructure
- `docker-compose up` starts: api (8080), frontend (4200), mongo (27017), redis (6379), minio (9000/9001)
- Rebuild backend: `docker-compose build --no-cache api && docker-compose restart api`
- Rebuild frontend: `docker-compose build --no-cache frontend && docker-compose restart frontend`
- MinIO console: http://localhost:9001 (minioadmin / minioadmin)

### Testing / Debugging
- Backend logs: `docker-compose logs -f api`
- Frontend logs: `docker-compose logs -f frontend`
- Check stream endpoint: `curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/books/<id>/stream -v`
- Check presigned URLs: response body of `GET /api/v1/books` — `cover_url` field should use `localhost:9000`, not `minio:9000`
