# Activity Tracker

Service API untuk mencatat aktivitas dan menampilkan ringkasan usage (daily/top) dengan PostgreSQL + Redis.

## Prasyarat

- Docker + Docker Compose (opsi termudah)
- Alternatif lokal: Go 1.22+ dan PostgreSQL + Redis

## Instalasi dengan Docker

1) Jalankan stack
```
docker compose up --build
```
2) API siap di `http://localhost:8080`

Catatan:
- Migrasi database dijalankan otomatis saat service API start.
- PostgreSQL tersedia di `localhost:5432` (user: `postgres`, password: `postgres`, db: `activity_tracker`).

## Instalasi Lokal (tanpa Docker)

1) Siapkan PostgreSQL dan Redis, lalu buat database `activity_tracker`.
2) Atur environment:
```
export DATABASE_URL="postgres://<user>:<pass>@<host>:<port>/activity_tracker?sslmode=disable"
export REDIS_ADDR="<host>:<port>"
export HTTP_ADDR=":8080"
```
3) Jalankan API:
```
go run ./cmd/api
```

## Konfigurasi Environment

- `DATABASE_URL` (default: `postgres://postgres:postgres@postgres:5432/activity_tracker?sslmode=disable`)
- `REDIS_ADDR` (default: `redis:6379`)
- `HTTP_ADDR` (default: `:8080`)
- `JWT_SECRET` (default: `rahasia`)

## Health Check

```
curl http://localhost:8080/health
```

## Troubleshooting Singkat

- Jika response `/api/usage/top` selalu `null`, pastikan ada log masuk lewat `/api/logs` dan cache `usage:top:last24h` tidak tersimpan dari state kosong (restart service atau hapus key di Redis).
