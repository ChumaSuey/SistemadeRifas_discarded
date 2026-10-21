# SistemaRifas — Project Retrospective and Technical Documentation

> **Status: Deprecated / archived experiment**
>
> This document describes an abandoned prototype. The project is not maintained,
> receives no security updates, and is not suitable for production use.
>
> **Do not deploy this application or expose it to real users, payments, or
> personal data.**

- **Project:** SistemaRifas
- **Type:** Experimental full-stack SaaS prototype
- **Last source changes observed:** July 29, 2026
- **Archived:** October 2026
- **License:** MIT

---

## 1. Purpose of this document

This is an archival record of the project: what it was meant to become, what was
built, what worked, what remained unfinished, and why it was dropped. It is
intended for anyone revisiting the repository, including the original author.

The project is preserved as a learning artifact. Nothing here should be read as a
commitment to maintain or revive it.

---

## 2. The original idea

SistemaRifas was conceived as a raffle platform for content creators. A creator
could publish a raffle, sell numbered tickets, allow participants to pay by manual
bank transfer or card, have payment proofs verified, and then run a draw that could
be audited by the public.

The intended differentiators were:

- **Transparency:** a commit/reveal draw flow with a public hash.
- **Manual payment support:** useful in regions where card payments are unreliable.
- **Verification workflow:** a verifier reviews payment screenshots before a ticket
  becomes eligible for the draw.
- **Live event experience:** a dedicated draw screen intended for streaming.
- **Multi-currency support:** USD, VES, and EUR.

The project was also a personal exercise in:

- structuring a Go API in layers,
- using MongoDB with the v2 driver,
- building a React + TypeScript SPA,
- using Docker Compose for local infrastructure,
- designing an API-first SaaS workflow.

It was never validated as a business, and the legal and operational complexity of
paid raffles was one reason the experiment stopped.

---

## 3. Scope at archival time

The prototype covered the main happy path of raffle creation and ticket
participation, but several integrations were stubs. The most important unfinished
areas were payments, OAuth, file storage, email verification, and production-grade
authorization.

The repository should be treated as a partial vertical slice, not a usable
product.

---

## 4. Status summary

| Area | Status | Summary |
|---|---|---|
| Backend architecture | Functional | Layered Go application with repositories, services, handlers, middleware, and DTOs |
| Database setup | Functional | Indexes created automatically at startup |
| Email/password auth | Functional | Registration, login, JWT access tokens, refresh tokens |
| Role-based authorization | Partial | Authentication worked; authorization was not consistently enforced |
| Raffle lifecycle | Functional | Create, edit draft, publish, cancel, list |
| Ticket purchase | Partial | Worked, but transactions and inventory accounting had gaps |
| Manual payment proof | Functional | Proof URL can be submitted and reviewed through the API |
| Verifier UI | Partial | Page existed; expected endpoint was missing |
| Draw logic | Partial | Commit/reveal and shuffle logic existed; the flow was not fully wired |
| Live draw UI | Partial | Page existed; called missing endpoints |
| Admin API | Functional but unsafe | Endpoints existed, but role enforcement was incomplete |
| Admin UI | Partial | User table worked; statistics endpoint was missing |
| Stripe | Not implemented | Webhook stub only; no PaymentIntent flow |
| Twitch OAuth | Not implemented | Endpoint returned HTTP 501 |
| S3/R2 uploads | Not implemented | Only a proof URL string was accepted |
| Email verification | Not implemented | Field existed; workflow did not |
| Pro tier | Not implemented | Tier existed in the model; billing did not |
| Profile editing | Not implemented | No `PATCH /users/me` endpoint |
| Frontend | Mostly built | Nine feature pages plus shared layout, auth, and theme |
| Tests | None | No Go or frontend test suite |
| CI | None | No GitHub Actions or equivalent |
| Docker | Partial | Compose file existed; Dockerfile Go version was inconsistent |
| Documentation | Partial | Setup and pending-work notes existed; this retrospective was added at archival |

---

## 5. Technology stack

### Backend

- **Language:** Go 1.25
- **Router:** `github.com/go-chi/chi/v5`
- **Database:** MongoDB via `go.mongodb.org/mongo-driver/v2`
- **Auth:** JWT via `github.com/golang-jwt/jwt/v5`
- **Passwords:** bcrypt via `golang.org/x/crypto`
- **Logging:** standard library `log/slog`

### Frontend

- **Framework:** React 19
- **Language:** TypeScript
- **Build tool:** Vite
- **UI library:** MUI
- **HTTP client:** Axios
- **Routing:** React Router
- **Package manager:** pnpm

### Infrastructure

- **MongoDB:** Docker Compose service
- **API:** Dockerfile for a containerized Go binary
- **Local development:** API and Vite dev server run separately

---

## 6. Architecture

### 6.1 System context

```text
Browser
  |
  | HTTP/JSON
  v
React SPA (Vite dev server, port 5173)
  |
  | Axios, Bearer access token
  v
Go API (port 8080)
  |
  | MongoDB driver
  v
MongoDB (port 27017)

Planned but not implemented:
  API -> Stripe
  API -> Twitch OAuth
  API -> S3 / Cloudflare R2
  API -> email provider
```

### 6.2 Backend structure

```text
cmd/api/main.go
  |
  +-- internal/config        Environment-based configuration
  +-- internal/db            MongoDB index setup
  +-- internal/domain        Domain structs and constants
  +-- internal/dto           Request and response shapes
  +-- internal/repository    MongoDB access layer
  +-- internal/service       Business rules and workflows
  +-- internal/handler       HTTP request/response handling
  +-- internal/middleware    Auth, CORS, logging, role helpers
  +-- internal/utils         JWT, password hashing, token hashing, errors
  +-- pkg/crypto             Commit/reveal seed generation and shuffle
```

The layering was intentionally conventional:

- **Handlers** decode requests, extract the authenticated user, and write responses.
- **Services** enforce business rules and orchestrate repositories.
- **Repositories** contain MongoDB queries.
- **Domain** contains database-facing structs.
- **DTOs** control what the API exposes.
- **Middleware** handles cross-cutting HTTP concerns.

### 6.3 Frontend structure

```text
frontend/src/
├── api/client.ts            Axios instance and token refresh interceptor
├── components/              Layout, navbar, route guard, raffle card, ticket card
├── context/                 Auth, snackbar, and theme providers
├── pages/                   Home, login, raffle detail, dashboard, manage,
│                            live draw, verifier, my tickets, admin
├── types/                   Shared TypeScript interfaces
└── utils/                   Formatters and token storage
```

The frontend used role-aware navigation and protected routes on the client side.
That was useful for UX, but it was not a substitute for server-side authorization.

### 6.4 Data model

The MongoDB collections were:

| Collection | Purpose |
|---|---|
| `users` | Accounts, roles, status, OAuth placeholders, social links |
| `raffles` | Raffle metadata, pricing, status, draw commit/result fields |
| `tickets` | Numbered tickets, payment method, status, proof, verifier |
| `audit_logs` | State-change records intended for traceability |
| `refresh_tokens` | Hashed refresh tokens with expiry |

Key fields included:

- **User:** name, email, email verification flag, password hash, OAuth provider/id,
  role, status, social links, timestamps.
- **Raffle:** raffler id, title, description, status, pricing tier, ticket price,
  currency, payment methods, total/sold ticket counts, winner ids, sort seed and
  commit hash, draw date, timestamps.
- **Ticket:** raffle id, participant id, ticket number, status, payment method,
  proof URL, Stripe payment intent id, rejection reason, verifier id, verification
  time, creation time.
- **Audit log:** raffle id, user id, action, details, IP address, user agent,
  creation time.
- **Refresh token:** user id, token hash, expiry, creation time.

### 6.5 Database indexes

Indexes were created at application startup:

- Unique index on `users.email`.
- Index on `users.oauth_provider` + `users.oauth_id`.
- Index on `raffles.raffler_id` + `raffles.status`.
- Unique index on `tickets.raffle_id` + `tickets.ticket_number`.
- Index on `tickets.raffle_id` + `tickets.status`.
- Index on `tickets.participant_id`.
- Index on `audit_logs.raffle_id` + `audit_logs.created_at`.
- Unique index on `refresh_tokens.token_hash`.
- TTL index on `refresh_tokens.expires_at`.

The `migrations/` directory was intended for migration scripts but remained empty
in the snapshot.

---

## 7. Feature inventory

### 7.1 Implemented and functional

- User registration with name, email, and password.
- User login.
- Password hashing with bcrypt.
- JWT access token generation and validation.
- Opaque refresh tokens stored as SHA-256 hashes.
- Refresh-token rotation on refresh.
- Logout by deleting a refresh token.
- Authenticated profile lookup.
- Raffle creation.
- Raffle title/description updates while in draft.
- Raffle publishing with a free-tier ticket limit check.
- Raffle cancellation.
- Public listing of active raffles with pagination.
- Public raffle detail lookup.
- Creator-specific raffle listing.
- Ticket purchase with sequential ticket numbers.
- Prevention of creator participation in their own raffle.
- Manual payment proof URL submission.
- Verifier approval and rejection through the API.
- Prevention of self-verification by the verifier.
- Admin user listing.
- Admin role updates.
- Admin status updates.
- Audit-log writes for key actions.
- Commit/reveal seed generation and hash verification.
- Deterministic Fisher-Yates shuffle implementation.
- Re-sort behavior that cancels the previous winner and starts a new draw.
- Responsive frontend navigation.
- Login and registration UI.
- Landing page.
- Public raffle cards and detail page.
- Raffle creation modal.
- Raffle management page.
- Live draw page with animation.
- Verifier page.
- My tickets page.
- Admin page.
- Dark mode.
- Snackbar notifications.
- Client-side route protection.
- Automatic access-token refresh using the stored refresh token.

### 7.2 Partial or broken at archival time

- **Draw flow through the UI:** the commit service and handler existed, but the
  route was not registered. Since the draw required a prior commit, the draw could
  not be completed through the provided interface.
- **Live draw page:** it called endpoints that were not present in the backend
  router, including a separate commit/reveal sequence.
- **Verifier page:** it requested a pending-tickets endpoint that was never
  implemented, so the page could not populate its queue.
- **Admin page:** user management worked, but the statistics card used a missing
  endpoint and remained at zero.
- **Authorization:** authentication was consistently applied, but role checks were
  not consistently wired into protected routes.
- **Ticket purchase atomicity:** each ticket in a multi-ticket purchase was created
  in its own transaction, so a partial success was possible if a later ticket
  failed.
- **Inventory accounting:** rejected and cancelled tickets did not reduce the sold
  count, so ticket numbers were not released as originally intended.
- **Payment-method validation:** the purchase endpoint accepted payment methods
  that the raffle did not necessarily offer.
- **Stripe as a purchase option:** the frontend offered Stripe, and the API
  accepted it, but no payment flow existed. Such tickets remained pending.
- **Draw reproducibility:** the verified-ticket query had no deterministic sort, the
  PRNG used platform-dependent integer arithmetic, and the audit entry did not
  include the full ticket order.
- **Draw date:** raffle creation did not set a draw date, so countdowns on cards did
  not appear in normal use.
- **Public raffle detail:** the response included ticket records, including
  participant and payment-proof metadata, which was more exposure than intended.
- **Result retrieval:** draw results were available in the immediate response, but
  there was no dedicated public endpoint for later verification.
- **Admin cancellation:** only the raffle owner could cancel; the original design
  allowed admins to cancel too.
- **Audit metadata:** audit entries were written with empty IP address and user
  agent values.
- **Environment loading:** the Go application read environment variables but did not
  load `.env` itself, so the documented local setup did not apply the file as
  expected.
- **Docker build:** the Dockerfile used Go 1.22 while `go.mod` required Go 1.25.

### 7.3 Not implemented

- Stripe PaymentIntents.
- Stripe webhook signature verification.
- Webhook-based ticket verification or cancellation.
- Twitch OAuth login.
- Google OAuth, despite the data model allowing an OAuth provider.
- Presigned S3/R2 uploads.
- Multipart proof uploads.
- Email verification.
- Transactional email.
- Pro-tier charging before publishing.
- Profile editing.
- Social-link editing.
- Admin bootstrap or seeded admin account.
- Dispute or suspicious-payment workflow.
- Notifications.
- Export features.
- Rate limiting.
- Automated tests.
- CI/CD.
- Production deployment configuration.

---

## 8. Business rules: intended vs implemented

The original design listed core rules. The implementation matched some of them only
partially.

| Rule | Intent | Archival status |
|---|---|---|
| R1 | Raffle creator cannot buy their own ticket | Implemented in the purchase service |
| R2 | Ticket purchase must be atomic and respect total inventory | Partial: per-ticket transactions, not whole-order atomicity |
| R3 | Draw uses only verified tickets and fails with zero verified tickets | Implemented |
| R4 | Raffle fields cannot be changed after activation | Stricter than planned: updates only allowed in draft |
| R5 | Ticket numbers are unique per raffle | Implemented with a unique compound index |
| R6 | A verifier cannot verify their own ticket; verification is audited | Implemented, though audit metadata was incomplete |
| R7 | Free tier limited to a fixed number of tickets | Checked at publish time, not throughout the lifecycle |
| R8 | Cancellation rules and paid-ticket cancellation | Implemented for owners; admin cancellation missing |

---

## 9. Authentication and authorization

### 9.1 What worked

- Passwords were hashed with bcrypt.
- Access tokens were signed JWTs with user id and role claims.
- Refresh tokens were random opaque values.
- Refresh tokens were stored as hashes, not plaintext.
- Refresh-token rotation invalidated the old token.
- A frontend interceptor retried requests after refreshing the access token.
- Client-side protected routes checked the stored user role.

### 9.2 What was incomplete

- The role middleware helper existed but was not consistently applied to the route
  tree.
- Admin and verifier operations relied mainly on authentication rather than robust
  role enforcement.
- Access tokens continued to represent a user even if that user's status changed.
- Suspended users were not fully blocked by every authentication path.
- Registration assigned the `raffler` role to every new account. The `participant`
  role was defined but not assigned through a normal flow.
- There was no bootstrap path for creating the first admin account.
- Role and status values accepted by admin endpoints were not strictly validated.

These points are recorded as design-level observations. They are one reason the
project must not be deployed.

---

## 10. Draw and transparency design

### 10.1 Intended flow

1. The raffler starts the commit phase.
2. The backend generates a random nonce and timestamp.
3. A seed string is formed from the raffle id, nonce, and timestamp.
4. The SHA-256 hash of the seed is published as the commit hash.
5. The raffler later reveals the seed and executes the draw.
6. Only verified tickets participate.
7. A Fisher-Yates shuffle uses the seed as deterministic input.
8. The ticket in the first position after shuffling wins.
9. The winner, winning ticket, seed, commit hash, and draw date are stored.
10. The audit log records the event.

### 10.2 What existed

- Seed generation from a 32-byte cryptographic nonce and timestamp.
- SHA-256 commit hash.
- Hash verification before execution.
- Fisher-Yates shuffle over ticket numbers.
- Winner lookup by ticket number.
- Raffle status transition to finished.
- Winning ticket status transition to winner.
- Re-sort path that cancels the previous winning ticket.

### 10.3 Why it was not yet trustworthy

- The verified-ticket query had no explicit deterministic ordering, so the shuffle
  could depend on database return order.
- The PRNG used architecture-dependent integer width and modulo arithmetic.
- A public verification tool was not built.
- The full ordered participant list was not stored in the audit entry.
- The commit route was not wired into the API.
- The frontend expected a different reveal endpoint.
- There was no persistent public result endpoint that exposed the seed for later
  verification.

The concept was a reasonable start, but more work was needed before claiming
external auditability.

---

## 11. API surface at archival time

### 11.1 Registered routes

#### Auth

| Method | Route | Intended access |
|---|---|---|
| POST | `/api/v1/auth/register` | Public |
| POST | `/api/v1/auth/login` | Public |
| POST | `/api/v1/auth/oauth/twitch` | Public, but returned 501 |
| POST | `/api/v1/auth/refresh` | Public with refresh token |
| POST | `/api/v1/auth/logout` | Authenticated |
| GET | `/api/v1/auth/me` | Authenticated |

#### Raffles

| Method | Route | Intended access |
|---|---|---|
| GET | `/api/v1/raffles` | Public |
| GET | `/api/v1/raffles/{id}` | Public |
| POST | `/api/v1/raffles` | Authenticated |
| GET | `/api/v1/raffles/my` | Authenticated |
| PATCH | `/api/v1/raffles/{id}` | Authenticated, owner |
| POST | `/api/v1/raffles/{id}/publish` | Authenticated, owner |
| POST | `/api/v1/raffles/{id}/cancel` | Authenticated, owner |
| POST | `/api/v1/raffles/{id}/sort` | Authenticated, owner |
| POST | `/api/v1/raffles/{id}/resort` | Authenticated, owner |
| POST | `/api/v1/raffles/{id}/tickets` | Authenticated |
| GET | `/api/v1/raffles/{id}/tickets` | Authenticated |

#### Tickets

| Method | Route | Intended access |
|---|---|---|
| GET | `/api/v1/tickets/my` | Authenticated |
| POST | `/api/v1/tickets/{id}/upload-proof` | Authenticated, ticket owner |
| POST | `/api/v1/tickets/{id}/verify` | Authenticated, verifier |
| POST | `/api/v1/tickets/{id}/reject` | Authenticated, verifier |

#### Webhooks

| Method | Route | Intended access |
|---|---|---|
| POST | `/api/v1/webhooks/stripe` | Stripe signature, stub only |

#### Users

| Method | Route | Intended access |
|---|---|---|
| GET | `/api/v1/users` | Admin |
| PATCH | `/api/v1/users/{id}/role` | Admin |
| PATCH | `/api/v1/users/{id}/status` | Admin |

### 11.2 Routes expected by the frontend but missing

| Method | Route | Used by | Backend state |
|---|---|---|---|
| POST | `/api/v1/raffles/{id}/commit` | Manage page, live draw | Handler existed; route not registered |
| POST | `/api/v1/raffles/{id}/recompute-and-reveal` | Live draw | Not implemented |
| GET | `/api/v1/tickets/pending` | Verifier page | Not implemented |
| GET | `/api/v1/admin/stats` | Admin page | Not implemented |

These mismatches meant the frontend and backend were not fully integrated at
archival time.

---

## 12. Payments

### 12.1 Manual payments

The intended manual flow was:

1. Participant selects tickets and chooses manual transfer.
2. Tickets are created as pending.
3. Participant submits a proof URL.
4. A verifier approves or rejects the ticket.
5. Approved tickets become eligible for the draw.

The API supported creating pending tickets, submitting a proof URL, and
approving/rejecting tickets. It did not actually upload or store files; it only
accepted a URL. The verifier UI could not load its queue because the expected
endpoint was missing.

### 12.2 Stripe

Stripe was planned but not implemented:

- No PaymentIntent creation in the ticket service.
- No Stripe SDK dependency.
- The webhook endpoint read the body and returned success without verifying a
  signature or changing ticket state.
- The frontend still showed a Stripe option, and the API accepted it, which could
  create tickets that had no automatic payment path.

### 12.3 File storage

S3/R2 was planned but not implemented:

- No upload endpoint.
- No presigned URL generation.
- No multipart handling.
- No storage SDK dependency.

---

## 13. Frontend retrospective

### 13.1 What was built

- **Landing page:** hero section, three-step explanation, active raffle list, and
  pagination.
- **Login/register:** tabbed form with redirect awareness.
- **Dashboard:** creator raffle list and creation dialog.
- **Raffle detail:** raffle information, ticket purchase form, and owner management
  link.
- **Manage raffle:** publish, cancel, commit, sort, re-sort, and verified-ticket list.
- **Live draw:** animated draw experience with commit hash and seed display.
- **Verifier:** pending-ticket review UI.
- **My tickets:** participant ticket list and proof URL submission.
- **Admin:** user table, role promotion, suspension dialog, and statistics cards.
- **Shared UI:** responsive navbar/drawer, dark mode, snackbar provider, loading
  states, empty states, and protected routes.

### 13.2 Quality checks at archival time

- TypeScript compilation with `tsc --noEmit` passed.
- Oxlint reported three fast-refresh warnings related to context files exporting
  both providers and hooks.
- No frontend tests existed.
- The production build was not re-run as part of the archival review.

### 13.3 Integration gaps

The frontend was more complete than the backend. Several screens were built against
an expected API that was not finished. The most visible consequences were:

- the draw could not be started with the commit endpoint;
- the live draw called an endpoint that did not exist;
- the verifier queue stayed empty;
- admin statistics stayed at zero.

The UI was still useful as a design and component exercise.

---

## 14. Security and design observations

This section stays at the design level. It intentionally does not provide
exploitation steps, secrets, or instructions for causing failures.

The prototype was not production-ready. Notable observations include:

- **Authorization was incomplete.** Authentication existed, but protected routes
  did not consistently enforce the required role. Some administrative or
  verification operations were more accessible than intended.
- **Client-side route guards were relied on for UX, not security.** Hiding a page
  in the frontend does not protect the corresponding API.
- **Public responses exposed too much ticket data.** The public raffle detail
  response included ticket records with participant and payment-proof metadata.
- **CORS was permissive.** The API allowed any origin.
- **A fallback JWT secret existed.** Running without configuration could produce a
  predictable signing setup.
- **Authentication state was not revalidated.** Account status changes did not
  immediately invalidate existing access tokens.
- **Some authentication paths did not fully account for suspended accounts.**
- **No rate limiting was implemented.** Authentication and purchase endpoints had
  no abuse controls.
- **Input validation was minimal.** Email format, password strength, allowed roles,
  allowed statuses, and several domain values were not comprehensively validated.
- **Audit metadata was incomplete.** IP address and user agent were not populated.
- **Refresh-token lifetime was hardcoded.** The configured refresh TTL was not used
  by the token service.
- **Environment loading was incomplete.** The app read environment variables but
  did not load `.env` for local `go run`.
- **Payment integration was unsafe by omission.** The UI and API accepted a payment
  method for which no payment flow existed.
- **No tests covered authorization or business rules.**

Any revival should begin with a security review and a clean authorization model.

---

## 15. Infrastructure and setup notes

### 15.1 Intended local workflow

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

### 15.2 Environment variables

The example environment file defined:

- server environment and port,
- JWT secret and token lifetimes,
- MongoDB URI and database name,
- Twitch OAuth credentials,
- Stripe credentials,
- S3/R2 credentials and bucket settings.

The example values were placeholders. No real credentials are part of this
documentation.

### 15.3 Known setup problems

- `.env` was not loaded automatically by the Go process.
- Docker Compose expected a `.env` file, but the snapshot did not contain one.
- The Dockerfile used Go 1.22 while the module declared Go 1.25.
- The frontend API base URL was hardcoded to `http://localhost:8080`.
- `migrations/` was empty; index setup lived in application code.
- MongoDB had to be running before the API could start.
- There was no production frontend serving strategy in the repository.

### 15.4 Docker

The Compose file defined:

- an API service built from the Dockerfile,
- a MongoDB 7 service with a persistent volume and health check.

The Compose setup was useful for local MongoDB practice. It was not production
configuration.

---

## 16. Testing and quality

At archival time:

- No Go tests existed.
- No frontend unit, integration, or end-to-end tests existed.
- No CI pipeline existed.
- No coverage reporting existed.
- No static analysis pipeline beyond local linting existed for the frontend.
- The frontend TypeScript check passed.
- Oxlint produced three fast-refresh warnings.
- The backend build and full application flow were not validated during the
  archival review.

The absence of tests is one of the largest gaps in the prototype.

---

## 17. Why the project was dropped

The project was a practice experiment, not a committed product. The main reasons it
stopped were:

1. **Scope:** A raffle platform touches payments, identity, fraud prevention,
   storage, notifications, auditing, and legal compliance.
2. **Payment complexity:** Manual verification alone is an operational burden.
   Stripe integration and webhook reliability add more work.
3. **Trust requirements:** A transparent draw needs a carefully designed,
   deterministic, independently verifiable process.
4. **Business and legal uncertainty:** Paid raffles are regulated differently across
   jurisdictions. That made the idea less attractive as a casual side project.
5. **Time:** The remaining work was larger than the initial prototype suggested.
6. **Learning goal reached:** The project achieved its purpose as an exercise in
   Go, MongoDB, React, and Docker Compose.

The repository was kept as a record of that experiment.

---

## 18. Revival checklist

If this project were ever restarted, the following work would be required before
calling it usable:

1. **Validate the business and legal model.** Confirm whether paid raffles are
   permitted in the target jurisdictions.
2. **Rebuild authorization.** Enforce roles server-side, define permissions
   explicitly, and add an admin bootstrap flow.
3. **Fix the API contract.** Register the commit route, implement the missing
   endpoints, or remove the frontend features that depend on them.
4. **Make purchases atomic.** Treat a multi-ticket order as one transaction and
   define how rejected or cancelled inventory is released.
5. **Redesign the draw.** Use deterministic ordering, a well-specified PRNG, a
   public verification tool, and a complete audit record.
6. **Implement payments properly.** Add Stripe PaymentIntents, idempotent webhook
   handling, signature verification, and clear failure states.
7. **Add real storage.** Use presigned uploads or multipart handling with an object
   store.
8. **Add email verification and notifications.**
9. **Harden configuration and security.** Remove fallback secrets, validate
   environment variables, restrict CORS, add rate limiting, and revalidate account
   status.
10. **Add tests.** Cover authorization, purchases, cancellation, verification, and
    draw determinism.
11. **Add CI.** Run formatting, linting, tests, and builds for backend and
    frontend.
12. **Add observability and deployment configuration.**
13. **Review all exposed data.** Use explicit response DTOs for public endpoints.

---

## 19. Repository metadata

- **Language:** Go and TypeScript
- **Primary use:** Learning / archival
- **Maintenance:** None
- **Support:** None
- **Tests:** None
- **CI:** None
- **Local git history at archival time:** None present in the snapshot
- **License:** MIT

Recommended GitHub settings:

- Archive the repository.
- Add a short repository description such as:
  `Deprecated experimental raffle SaaS prototype built with Go, MongoDB, React, and Docker Compose.`
- Suggested topics: `deprecated`, `experiment`, `go`, `mongodb`, `react`,
  `typescript`, `docker-compose`, `raffle`, `saas-prototype`.

---

## 20. Final note

SistemaRifas was an ambitious experiment for a small practice project. It produced
a working core API, a broad frontend, and a useful exploration of layered Go
services, MongoDB, JWT authentication, and Docker Compose.

It also demonstrated why production software in this domain is difficult:
authorization, payment integrity, deterministic auditing, and operational trust are
not optional details. The project is preserved with those lessons documented rather
than hidden.

**Do not deploy it.**
