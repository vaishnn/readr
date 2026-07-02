# Readr

A self-hosted digital library for reading and managing PDFs (and eventually EPUBs). Upload your books, read them in the browser, and track your progress — no cloud required.

![License](https://img.shields.io/badge/license-MIT-blue)
![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)
![Angular](https://img.shields.io/badge/Angular-17-DD0031?logo=angular)

## Features

- **Upload & manage** PDF books with auto-generated cover thumbnails
- **In-browser reader** powered by [ngx-extended-pdf-viewer](https://github.com/stephanrauh/ngx-extended-pdf-viewer)
- **Reading progress** saved automatically (page, session time)
- **Highlights & notes** with color-coded text selection
- **Dark / light theme** with persistent preference
- **JWT authentication** with refresh tokens
- Self-hosted — all data stays on your server

## Tech Stack

| Layer | Tech |
|---|---|
| Frontend | Angular 17 (standalone components, signals), Tailwind CSS v3 |
| Backend | Go 1.22, Chi router |
| Database | MongoDB 7 |
| Storage | MinIO (S3-compatible object storage) |
| Cache | Redis 7 |
| Infrastructure | Docker Compose |

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose

### Run locally

```bash
git clone https://github.com/vaishnn/readr.git
cd readr
docker compose up --build
```

The app will be available at **http://localhost:4200**.

| Service | URL |
|---|---|
| Frontend | http://localhost:4200 |
| API | http://localhost:8080 |
| MinIO console | http://localhost:9001 |

Default MinIO credentials: `minioadmin` / `minioadmin`

### Environment variables

The default `docker-compose.yml` includes sensible development defaults. For production, override these at minimum:

```env
JWT_SECRET=<strong-random-secret>
JWT_REFRESH_SECRET=<strong-random-secret>
MINIO_ACCESS_KEY=<your-key>
MINIO_SECRET_KEY=<your-secret>
MINIO_PUBLIC_URL=https://your-domain.com:9000
```

## Project Structure

```
readr/
├── backend/
│   ├── cmd/server/        # Entry point
│   └── internal/
│       ├── config/        # Env-based config
│       ├── handlers/      # HTTP handlers (books, auth, users)
│       ├── middleware/    # JWT auth, CORS
│       ├── models/        # BSON/JSON structs
│       ├── services/      # Business logic
│       └── storage/       # MinIO client
├── frontend/
│   └── src/app/
│       ├── core/          # Models, services
│       ├── features/      # Library, reader, settings, auth
│       └── shared/        # Spinner, toast, navbar
├── k8s/                   # Kubernetes manifests (optional)
└── docker-compose.yml
```

## Development

### Backend only

```bash
cd backend
go run ./cmd/server
```

### Frontend only

```bash
cd frontend
npm install
npm start
```

### Rebuild after Go changes

```bash
docker compose build --no-cache api && docker compose up api
```

## Contributing

Pull requests are welcome. For larger changes, please open an issue first to discuss what you'd like to change.

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes
4. Open a pull request

## License

[MIT](LICENSE)
