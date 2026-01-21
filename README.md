# Drakor Backend

Backend API untuk aplikasi streaming drakor menggunakan Golang + Gin + PostgreSQL (Supabase).

## Tech Stack

- **Framework**: Gin
- **Database**: PostgreSQL (Supabase)
- **Auth**: JWT (golang-jwt/jwt)
- **Validation**: go-playground/validator

## Setup

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Update `.env` dengan credentials Anda:
   ```
   DATABASE_URL=postgresql://postgres:PASSWORD@db.xxx.supabase.co:5432/postgres
   JWT_SECRET=your-secret-key
   PORT=8080
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run server:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

### Health Check
- `GET /health` - Check server status

### Auth (Phase 4)
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/auth/me`
- `PUT /api/auth/profile`

### More endpoints coming in future phases...

## Project Structure

```
server/
├── cmd/api/main.go      # Entry point
├── internal/            # Application logic
├── pkg/                 # Shared packages
├── migrations/          # SQL migrations
├── .env.example         # Environment template
└── go.mod               # Go module
```
