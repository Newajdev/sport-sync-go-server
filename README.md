# SpotSync

SpotSync is a Go REST API for smart parking and EV charging reservations. Airports and malls can define parking zones, drivers can reserve spots, and admins can manage zones and view all bookings.

Built with Echo, GORM, PostgreSQL, JWT authentication, and role-based access (`driver` / `admin`).

## Features

- User registration and login with JWT
- Parking zone management (admin create, public list/get with live `available_spots`)
- Concurrency-safe reservations using GORM transactions and `FOR UPDATE` row locks
- Drivers can view and cancel their own reservations
- Admins can list all reservations in the system

## Tech Stack

- Go 
- Echo 
- GORM + PostgreSQL
- JWT 
- bcrypt
- go-playground/validator

## Project Structure

```text
sport-sync/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── auth/                   # JWT create/validate
│   ├── config/                 # Environment and database config
│   ├── domain/
│   │   ├── user/               # Auth (register, login)
│   │   ├── zone/               # Parking zones
│   │   └── reservation/        # Spot reservations
│   ├── httpresponse/           # SpotSync success/error envelope
│   ├── middlewares/            # Auth and admin middleware
│   └── server/                 # Echo setup and route registration
├── .air.toml                   # Air config for live reload
├── .env.example                # Example environment variables
├── go.mod
└── go.sum
```

## Architecture

Request flow follows strict layered architecture:

```text
Route -> Handler -> Service -> Repository -> Database
```

- **Handler** binds and validates DTOs, maps HTTP status codes, returns JSON envelopes.
- **Service** contains business rules (password hashing, capacity checks, ownership checks).
- **Repository** runs GORM queries, transactions, and row locks.
- **DTO** defines request/response shapes — GORM models are never exposed directly.

Example:

```text
POST /api/v1/reservations
    -> reservation handler
    -> reservation service
    -> reservation repository (transaction + FOR UPDATE)
    -> PostgreSQL
```

## Requirements

- Go
- PostgreSQL
- Git

## Setup

1. Clone the repository and enter the project directory:

```bash
git clone https://github.com/Newajdev/sport-sync-go-server.git
cd sport-sync
```

2. Create a PostgreSQL database:

```sql
CREATE DATABASE spotsync;
```

3. Copy the example environment file:

```bash
cp .env.example .env
```

4. Update `.env` with your database credentials:

```env
DSN="host=localhost user=postgres password=postgres dbname=spotsync port=5432 sslmode=disable"
PORT=8080
JWT_SECRET=change-this-secret
```

## Run Locally

Install dependencies:

```bash
go mod tidy
```

Start the server:

```bash
go run cmd/main.go
or
air
```



Health check:

```bash
curl http://localhost:8080/health
```

Expected response:

```text
running
```

## Build

```bash
go build -o server ./cmd/main.go
```


## API Endpoints

All responses use the SpotSync envelope:

```json
{ "success": true, "message": "...", "data": { } }
```

Errors:

```json
{ "success": false, "message": "...", "errors": "..." }
```

| Method | Path | Access |
|--------|------|--------|
| POST | `/api/v1/auth/register` | Public |
| POST | `/api/v1/auth/login` | Public |
| POST | `/api/v1/zones` | Admin |
| GET | `/api/v1/zones` | Public |
| GET | `/api/v1/zones/:id` | Public |
| POST | `/api/v1/reservations` | Auth (driver, admin) |
| GET | `/api/v1/reservations/my-reservations` | Auth |
| DELETE | `/api/v1/reservations/:id` | Auth (own reservation) |
| GET | `/api/v1/reservations` | Admin |

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@spotsync.com","password":"securePassword123","role":"driver"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@spotsync.com","password":"securePassword123"}'
```

Save the `token` from the response for protected routes.

### Create Zone (admin)

```bash
curl -X POST http://localhost:8080/api/v1/zones \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Terminal 1 EV Charging","type":"ev_charging","total_capacity":20,"price_per_hour":5.50}'
```

### List Zones

```bash
curl http://localhost:8080/api/v1/zones
```

### Create Reservation

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"zone_id":1,"license_plate":"ABC-1234"}'
```

### My Reservations

```bash
curl http://localhost:8080/api/v1/reservations/my-reservations \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Cancel Reservation

```bash
curl -X DELETE http://localhost:8080/api/v1/reservations/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### All Reservations (admin)

```bash
curl http://localhost:8080/api/v1/reservations \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

## Concurrency

Reservation creation uses a database transaction with a row-level lock (`FOR UPDATE`) on the parking zone. This prevents two drivers from booking the last available EV spot at the same time.

## Database Tables

GORM `AutoMigrate` creates these tables on startup:

- `users`
- `zones`
- `reservations`

# admin
```bash
ADMIN EMAIL= admin@spotsync.com
ADMIN PASSWORD= admin123456
ADMIN NAME= SpotSync Admin
```

# driver
```bash
ADMIN EMAIL= system.driver@spotsync.com
ADMIN PASSWORD= driver123456
ADMIN NAME= System Driver
```