# 🎟️ START.md — Sistema de Rifas (SaaS Comercial)
# Versión: 1.0 | Estado: Diseño Técnico Completo

### 6. Modelo de datos en GO

package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string        `bson:"name" json:"name"`
	Email         string        `bson:"email" json:"email"`
	EmailVerified bool          `bson:"email_verified" json:"email_verified"`
	PasswordHash  *string       `bson:"password_hash,omitempty" json:"-"`
	OAuthProvider *string       `bson:"oauth_provider,omitempty" json:"oauth_provider,omitempty"`
	OAuthID       *string       `bson:"oauth_id,omitempty" json:"oauth_id,omitempty"`
	Role          string        `bson:"role" json:"role"` // admin, raffler, participant, verifier
	Status        string        `bson:"status" json:"status"` // active, suspended, banned
	SocialLinks   struct {
		Twitch    string `bson:"twitch,omitempty" json:"twitch,omitempty"`
		Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"`
	} `bson:"social_links" json:"social_links"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Raffle struct {
	ID               bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	RafflerID        bson.ObjectID  `bson:"raffler_id" json:"raffler_id"`
	Title            string         `bson:"title" json:"title"`
	Description      string         `bson:"description" json:"description"`
	Status           string         `bson:"status" json:"status"` // draft, pending_verification, active, finished, cancelled
	PricingTier      string         `bson:"pricing_tier" json:"pricing_tier"` // free, pro
	TicketPrice      float64        `bson:"ticket_price" json:"ticket_price"`
	Currency         string         `bson:"currency" json:"currency"` // USD, VES, EUR
	PaymentMethods   []string       `bson:"payment_methods" json:"payment_methods"` // manual_transfer, stripe
	TicketsTotal     int            `bson:"tickets_total" json:"tickets_total"`
	TicketsSold      int            `bson:"tickets_sold" json:"tickets_sold"`
	WinnerID         *bson.ObjectID `bson:"winner_id,omitempty" json:"winner_id,omitempty"`
	WinningTicketID  *bson.ObjectID `bson:"winning_ticket_id,omitempty" json:"winning_ticket_id,omitempty"`
	SortSeed         *string        `bson:"sort_seed,omitempty" json:"sort_seed,omitempty"`
	SortCommitHash   *string        `bson:"sort_commit_hash,omitempty" json:"sort_commit_hash,omitempty"`
	DrawDate         *time.Time     `bson:"draw_date,omitempty" json:"draw_date,omitempty"`
	CreatedAt        time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time      `bson:"updated_at" json:"updated_at"`
}

type Ticket struct {
	ID                     bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	RaffleID               bson.ObjectID  `bson:"raffle_id" json:"raffle_id"`
	ParticipantID          bson.ObjectID  `bson:"participant_id" json:"participant_id"`
	TicketNumber           int            `bson:"ticket_number" json:"ticket_number"`
	Status                 string         `bson:"status" json:"status"` // pending, verified, rejected, cancelled, winner
	PaymentMethod          string         `bson:"payment_method" json:"payment_method"` // manual_transfer, stripe
	PaymentProofURL        *string        `bson:"payment_proof_url,omitempty" json:"payment_proof_url,omitempty"`
	StripePaymentIntentID  *string        `bson:"stripe_payment_intent_id,omitempty" json:"stripe_payment_intent_id,omitempty"`
	RejectionReason        *string        `bson:"rejection_reason,omitempty" json:"rejection_reason,omitempty"`
	VerifiedBy             *bson.ObjectID `bson:"verified_by,omitempty" json:"verified_by,omitempty"`
	VerifiedAt             *time.Time     `bson:"verified_at,omitempty" json:"verified_at,omitempty"`
	CreatedAt              time.Time      `bson:"created_at" json:"created_at"`
}

type AuditLog struct {
	ID        bson.ObjectID          `bson:"_id,omitempty" json:"id"`
	RaffleID  bson.ObjectID          `bson:"raffle_id" json:"raffle_id"`
	UserID    bson.ObjectID          `bson:"user_id" json:"user_id"`
	Action    string                 `bson:"action" json:"action"`
	Details   map[string]interface{} `bson:"details" json:"details"`
	IPAddress string                 `bson:"ip_address" json:"ip_address"`
	UserAgent string                 `bson:"user_agent" json:"user_agent"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
}

type RefreshToken struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id" json:"user_id"`
	TokenHash string        `bson:"token_hash" json:"-"`
	ExpiresAt time.Time     `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}


### 7. Indices de Base de datos

Ejecutar al iniciar la aplicación (script en migrations/ o en main.go):

// go
// users
db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys:    bson.D{{Key: "email", Value: 1}},
    Options: options.Index().SetUnique(true),
})
db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{Key: "oauth_provider", Value: 1}, {Key: "oauth_id", Value: 1}},
})

// raffles
db.Collection("raffles").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{Key: "raffler_id", Value: 1}, {Key: "status", Value: 1}},
})

// tickets — ÍNDICE ÚNICO CRÍTICO para evitar números duplicados
db.Collection("tickets").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys:    bson.D{{Key: "raffle_id", Value: 1}, {Key: "ticket_number", Value: 1}},
    Options: options.Index().SetUnique(true),
})
db.Collection("tickets").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{Key: "raffle_id", Value: 1}, {Key: "status", Value: 1}},
})
db.Collection("tickets").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{Key: "participant_id", Value: 1}},
})

// audit_logs
db.Collection("audit_logs").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys: bson.D{{Key: "raffle_id", Value: 1}, {Key: "created_at", Value: -1}},
})

// refresh_tokens
db.Collection("refresh_tokens").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys:    bson.D{{Key: "token_hash", Value: 1}},
    Options: options.Index().SetUnique(true),
})
db.Collection("refresh_tokens").Indexes().CreateOne(ctx, mongo.IndexModel{
    Keys:    bson.D{{Key: "expires_at", Value: 1}},
    Options: options.Index().SetExpireAfterSeconds(0), // TTL
})

### 8. Flujo de Pagos
#### 8.1. Pago Manual (Transferencia Bancaria)
Participant selecciona tickets y elige "Transferencia Bancaria".
Sistema reserva los números (crea tickets con status: "pending").
Participant sube screenshot del comprobante a S3/R2 → payment_proof_url.
Verifier revisa el panel de verificación, aprueba o rechaza.
Si aprueba: status → "verified", verified_by → verifier ID.
Si rechaza: status → "rejected", rejection_reason → texto, números liberados.
#### 8.2. Pago Automático (Stripe)
Participant selecciona tickets y elige "Tarjeta (Stripe)".
Backend crea PaymentIntent con Stripe.
Frontend confirma el pago con Stripe.js.
Webhook payment_intent.succeeded actualiza tickets a status: "verified".
Si falla: status → "cancelled", números liberados.
#### 8.3. Anti-fraude en Pagos Manuales
Hash del archivo de comprobante se almacena en audit_logs.
Validación: un mismo payment_proof_url no puede usarse para múltiples tickets de la misma rifa.
Verifier puede marcar como "sospechoso" para revisión de Admin.
### 9. Algoritmo de Sorteo Transparente
#### 9.1. Fase de Commit (Antes del Sorteo)
El Raffler inicia el sorteo desde el dashboard.
Backend genera:
nonce = 32 bytes aleatorios criptográficos
timestamp = Unix nanoseconds
seed_input = SHA-256(raffle_id + nonce + timestamp)
sort_commit_hash = SHA-256(seed_input) (se publica antes)
Se guarda sort_commit_hash en la rifa.
Se muestra públicamente: hash, cantidad de tickets verificados, timestamp.
#### 9.2. Fase de Reveal (Ejecución del Sorteo)
Se revela seed_input (la semilla real).
Se recolectan todos los tickets con status == "verified" de la rifa.
Se ejecuta Fisher-Yates usando seed_input como semilla de un PRNG determinista.
El ticket en posición 0 del array barajado es el ganador.
Se guarda:
sort_seed = seed_input
winner_id = participant del ticket ganador
winning_ticket_id = ID del ticket
draw_date = timestamp
status = "finished"
Se registra en audit_logs con lista completa de tickets participantes y orden final del shuffle.
#### 9.3. Verificación Externa
Cualquier persona puede verificar:
plain
1. Tomar el sort_commit_hash publicado antes del sorteo.
2. Tomar la seed revelada después.
3. Calcular SHA-256(seed) y comparar con sort_commit_hash.
4. Ejecutar Fisher-Yates con la seed sobre la lista de tickets verificados.
5. Confirmar que el resultado coincide con el ganador oficial.
#### 9.4. Re-Sorteo
Solo permitido si el ganador es marcado como "disqualified" por el Raffler (con razón documentada).
Se genera un nuevo sort_commit_hash independiente.
El ticket anteriormente ganador queda con status: "cancelled".
Se ejecuta el proceso de commit-reveal nuevamente.



#### 10. API REST — Endpoints

🔐 Auth

| Método | Endpoint                    | Descripción                     | Auth                    |
| ------ | --------------------------- | ------------------------------- | ----------------------- |
| POST   | `/api/v1/auth/register`     | Registro email/password         | Público                 |
| POST   | `/api/v1/auth/login`        | Login email/password            | Público                 |
| POST   | `/api/v1/auth/oauth/twitch` | Login/Registro con Twitch OAuth | Público                 |
| POST   | `/api/v1/auth/refresh`      | Refresh access token            | Público (refresh token) |
| POST   | `/api/v1/auth/logout`       | Invalidar refresh token         | Bearer                  |
| GET    | `/api/v1/auth/me`           | Perfil del usuario autenticado  | Bearer                  |


🎟️ Raffles

| Método | Endpoint                      | Descripción                            | Auth    | Roles                |
| ------ | ----------------------------- | -------------------------------------- | ------- | -------------------- |
| POST   | `/api/v1/raffles`             | Crear rifa                             | Bearer  | Raffler, Admin       |
| GET    | `/api/v1/raffles`             | Listar rifas públicas (activas)        | Público | —                    |
| GET    | `/api/v1/raffles/:id`         | Detalle de rifa pública                | Público | —                    |
| GET    | `/api/v1/raffles/my`          | Listar mis rifas (creador)             | Bearer  | Raffler              |
| PATCH  | `/api/v1/raffles/:id`         | Actualizar rifa (solo draft)           | Bearer  | Raffler dueño        |
| POST   | `/api/v1/raffles/:id/publish` | Publicar rifa (draft → active)         | Bearer  | Raffler dueño        |
| POST   | `/api/v1/raffles/:id/cancel`  | Cancelar rifa                          | Bearer  | Raffler dueño, Admin |
| POST   | `/api/v1/raffles/:id/sort`    | Ejecutar sorteo                        | Bearer  | Raffler dueño        |
| POST   | `/api/v1/raffles/:id/resort`  | Re-sortear (disqualify + nuevo sorteo) | Bearer  | Raffler dueño        |


🎫 Tickets

| Método | Endpoint                           | Descripción                | Auth   | Roles                          |
| ------ | ---------------------------------- | -------------------------- | ------ | ------------------------------ |
| POST   | `/api/v1/raffles/:id/tickets`      | Comprar/reservar tickets   | Bearer | Participant                    |
| GET    | `/api/v1/raffles/:id/tickets`      | Listar tickets de una rifa | Bearer | Raffler dueño, Verifier, Admin |
| GET    | `/api/v1/tickets/my`               | Mis tickets (participante) | Bearer | Participant                    |
| POST   | `/api/v1/tickets/:id/upload-proof` | Subir comprobante de pago  | Bearer | Participant dueño              |
| POST   | `/api/v1/tickets/:id/verify`       | Verificar ticket (aprobar) | Bearer | Verifier, Admin                |
| POST   | `/api/v1/tickets/:id/reject`       | Rechazar ticket            | Bearer | Verifier, Admin                |

💳 Stripe Webhooks
| Método | Endpoint | Descripción | Auth |
|---|---|---|
| POST | /api/v1/webhooks/stripe | Webhook de Stripe | Stripe Signature |


### 11. Reglas de Negocio Core
#### R1. Restricción de Participación
ParticipantID NUNCA puede ser igual a Raffle.RafflerID.
Validación en capa de servicio antes de crear ticket. Error: 409 Conflict.
#### R2. Concurrencia en Compra
Usar transacción multi-documento de MongoDB:
Crear ticket con ticket_number único.
Incrementar tickets_sold en raffles.
Validar tickets_sold <= tickets_total.
Índice único (raffle_id, ticket_number) garantiza atomicidad.
#### R3. Sorteo Solo con Tickets Verificados
Sort() solo considera tickets con status == "verified".
Si hay 0 tickets verificados, error 422 Unprocessable Entity.
#### R4. Raffle No Modificable Después de Activa
Solo campos editables en status "active": description (menor).
No se puede cambiar ticket_price, tickets_total, currency.
#### R5. Un Ticket por Número
Garantizado por índice único compuesto en MongoDB.
#### R6. Verificación de Comprobantes
Un Verifier no puede verificar su propio ticket (auto-validación prohibida).
Cada verificación genera entrada en audit_logs.
#### R7. Límite de Tickets (Tier Free)
Tier "free": máximo 100 tickets por rifa.
Tier "pro": sin límite (o límite configurable por Admin).
#### R8. Cancelación de Rifa
Solo permitida si status es "draft" o "active".
Si se cancela en "active": todos los tickets pasan a "cancelled", se notifica a participants.
No se puede cancelar si status == "finished".

### 12. Modelo de Monetización

| Tier     | Precio         | Límite de Tickets | Funciones                                                                                                   |
| -------- | -------------- | ----------------- | ----------------------------------------------------------------------------------------------------------- |
| **Free** | \$0            | 100 tickets       | Sorteo básico, 1 método de pago manual                                                                      |
| **Pro**  | \$X por evento | Ilimitado         | Múltiples métodos de pago, Stripe integrado, soporte prioritario, verificadores ilimitados, export de datos |

El pago del tier Pro se procesa antes de publicar la rifa (publish endpoint).


#### 13. Hoja de Ruta de Implementación 

Esta seccion quedara pendiente, se ira trabajando en la marcha

#### 14. Variables de Entorno

bash

# Servidor
APP_ENV=development
APP_PORT=8080
APP_JWT_SECRET=super-secret-min-32-chars-long
APP_JWT_ACCESS_TTL=15m
APP_JWT_REFRESH_TTL=7d

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=raffle_saas

# OAuth Twitch
TWITCH_CLIENT_ID=xxx
TWITCH_CLIENT_SECRET=xxx
TWITCH_REDIRECT_URI=http://localhost:5173/auth/twitch/callback

# Stripe
STRIPE_SECRET_KEY=sk_test_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx

# Storage (S3/R2)
S3_ENDPOINT=https://xxx.r2.cloudflarestorage.com
S3_ACCESS_KEY_ID=xxx
S3_SECRET_ACCESS_KEY=xxx
S3_BUCKET=raffle-proofs
S3_REGION=auto


📎 Notas para el Equipo de Desarrollo
MongoDB Driver: Usar go.mongodb.org/mongo-driver/v2 (versión 2, no la 1).
Transacciones: Todas las operaciones de compra de tickets DEBEN usar sesiones de transacción.
Sorteo: El PRNG debe ser 100% determinista dada la misma seed. Documentar el algoritmo exacto para auditoría externa.
Seguridad: Nunca exponer password_hash, token_hash, ni stripe_payment_intent_id en responses JSON.
Logging: Todo cambio de estado crítico (ticket verificado, sorteo ejecutado, rifa cancelada) debe generar audit_log.


---

Guarda esto como `START.md` en la raíz de tu repo. Está listo para que tu equipo empiece a codear sin tener que adivinar nada. Si necesitas que profundice en alguna sección específica (por ejemplo, el diagrama de flujo de la transacción atómica de compra o la implementación exacta del PRNG en Go), me avisas.

### Correccion de Developer

Este documento se dividio en 2.
