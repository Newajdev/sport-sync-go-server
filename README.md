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

Protected routes require the header:

```text
Authorization: Bearer <token>
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

---

### Register User

**Endpoint:** `POST /api/v1/auth/register`

**Access:** Public

**Request Body**

```json
{
  "name": "John Doe",
  "email": "john@spotsync.com",
  "password": "securePassword123",
  "role": "driver"
}
```

**Success Response (201 Created)**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@spotsync.com",
    "role": "driver",
    "created_at": "2026-06-20T15:30:00Z",
    "updated_at": "2026-06-20T15:30:00Z"
  }
}
```

---

### Login

**Endpoint:** `POST /api/v1/auth/login`

**Access:** Public

**Request Body**

```json
{
  "email": "john@spotsync.com",
  "password": "securePassword123"
}
```

**Success Response (200 OK)**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@spotsync.com",
      "role": "driver"
    }
  }
}
```

---

### Create Zone

**Endpoint:** `POST /api/v1/zones`

**Access:** Admin

**Request Body**

```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 20,
  "price_per_hour": 5.50
}
```

**Success Response (201 Created)**

```json
{
  "success": true,
  "message": "Parking zone created successfully",
  "data": {
    "id": 5,
    "name": "Terminal 1 EV Charging",
    "type": "ev_charging",
    "total_capacity": 20,
    "price_per_hour": 5.50,
    "created_at": "2026-06-20T15:30:00Z",
    "updated_at": "2026-06-20T15:30:00Z"
  }
}
```

---

### List Zones

**Endpoint:** `GET /api/v1/zones`

**Access:** Public

**Request Body**

None

**Success Response (200 OK)**

```json
{
  "success": true,
  "message": "Parking zones retrieved successfully",
  "data": [
    {
      "id": 5,
      "name": "Terminal 1 EV Charging",
      "type": "ev_charging",
      "total_capacity": 20,
      "available_spots": 15,
      "price_per_hour": 5.50,
      "created_at": "2026-06-20T15:30:00Z",
      "updated_at": "2026-06-20T15:30:00Z"
    }
  ]
}
```

---

### Get Zone by ID

**Endpoint:** `GET /api/v1/zones/:id`

**Access:** Public

**Request Body**

None

**Success Response (200 OK)**

```json
{
  "success": true,
  "message": "Parking zone retrieved successfully",
  "data": {
    "id": 5,
    "name": "Terminal 1 EV Charging",
    "type": "ev_charging",
    "total_capacity": 20,
    "available_spots": 15,
    "price_per_hour": 5.50,
    "created_at": "2026-06-20T15:30:00Z",
    "updated_at": "2026-06-20T15:30:00Z"
  }
}
```

---

### Create Reservation

**Endpoint:** `POST /api/v1/reservations`

**Access:** Auth (driver, admin)

**Request Body**

```json
{
  "zone_id": 5,
  "license_plate": "ABC-1234"
}
```

**Success Response (201 Created)**

```json
{
  "success": true,
  "message": "Reservation confirmed successfully",
  "data": {
    "id": 105,
    "user_id": 1,
    "zone_id": 5,
    "license_plate": "ABC-1234",
    "status": "active",
    "created_at": "2026-06-20T15:30:00Z",
    "updated_at": "2026-06-20T15:30:00Z"
  }
}
```

---

### My Reservations

**Endpoint:** `GET /api/v1/reservations/my-reservations`

**Access:** Auth

**Request Body**

None

**Success Response (200 OK)**

```json
{
  "success": true,
  "message": "My reservations retrieved successfully",
  "data": [
    {
      "id": 105,
      "license_plate": "ABC-1234",
      "status": "active",
      "zone": {
        "id": 5,
        "name": "Terminal 1 EV Charging",
        "type": "ev_charging"
      },
      "created_at": "2026-06-20T15:30:00Z"
    }
  ]
}
```

---

### Cancel Reservation

**Endpoint:** `DELETE /api/v1/reservations/:id`

**Access:** Auth (own reservation only)

**Request Body**

None

**Success Response (200 OK)**

```json
{
  "success": true,
  "message": "Reservation cancelled successfully",
  "data": null
}
```

---

### All Reservations (Admin)

**Endpoint:** `GET /api/v1/reservations`

**Access:** Admin

**Request Body**

None

**Success Response (200 OK)**

```json
{
  "success": true,
  "message": "Reservations retrieved successfully",
  "data": [
    {
      "id": 105,
      "user_id": 1,
      "zone_id": 5,
      "license_plate": "ABC-1234",
      "status": "active",
      "user": {
        "id": 1,
        "name": "John Doe",
        "email": "john@spotsync.com",
        "role": "driver"
      },
      "zone": {
        "id": 5,
        "name": "Terminal 1 EV Charging",
        "type": "ev_charging"
      },
      "created_at": "2026-06-20T15:30:00Z",
      "updated_at": "2026-06-20T15:30:00Z"
    }
  ]
}
```

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