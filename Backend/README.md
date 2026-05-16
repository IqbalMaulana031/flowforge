# FlowForge Backend

Backend FlowForge adalah REST API berbasis Go/Gin dengan PostgreSQL dan Redis.

## Prasyarat

- Go 1.23 atau lebih baru
- PostgreSQL 15 atau kompatibel
- Redis 7 atau kompatibel
- golang-migrate CLI
- Make

Jika ingin menjalankan PostgreSQL dan Redis dari Docker Compose project ini:

```bash
docker-compose up -d postgres redis
```

## Setup local

1. Masuk ke folder backend.

```bash
cd App/Backend
```

2. Salin file environment.

```bash
cp .env.sample .env
```

3. Sesuaikan isi `.env` sesuai database dan Redis lokal.

Default yang sudah cocok dengan `docker-compose.yml`:

```env
APP_PORT=8080
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=flowforge
POSTGRES_PASSWORD=secret
POSTGRES_NAME=flowforge_db
POSTGRES_SSLMODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
```

4. Isi key enkripsi password.

Backend membutuhkan:

```env
PASSWORD_ENCRYPTION_PRIVATE_KEY_B64=<base64-pem-rsa-private-key>
```

Nilai public key pasang di frontend sebagai `VITE_PASSWORD_ENCRYPTION_PUBLIC_KEY_B64`.

5. Install dependency Go.

```bash
go mod download
```

Atau:

```bash
make dep-download
```

## Migrasi database

Migration menggunakan `golang-migrate`. Untuk command `make migrate` dan `make rollback`, parameter wajibnya adalah:

- `url`: connection string PostgreSQL tanpa query string `sslmode`.
- `module`: nama schema/module migration.

Contoh URL local:

```bash
export DB_URL="postgres://flowforge:secret@localhost:5432/flowforge_db"
```

### Urutan migrasi pertama kali

Migrasi schema harus dibuat dulu sebelum table, karena table berada di schema terpisah.

Urutan lengkap yang harus dijalankan:

1. Schema: `master`, `auth`, `workflow`, `execution`, `scheduler`
2. Extensions: `pgcrypto`
3. Module `master`
4. Module `auth`
5. Module `workflow`
6. Module `execution`
7. Module `scheduler`

Alasan dependency:

- `master.tenants` dipakai oleh `auth.users`, `workflow.workflows`, `execution.workflow_runs`, dan `scheduler.schedules`.
- `auth.users` dipakai oleh `workflow.workflows`, `workflow.workflow_versions`, dan `scheduler.schedules`.
- `workflow.workflows` dan `workflow.workflow_versions` dipakai oleh `execution.workflow_runs`.
- `execution.workflow_runs` dipakai oleh `execution.run_steps` dan `execution.execution_logs`.
- `scheduler.schedules` bergantung ke schema `workflow`, `master`, dan `auth`.

### Migrasi cepat semua schema dan module

Project sudah menyediakan target untuk menjalankan semua migrasi sesuai urutan di atas:

```bash
make all-db-migrate url=$DB_URL
```

Target ini menjalankan script `bin/migrate.sh` dengan urutan:

```text
db/schemas -> db/migrations/extensions -> db/migrations/master -> db/migrations/auth -> db/migrations/workflow -> db/migrations/execution -> db/migrations/scheduler
```

### Migrasi manual per schema/module

1. Buat schema terlebih dahulu.

```bash
make migrate-schema url=$DB_URL
```

Schema yang dibuat berurutan dari file migration:

```text
000001_master
000002_auth
000003_workflow
000004_execution
000005_scheduler
```

2. Jalankan migrasi table per module sesuai urutan dependency.

```bash
make migrate module=extensions url=$DB_URL
make migrate module=master url=$DB_URL
make migrate module=auth url=$DB_URL
make migrate module=workflow url=$DB_URL
make migrate module=execution url=$DB_URL
make migrate module=scheduler url=$DB_URL
```

### Cara rollback lalu migrate lagi

Gunakan rollback dari module terakhir ke module pertama agar foreign key tidak konflik.

Contoh rollback satu step untuk semua module, lalu migrate lagi:

```bash
make rollback module=scheduler url=$DB_URL
make rollback module=execution url=$DB_URL
make rollback module=workflow url=$DB_URL
make rollback module=auth url=$DB_URL
make rollback module=master url=$DB_URL
make rollback module=extensions url=$DB_URL
make rollback-schema url=$DB_URL
```

Setelah rollback, jalankan migrasi lagi dari awal:

```bash
make migrate-schema url=$DB_URL
make migrate module=extensions url=$DB_URL
make migrate module=master url=$DB_URL
make migrate module=auth url=$DB_URL
make migrate module=workflow url=$DB_URL
make migrate module=execution url=$DB_URL
make migrate module=scheduler url=$DB_URL
```

Untuk rollback semua version pada satu module:

```bash
make rollback-all module=<nama-module> url=$DB_URL
```

Untuk rollback semua schema:

```bash
make rollback-schema-all url=$DB_URL
```

## Seed data

Setelah migrasi berhasil, jalankan seeder jika membutuhkan data awal:

```bash
make seed
```

Seeder membaca file SQL di `db/seeds`.

## Menjalankan backend di local

Pastikan PostgreSQL dan Redis aktif, `.env` sudah benar, dan migrasi sudah dijalankan.

```bash
go run ./main.go
```

API akan berjalan di port sesuai `APP_PORT`, default:

```text
http://localhost:8080
```

Health endpoint:

```text
GET http://localhost:8080/v1/health
```

## Command Makefile yang sering dipakai

```bash
make dep-download        # download dependency Go
make tidy                # rapikan go.mod/go.sum
make test.unit           # jalankan unit test
make migrate-schema url=$DB_URL
make rollback-schema url=$DB_URL
make migrate module=<module> url=$DB_URL
make rollback module=<module> url=$DB_URL
make seed
```
