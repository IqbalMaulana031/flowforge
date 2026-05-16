# Backend Setup

FlowForge backend starter berada di folder `App/Backend` dan menggunakan Go, Gin, Uber Fx, GORM, dan Redis.

## Prerequisite

- Go 1.23+
- PostgreSQL 15+
- Redis 7+
- golang-migrate CLI jika ingin menjalankan migration

## Local Run

```bash
cp .env.sample .env
go mod tidy
go run ./main.go
```

API berjalan di `http://localhost:8080` secara default.

## Health Check

```bash
curl http://localhost:8080/
curl http://localhost:8080/v1/health/ping
curl http://localhost:8080/v1/health
```

Endpoint `/v1/health` akan melakukan ping ke PostgreSQL dan Redis menggunakan konfigurasi `.env`.
