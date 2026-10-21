# PENDING.md — Stubs to implement later

## 1. Stripe Integration
- **File:** `internal/handler/webhook_handler.go`
- **Endpoint:** `POST /api/v1/webhooks/stripe`
- **What's needed:**
  - Verify Stripe webhook signature (`STRIPE_WEBHOOK_SECRET`)
  - Handle `payment_intent.succeeded` → update ticket status to `verified`
  - Handle `payment_intent.payment_failed` → cancel ticket
  - Create PaymentIntent in `ticket_service.go` for automatic payments
- **Dependency:** `github.com/stripe/stripe-go/v78`

## 2. Twitch OAuth
- **File:** `internal/handler/auth_handler.go` — `POST /api/v1/auth/oauth/twitch`
- **What's needed:**
  - Redirect to Twitch OAuth authorize URL
  - Exchange authorization code for access token
  - Fetch Twitch user profile
  - Create/find user by `oauth_provider: "twitch"` + `oauth_id`
  - Generate JWT tokens and return
- **Dependencies:** Twitch API HTTP client, OAuth2 flow

## 3. S3/R2 File Upload
- **File:** `internal/service/ticket_service.go` — UploadProof
- **What's needed:**
  - Generate presigned upload URLs or accept multipart uploads
  - Store files in Cloudflare R2 / AWS S3
  - Store URL in `ticket.payment_proof_url`
- **Dependency:** `github.com/aws/aws-sdk-go-v2`

## 4. Email Verification
- **Field:** `user.email_verified`
- **What's needed:**
  - Send verification email on registration
  - Endpoint to verify email token
- **Dependency:** SMTP client or email service (Resend, SendGrid, etc.)

## 5. Pro Tier Payment
- **Endpoint:** `POST /api/v1/raffles/:id/publish` — charge Stripe before publishing
- **What's needed:**
  - Check raffle `pricing_tier == "pro"`
  - Create Stripe Checkout session or charge before activating

---

## 6. Frontend UI — Current State ✅ (Completed 2026-07-29)

### What works now
- Login / Register with tabs + post-login redirect awareness
- Landing page with hero section + "Cómo funciona" 3-step guide
- Public raffle listing (paginated cards, hover animations)
- Raffle detail page with ticket purchase form + owner "Gestionar Rifa" button
- Proof upload nudge after manual transfer purchase
- Dashboard (my raffles + create modal) with loading states
- Raffle manage page (publish, cancel, sort, commit, re-sort) + cancel confirmation dialog
- Live draw page with ticket-spinning animation + suspense build-up
- Verifier panel (approve/reject pending tickets) with loading states
- My tickets page (list + upload receipt button) + ticket links to raffle
- Admin page (user table, promote/suspend, dashboard stats) + suspend confirmation dialog
- Global snackbar/toast feedback system on all pages
- Loading states (CircularProgress) on every page that fetches data
- Empty states with icons + guidance CTAs on all lists
- Mobile responsive: hamburger drawer menu, responsive spacing, stacked grids
- JWT auto-refresh interceptor
- Role-based navigation and route protection

### Remaining nice-to-haves (Low Priority)
- [x] ~~Dark mode toggle~~ ✅ (2026-07-29)
- [x] ~~Raffle countdown timer on cards (draw_date)~~ ✅ (2026-07-29)
- [x] ~~Copy-to-clipboard for commit hash (live draw verification)~~ ✅ (2026-07-29)
- [x] ~~Pagination on every list (verifier, manage tickets, admin)~~ ✅ (2026-07-29)
- [ ] **Profile page (edit name, social links):** Needs `PATCH /users/me` backend endpoint.
- [ ] **Attach actual Stripe.js frontend flow for card payments:** Needs Stripe backend (#1).

---

## 7. MongoDB Without Docker

If you don't want Docker Desktop, use **MongoDB Atlas** (free cloud tier):

1. Sign up at https://www.mongodb.com/atlas — free M0 cluster
2. Create a cluster, whitelist your IP (or use 0.0.0.0/0 for dev)
3. Get the connection string: `mongodb+srv://user:pass@cluster.mongodb.net`
4. Set in `.env`: `MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net`
5. No Docker needed — runs locally with `go run ./cmd/api`

Alternative: install MongoDB Community Edition directly on Windows from https://www.mongodb.com/try/download/community
