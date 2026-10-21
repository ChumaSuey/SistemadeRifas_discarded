# Setup Guide — Sistema de Rifas

## Prerequisites

Install these once:

### Windows

1. **Docker Desktop** — download and install from https://www.docker.com/products/docker-desktop/
   - After installation, open Docker Desktop and wait until the whale icon in the system tray says "Engine running"
   - This gives you both `docker` and `docker compose` commands

2. **Go 1.22+** — download from https://go.dev/dl/
   - Run the installer, accept defaults
   - Verify with: `go version`

3. **Node.js 18+** — download from https://nodejs.org/
   - Install the LTS version
   - This gives you `node` and `npm`

4. **pnpm** — install with npm:
   ```powershell
   npm install -g pnpm
   ```
   - Verify with: `pnpm --version`

### Verify installation

Open PowerShell and run:

```powershell
docker --version
docker compose version
go version
node --version
pnpm --version
```

All should print version numbers.

---

## Quick Start (First Time)

```powershell
# 1. Copy env file
Copy-Item .env.example .env

# 2. Start MongoDB ONLY (not the API, we'll run that manually for development)
docker compose up -d mongo
```

Wait ~10 seconds for MongoDB to boot. You'll see:
```
[+] Running 2/2
 ✔ Network rafflesaas_default  Created
 ✔ Container rifas-mongo-1     Started
```

```powershell
# 3. Install frontend dependencies
cd frontend
pnpm install
cd ..

# 4. Run the API (Terminal 1)
go run ./cmd/api

# 5. Run the frontend (Terminal 2)
cd frontend
pnpm dev
```

You should see:
```
INFO connecting to MongoDB uri=mongodb://localhost:27017
INFO connected to MongoDB
INFO running database index migrations
INFO database index migrations completed
INFO server starting port=8080
```

The API is now listening at `http://localhost:8080` and the frontend at `http://localhost:5173`.
Open `http://localhost:5173` in your browser.

---

## Quick Test

```powershell
# Register a user
$body = @{name="Test"; email="test@test.com"; password="password123"} | ConvertTo-Json
$response = Invoke-RestMethod -Method POST -Uri http://localhost:8080/api/v1/auth/register -Body $body -ContentType "application/json"
$response | ConvertTo-Json

# Save the access token
$token = $response.access_token

# Check your profile
Invoke-RestMethod -Uri http://localhost:8080/api/v1/auth/me -Headers @{Authorization="Bearer $token"} | ConvertTo-Json

# Create a raffle
$raffle = @{title="Mi primera rifa"; description="Probando"; ticket_price=5.00; currency="USD"; tickets_total=50; payment_methods=@("manual_transfer")} | ConvertTo-Json
Invoke-RestMethod -Method POST -Uri http://localhost:8080/api/v1/raffles -Body $raffle -ContentType "application/json" -Headers @{Authorization="Bearer $token"} | ConvertTo-Json
```

---

## Day-to-Day Commands

```powershell
# Start MongoDB
docker compose up -d mongo

# Check if MongoDB is running
docker compose ps

# Stop MongoDB (keeps your data)
docker compose stop mongo

# Stop MongoDB AND delete all data (fresh start)
docker compose down -v

# View MongoDB logs
docker compose logs mongo

# Run the API (Terminal 1)
go run ./cmd/api

# Run the frontend (Terminal 2)
cd frontend
pnpm dev

# Build frontend for production
cd frontend
pnpm build
```

---

## MongoDB Cheat Sheet

Connect to the database directly to inspect data:

```powershell
# Open a MongoDB shell inside the container
docker compose exec mongo mongosh raffle_saas
```

Inside the shell:

```js
// Show all collections
show collections

// See all users
db.users.find().pretty()

// See all raffles
db.raffles.find().pretty()

// See all tickets
db.tickets.find().pretty()

// Delete everything (reset)
db.users.drop()
db.raffles.drop()
db.tickets.drop()
db.audit_logs.drop()
db.refresh_tokens.drop()

// Exit
exit
```

---

## API Endpoints Quick Reference

### Auth (public)
| Method | Endpoint | Body |
|--------|----------|------|
| POST | `/api/v1/auth/register` | `{name, email, password}` |
| POST | `/api/v1/auth/login` | `{email, password}` |
| POST | `/api/v1/auth/refresh` | `{refresh_token}` |

### Auth (authenticated)
| Method | Endpoint | Notes |
|--------|----------|-------|
| POST | `/api/v1/auth/logout` | `{refresh_token}` |
| GET | `/api/v1/auth/me` | User profile |

### Raffles (public)
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/api/v1/raffles` | List active raffles |
| GET | `/api/v1/raffles/{id}` | Detail + tickets |

### Raffles (raffler only)
| Method | Endpoint | Notes |
|--------|----------|-------|
| POST | `/api/v1/raffles` | Create (draft) |
| GET | `/api/v1/raffles/my` | My raffles |
| PATCH | `/api/v1/raffles/{id}` | Update (draft only) |
| POST | `/api/v1/raffles/{id}/publish` | Draft → Active |
| POST | `/api/v1/raffles/{id}/cancel` | Cancel raffle |
| POST | `/api/v1/raffles/{id}/sort` | Execute draw |
| POST | `/api/v1/raffles/{id}/resort` | Re-sort (disqualify winner) |

### Tickets (authenticated)
| Method | Endpoint | Notes |
|--------|----------|-------|
| POST | `/api/v1/raffles/{id}/tickets` | Buy tickets `{quantity, payment_method}` |
| GET | `/api/v1/raffles/{id}/tickets` | List tickets (raffler/verifier) |
| GET | `/api/v1/tickets/my` | My tickets (participant) |
| POST | `/api/v1/tickets/{id}/upload-proof` | `{payment_proof_url}` |
| POST | `/api/v1/tickets/{id}/verify` | Approve (verifier) |
| POST | `/api/v1/tickets/{id}/reject` | Reject `{reason}` (verifier) |

### Admin
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/api/v1/users` | List users (admin) |
| PATCH | `/api/v1/users/{id}/role` | Update role `{role}` (admin) |
| PATCH | `/api/v1/users/{id}/status` | Update status `{status}` (admin) |

---

## Troubleshooting

**"MongoDB connection refused"**
- Docker Desktop might not be running. Open it and wait for "Engine running".
- MongoDB might still be starting. Wait 10-15 seconds and retry.

**"Port 27017 already in use"**
- You might have another MongoDB running. Stop it or change the port in `.env` and `docker-compose.yml`.

**"Port 8080 already in use"**
- Change `APP_PORT` in `.env` to something else like `8081`.

**Docker takes forever to start**
- First run downloads the MongoDB image (~700MB). Subsequent starts are instant.

**Reset everything**
```powershell
docker compose down -v
docker compose up -d mongo
go run ./cmd/api
```
