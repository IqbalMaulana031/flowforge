# FlowForge Frontend

Frontend FlowForge adalah aplikasi React + Vite.

## Prasyarat

- Node.js 18 atau lebih baru
- npm
- Backend FlowForge berjalan di local

## Setup local

1. Masuk ke folder frontend.

```bash
cd "App/Front End"
```

2. Install dependency.

```bash
npm install
```

3. Salin file environment.

```bash
cp .env.sample .env
```

4. Sesuaikan isi `.env`.

Default local:

```env
VITE_API_BASE_URL=http://localhost:8080/v1
VITE_WS_URL=ws://localhost:8080/ws
VITE_PASSWORD_ENCRYPTION_PUBLIC_KEY_B64=<base64-der-rsa-public-key>
```

Catatan:

- `VITE_API_BASE_URL` harus mengarah ke base URL backend.
- `VITE_WS_URL` harus mengarah ke WebSocket backend.
- `VITE_PASSWORD_ENCRYPTION_PUBLIC_KEY_B64` harus pasangan public key dari private key backend `PASSWORD_ENCRYPTION_PRIVATE_KEY_B64`.

## Menjalankan frontend di local

Pastikan backend sudah berjalan terlebih dahulu.

```bash
npm run dev
```

Vite akan menampilkan URL local di terminal. Secara default aplikasi dapat dibuka di:

```text
http://localhost:5173
```

Script `dev` menggunakan `vite --host 0.0.0.0`, sehingga aplikasi juga dapat diakses dari network yang sama jika firewall mengizinkan.

## Build production

```bash
npm run build
```

Output build akan dibuat oleh Vite di folder `dist`.

## Preview build production

```bash
npm run preview
```

## Urutan menjalankan full local

1. Jalankan PostgreSQL dan Redis.
2. Setup `.env` backend.
3. Jalankan migrasi database backend.
4. Jalankan backend dengan `go run ./main.go` dari folder `App/Backend`.
5. Setup `.env` frontend.
6. Jalankan frontend dengan `npm run dev` dari folder `App/Front End`.
