# Backend Architecture

Starter ini mengikuti referensi `Documents/Project/folder-structure-ai-starter.md`.

## Request Flow

```text
Client
  -> Gin router di app/routes.go
  -> Middleware di middleware/
  -> Handler di modules/<domain>/<version>/handler/
  -> Service di modules/<domain>/<version>/service/
  -> Repository di modules/<domain>/<version>/repository/
  -> Entity di entity/ + Database
```

## Folder Utama

- `app/`: registrasi route HTTP Gin.
- `config/`: loader environment dan config struct.
- `middleware/`: request ID, CORS, auth placeholder, RBAC placeholder, tenant context, rate limit placeholder.
- `modules/health/v1/`: contoh module lengkap dengan builder, handler, dan service.
- `response/`: wrapper response standar.
- `utils/`: helper koneksi PostgreSQL, Redis, HTTP, JWT placeholder, pagination, encryption.
- `entity/`: model awal tenant, user, dan workflow.

Starter ini belum mengimplementasikan auth, workflow CRUD, execution engine, schedule, atau AI business logic. Folder dan pola DI sudah disiapkan agar module tersebut dapat ditambahkan bertahap.
