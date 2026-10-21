# SistemaRifas

> **Status: Deprecated / archived experiment**
>
> This repository is a historical snapshot of an experimental project. It is no
> longer maintained and receives no support, updates, or security fixes.
>
> **Do not deploy this application or expose it to the public.** It was built as
> a learning exercise and is not production-ready. See
> [Documentation.md](Documentation.md) for the full retrospective.

SistemaRifas was an experimental raffle SaaS aimed at content creators. The idea
was to let streamers and social-media creators run paid raffles with manual
payment verification and a publicly auditable draw, while participants could buy
tickets and submit proof of payment.

The project was also a practical exercise in building a full-stack application
with Go, MongoDB, React, TypeScript, and Docker Compose. It was dropped after the
initial prototype. The value of the repository is the experiment itself, not the
product.

## Tech stack

| Layer | Technology |
|---|---|
| Backend | Go 1.25, Chi router, MongoDB driver v2 |
| Database | MongoDB |
| Auth | JWT access tokens + refresh tokens, bcrypt |
| Frontend | React 19, TypeScript, Vite, MUI |
| Infrastructure | Docker Compose, Dockerfile |
| Documentation | Markdown |

## What worked

- Email/password registration and login.
- JWT access tokens with refresh-token rotation.
- Password hashing with bcrypt and refresh-token hashing.
- Raffle creation, editing, publishing, cancellation, and listing.
- Ticket purchase with per-raffle sequential ticket numbers.
- Manual payment proof submission and verifier approval/rejection through the API.
- Admin user listing, role updates, and status updates through the API.
- MongoDB indexes created automatically on startup.
- Commit/reveal draw logic using SHA-256 and a Fisher-Yates shuffle, though the
  full flow was not wired end-to-end.
- A responsive React frontend with landing, login, dashboard, raffle detail,
  management, live draw, verifier, my-tickets, and admin screens.
- Dark mode, snackbar notifications, loading states, and empty states.
- Frontend TypeScript checks passed at the time of archival; lint reported only
  three fast-refresh warnings.

## What was partial

- **Draw flow:** The draw service existed, but the API route needed to create the
  commit was never registered, so the draw could not be completed through the
  provided UI.
- **Verifier panel:** The page expected a pending-tickets endpoint that was not
  implemented. Verification worked through the API but not through that screen.
- **Admin panel:** The user-management table worked, but the statistics endpoint
  used by the dashboard was missing.
- **Live draw page:** It called additional draw endpoints that were not part of
  the backend router.
- **Authorization:** Authentication was implemented; role-based authorization was
  not consistently enforced on the server.
- **Ticket accounting:** Multi-ticket purchases and rejected-ticket inventory
  handling were not fully atomic or consistent.
- **Draw reproducibility:** Ticket ordering and random-number generation needed
  more work before external verification could be considered reliable.
- **Setup:** Local `.env` loading was incomplete, and the Dockerfile targeted an
  older Go version than `go.mod` required.

## What was never implemented

- Stripe PaymentIntents and webhook signature verification.
- Twitch OAuth.
- Real S3/Cloudflare R2 file uploads.
- Email verification.
- Pro-tier billing and publish gating.
- Profile editing.
- Admin bootstrap/seed flow.
- Automated tests, CI, and production hardening.

## Repository contents

```text
.
├── cmd/api/                  # Application entry point
├── internal/
│   ├── config/               # Environment configuration
│   ├── db/                   # MongoDB index setup
│   ├── domain/               # Domain models and constants
│   ├── dto/                  # Request/response DTOs
│   ├── handler/              # HTTP handlers
│   ├── middleware/           # Auth, CORS, logging, role helpers
│   ├── repository/           # MongoDB data access
│   ├── service/              # Business logic
│   └── utils/                # Hashing, JWT, MongoDB, errors
├── pkg/crypto/               # Commit/reveal seed and PRNG logic
├── migrations/               # Intended migration location (empty snapshot)
├── frontend/                 # React + TypeScript + Vite application
├── docker-compose.yml
├── Dockerfile
├── Documentation.md          # In-depth retrospective
└── LICENSE                   # MIT
```

## Running the prototype (historical)

These steps describe the intended development setup. They are preserved for
context and are not guaranteed to work in the current snapshot.

```powershell
Copy-Item .env.example .env
docker compose up -d mongo

# Terminal 1
go run ./cmd/api

# Terminal 2
cd frontend
pnpm install
pnpm dev
```

Known caveats:

- The Go application expects environment variables but does not load `.env`
  itself; it works with defaults only when those defaults are acceptable.
- `docker-compose.yml` expects a `.env` file.
- The Dockerfile targets an older Go version than `go.mod` requires.
- `migrations/` is empty; indexes are created by the application at startup.
- MongoDB must be running before the API starts.

## Security notice

Security was not production-grade. In particular, the prototype relied on
authentication without consistently enforcing authorization, used permissive CORS,
had a fallback JWT secret, exposed more ticket/payment metadata than intended in
some responses, and lacked rate limiting. These observations are documented at a
design level in [Documentation.md](Documentation.md).

Do not reuse this code as-is for anything involving real payments, personal data,
or public traffic.

## Documentation

For the full project history, architecture, feature-by-feature status, security
observations, and revival checklist, see [Documentation.md](Documentation.md).

## License

MIT. See [LICENSE](LICENSE).

The software is provided "as is", without warranty of any kind. The MIT license
does not make this abandoned prototype suitable for production use.
