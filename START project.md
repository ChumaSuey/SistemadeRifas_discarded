# 🎟️ START.md — Sistema de Rifas (SaaS Comercial)
# Versión: 1.0 | Estado: Diseño Técnico Completo

---

## 📑 Índice

1. [Visión General](#1-visión-general)
2. [Stack Tecnológico](#2-stack-tecnológico)
3. [Arquitectura de Carpetas](#3-arquitectura-de-carpetas)
4. [Roles del Sistema](#4-roles-del-sistema)
5. [Esquema de Base de Datos (MongoDB)](#5-esquema-de-base-de-datos-mongodb)
6. [Modelos de Datos en Go](#6-modelos-de-datos-en-go)
7. [Índices de Base de Datos](#7-índices-de-base-de-datos)
8. [Flujo de Pagos](#8-flujo-de-pagos)
9. [Algoritmo de Sorteo Transparente](#9-algoritmo-de-sorteo-transparente)
10. [API REST — Endpoints](#10-api-rest--endpoints)
11. [Reglas de Negocio Core](#11-reglas-de-negocio-core)
12. [Modelo de Monetización](#12-modelo-de-monetización)
13. [Hoja de Ruta de Implementación](#13-hoja-de-ruta-de-implementación)
14. [Variables de Entorno](#14-variables-de-entorno)

---

## 1. Visión General

Plataforma SaaS para creadores de contenido (Twitch, redes sociales) que permite crear, gestionar y ejecutar rifas con flujo de compra, verificación de pagos/requisitos y selección de ganadores en directo de forma auditada y transparente.

### 🎯 Características Clave
- **Transparencia Auditada:** Algoritmo de sorteo con semilla hash pública verificable externamente.
- **Verificación de Pagos:** Panel de trabajo para verificadores globales que validan comprobantes de pago manual o transacciones automáticas.
- **Pagos Dual:** Soporte para pagos manuales (transferencia + screenshot) y pagos automáticos (Stripe/PayPal).
- **Experiencia en Vivo:** Dashboard optimizado para retransmisiones en directo con sorteo rápido y re-sorteo instantáneo.
- **Multi-país:** Moneda de referencia VES (Venezuela) con soporte para USD y otras divisas.

---

## 2. Stack Tecnológico

| Capa | Tecnología |
|---|---|
| **Frontend** | React 18+ + TypeScript + Vite |
| **Gestor de Paquetes** | `pnpm` |
| **Backend** | Go 1.22+ — API REST |
| **Base de Datos** | MongoDB 6.0+ |
| **Driver MongoDB** | `go.mongodb.org/mongo-driver/v2` |
| **Auth** | JWT (access + refresh tokens) + OAuth 2.0 (Twitch) |
| **Storage de Archivos** | AWS S3 / Cloudflare R2 |
| **Pagos Automáticos** | Stripe |
| **Pagos Manuales** | Transferencia bancaria + upload de comprobante |
| **Containerización** | Docker + Docker Compose |
| **Logging** | `slog` (estándar Go) |

---

## 3. Arquitectura de Carpetas (Backend Go)
raffle-saas/
├── cmd/
│   └── api/
│       └── main.go                 # Punto de entrada
├── internal/
│   ├── config/                     # Variables de entorno
│   ├── domain/                     # Entidades de negocio (structs puros)
│   ├── repository/                 # Capa de acceso a MongoDB
│   ├── service/                    # Lógica de negocio
│   ├── handler/                    # HTTP handlers (controllers)
│   ├── middleware/                 # Auth, CORS, rate limiting, logging
│   ├── dto/                        # Request/Response structs
│   ├── validator/                  # Validaciones custom
│   └── utils/                      # Helpers (hash, jwt, crypto)
├── pkg/
│   └── crypto/                     # PRNG verificable, hash utilities
├── migrations/                     # Scripts de setup de índices MongoDB
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum


---

## 4. Roles del Sistema

| Rol | Descripción | Alcance |
|---|---|---|
| **Admin** | Supervisión global, métricas, resolución de disputas, gestión de verifiers. | Plataforma |
| **Raffler** | Crea rifas, configura parámetros, ejecuta sorteos, designa verificadores. | Sus rifas |
| **Participant** | Compra/reclama boletos, sube comprobantes de pago. | Sus tickets |
| **Verifier** | Revisa y aprueba/rechaza comprobantes de pago manuales. Validan que la rifa se ejecutó correctamente. | Global (todas las rifas) |

> **Nota:** Los Verifiers son usuarios con rol global. Un Admin puede promover a cualquier usuario a Verifier. No son asignados por rifa.

---

## 5. Esquema de Base de Datos (MongoDB)

### 5.1. `users`
```json
{
  "_id": "ObjectID",
  "name": "string",
  "email": "string (unique)",
  "email_verified": "boolean",
  "password_hash": "string | null",
  "oauth_provider": "string | null",      // "twitch", "google"
  "oauth_id": "string | null",
  "role": "admin | raffler | participant | verifier",
  "status": "active | suspended | banned",
  "social_links": {
    "twitch": "string",
    "instagram": "string"
  },
  "created_at": "ISODate",
  "updated_at": "ISODate"
}


### 5.2 'raffles'
```json
{
  "_id": "ObjectID",
  "raffler_id": "ObjectID",
  "title": "string",
  "description": "string",
  "status": "draft | pending_verification | active | finished | cancelled",
  "pricing_tier": "free | pro",
  "ticket_price": 0.00,
  "currency": "USD | VES | EUR",
  "payment_methods": ["manual_transfer", "stripe"],
  "tickets_total": 500,
  "tickets_sold": 0,
  "winner_id": "ObjectID | null",
  "winning_ticket_id": "ObjectID | null",
  "sort_seed": "string | null",
  "sort_commit_hash": "string | null",
  "draw_date": "ISODate | null",
  "created_at": "ISODate",
  "updated_at": "ISODate"
}

START pagos dash
Aquí tienes el START.md completo y actualizado. Está diseñado para que tu equipo lo lea y empiece a codear sin ambigüedades.
Markdown
Copiar
Código
Vista previa
# 🎟️ START.md — Sistema de Rifas (SaaS Comercial)
# Versión: 1.0 | Estado: Diseño Técnico Completo

---

## 📑 Índice

1. [Visión General](#1-visión-general)
2. [Stack Tecnológico](#2-stack-tecnológico)
3. [Arquitectura de Carpetas](#3-arquitectura-de-carpetas)
4. [Roles del Sistema](#4-roles-del-sistema)
5. [Esquema de Base de Datos (MongoDB)](#5-esquema-de-base-de-datos-mongodb)
6. [Modelos de Datos en Go](#6-modelos-de-datos-en-go)
7. [Índices de Base de Datos](#7-índices-de-base-de-datos)
8. [Flujo de Pagos](#8-flujo-de-pagos)
9. [Algoritmo de Sorteo Transparente](#9-algoritmo-de-sorteo-transparente)
10. [API REST — Endpoints](#10-api-rest--endpoints)
11. [Reglas de Negocio Core](#11-reglas-de-negocio-core)
12. [Modelo de Monetización](#12-modelo-de-monetización)
13. [Hoja de Ruta de Implementación](#13-hoja-de-ruta-de-implementación)
14. [Variables de Entorno](#14-variables-de-entorno)

---

## 1. Visión General

Plataforma SaaS para creadores de contenido (Twitch, redes sociales) que permite crear, gestionar y ejecutar rifas con flujo de compra, verificación de pagos/requisitos y selección de ganadores en directo de forma auditada y transparente.

### 🎯 Características Clave
- **Transparencia Auditada:** Algoritmo de sorteo con semilla hash pública verificable externamente.
- **Verificación de Pagos:** Panel de trabajo para verificadores globales que validan comprobantes de pago manual o transacciones automáticas.
- **Pagos Dual:** Soporte para pagos manuales (transferencia + screenshot) y pagos automáticos (Stripe/PayPal).
- **Experiencia en Vivo:** Dashboard optimizado para retransmisiones en directo con sorteo rápido y re-sorteo instantáneo.
- **Multi-país:** Moneda de referencia VES (Venezuela) con soporte para USD y otras divisas.

---

## 2. Stack Tecnológico

| Capa | Tecnología |
|---|---|
| **Frontend** | React 18+ + TypeScript + Vite |
| **Gestor de Paquetes** | `pnpm` |
| **Backend** | Go 1.22+ — API REST |
| **Base de Datos** | MongoDB 6.0+ |
| **Driver MongoDB** | `go.mongodb.org/mongo-driver/v2` |
| **Auth** | JWT (access + refresh tokens) + OAuth 2.0 (Twitch) |
| **Storage de Archivos** | AWS S3 / Cloudflare R2 |
| **Pagos Automáticos** | Stripe |
| **Pagos Manuales** | Transferencia bancaria + upload de comprobante |
| **Containerización** | Docker + Docker Compose |
| **Logging** | `slog` (estándar Go) |

---

## 3. Arquitectura de Carpetas (Backend Go)
raffle-saas/
├── cmd/
│   └── api/
│       └── main.go                 # Punto de entrada
├── internal/
│   ├── config/                     # Variables de entorno
│   ├── domain/                     # Entidades de negocio (structs puros)
│   ├── repository/                 # Capa de acceso a MongoDB
│   ├── service/                    # Lógica de negocio
│   ├── handler/                    # HTTP handlers (controllers)
│   ├── middleware/                 # Auth, CORS, rate limiting, logging
│   ├── dto/                        # Request/Response structs
│   ├── validator/                  # Validaciones custom
│   └── utils/                      # Helpers (hash, jwt, crypto)
├── pkg/
│   └── crypto/                     # PRNG verificable, hash utilities
├── migrations/                     # Scripts de setup de índices MongoDB
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
plain

---

## 4. Roles del Sistema

| Rol | Descripción | Alcance |
|---|---|---|
| **Admin** | Supervisión global, métricas, resolución de disputas, gestión de verifiers. | Plataforma |
| **Raffler** | Crea rifas, configura parámetros, ejecuta sorteos, designa verificadores. | Sus rifas |
| **Participant** | Compra/reclama boletos, sube comprobantes de pago. | Sus tickets |
| **Verifier** | Revisa y aprueba/rechaza comprobantes de pago manuales. Validan que la rifa se ejecutó correctamente. | Global (todas las rifas) |

> **Nota:** Los Verifiers son usuarios con rol global. Un Admin puede promover a cualquier usuario a Verifier. No son asignados por rifa.

---

## 5. Esquema de Base de Datos (MongoDB)

### 5.1. `users`
```json
{
  "_id": "ObjectID",
  "name": "string",
  "email": "string (unique)",
  "email_verified": "boolean",
  "password_hash": "string | null",
  "oauth_provider": "string | null",      // "twitch", "google"
  "oauth_id": "string | null",
  "role": "admin | raffler | participant | verifier",
  "status": "active | suspended | banned",
  "social_links": {
    "twitch": "string",
    "instagram": "string"
  },
  "created_at": "ISODate",
  "updated_at": "ISODate"
}
5.2. raffles
JSON
{
  "_id": "ObjectID",
  "raffler_id": "ObjectID",
  "title": "string",
  "description": "string",
  "status": "draft | pending_verification | active | finished | cancelled",
  "pricing_tier": "free | pro",
  "ticket_price": 0.00,
  "currency": "USD | VES | EUR",
  "payment_methods": ["manual_transfer", "stripe"],
  "tickets_total": 500,
  "tickets_sold": 0,
  "winner_id": "ObjectID | null",
  "winning_ticket_id": "ObjectID | null",
  "sort_seed": "string | null",
  "sort_commit_hash": "string | null",
  "draw_date": "ISODate | null",
  "created_at": "ISODate",
  "updated_at": "ISODate"
}

### 5.3. tickets
JSON
{
  "_id": "ObjectID",
  "raffle_id": "ObjectID",
  "participant_id": "ObjectID",
  "ticket_number": "int",
  "status": "pending | verified | rejected | cancelled | winner",
  "payment_method": "manual_transfer | stripe",
  "payment_proof_url": "string | null",
  "stripe_payment_intent_id": "string | null",
  "rejection_reason": "string | null",
  "verified_by": "ObjectID | null",
  "verified_at": "ISODate | null",
  "created_at": "ISODate"
}

### 5.4. audit_logs
JSON
{
  "_id": "ObjectID",
  "raffle_id": "ObjectID",
  "user_id": "ObjectID",
  "action": "string",
  "details": {},
  "ip_address": "string",
  "user_agent": "string",
  "created_at": "ISODate"
}

### 5.5. refresh_tokens
JSON
{
  "_id": "ObjectID",
  "user_id": "ObjectID",
  "token_hash": "string",
  "expires_at": "ISODate",
  "created_at": "ISODate"
}

