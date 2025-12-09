# Project Specification: dysv.de (Go + TanStack Start)

## 1. Project Overview
**Name:** dysv.de
**Architecture:** Headless (Decoupled).
**Repo Strategy:** Polyrepo-style Monorepo.
  - **Root:** Go Backend (built with `ko`).
  - **Web Subdirectory:** TanStack Start Frontend (built with Docker).
**Infrastructure:** Kubernetes (Netcup Control Plane + Hetzner Nodes), Nginx Gateway Fabric.
**Goal:** High-performance, multilingual hosting platform for Next.js/Nuxt developers.

---

## 2. Content & Business Strategy

### A. Pricing Plans & Logic
The app must display these 3 tiers. Toggle between **Monthly** and **Yearly** (Yearly = 2 months free).

| Tier Name | Monthly Price | Target Audience | Hard Limits (Enforced by K8s Quotas) |
| :--- | :--- | :--- | :--- |
| **Static Micro** | **€ 3.90** | React/Vue SPAs | Shared RAM, 1GB Storage |
| **Node Starter** | **€ 9.90** | Personal Blogs | 0.5 vCPU, 512MB RAM, 5GB Storage |
| **Node Pro** | **€ 39.90** | E-commerce/SaaS | 2 vCPU, 4GB RAM, 20GB Storage |

**Upsells:**
* **.de Domain:** +€1.00 / month (Add-on at checkout).

### B. Brand Identity (AI Prompts)
*To be generated using Nano Banana 3 Pro or Veo3.*
* **Aesthetic:** "German Engineering" – Precise, Dark Mode, Electric Blue & Slate Grey.
* **Logo:** Minimalist isometric server block transforming into a shield.
* **Hero Video:** Drone shot flying over a digital map of Germany -> Splitting into two glowing data beams landing in **Nürnberg** and **Falkenstein** -> Rack close-up.
* **Social Banner:** Futuristic datacenter, depth of field, blue LEDs.

### C. Copywriting Key Points
* **Geo-Redundancy:** "Unkillable Uptime. Your site runs across **Nürnberg and Falkenstein** simultaneously."
* **Cost Certainty:** "No metering. No surprise bills. Just flat pricing."
* **Data Sovereignty:** "100% German Infrastructure. ISO 27001 Certified Data Centers."

### D. Copywriting Tone (Strict Guidelines)
* **Do:** Be direct, professional, confident. Use technical terms correctly (Kubernetes, NVMe, SSR, Geo-Redundancy).
* **Don't:** Use marketing fluff ("Skyrocket your business", "Game changer", "Best in class"). Keep it grounded and factual.

### E. Legal & SEO (German Market)
* **Impressum:** Mandatory legal notice.
* **Datenschutz:** Privacy Policy (Must mention Stripe US transfer & Hetzner locations).
* **AGB:** Terms of Service (Specifically covering the 1€ domain cancellation policy).
* **SEO:** Hreflang tags required for `/de`, `/en`, `/hr`.

---

## 3. Backend Specification (Go 1.25)

### A. Core Stack
* **Language:** Go 1.25.
* **Build Tool:** `ko` (No Dockerfile required).
* **CLI:** `spf13/cobra`.
* **Config:** `spf13/viper`.
* **Database:** `go.mongodb.org/mongo-driver/v2`.

### B. Configuration
Load from Environment Variables (injected by K8s Secrets).
```go
type Config struct {
    Port         string `mapstructure:"PORT"`
    MongoURI     string `mapstructure:"MONGODB_URI"`
    StripeSecret string `mapstructure:"STRIPE_SECRET"` // SK_LIVE_...
}
````

### C. Stripe Implementation

The Go backend handles the **Checkout Session Creation**.

  * **Endpoint:** `POST /api/checkout`
  * **Payload:** `{ planId: "node-starter", interval: "yearly", domain: true }`
  * **Response:** `{ sessionId: "cs_test_..." }`
  * **Webhook:** Handle `checkout.session.completed` to provision the K8s namespace.

### D. Build Strategy (`ko`)

**File:** `.ko.yaml`

```yaml
defaultBaseImage: cgr.dev/chainguard/static:latest
builds:
  - id: dysv-api
    main: ./cmd/api
    env:
      - GOFLAGS=-trimpath
      - CGO_ENABLED=0
```

-----

## 4\. Frontend Specification (TanStack Start)

### A. Core Stack

  * **Framework:** TanStack Start (SSR, React 19).
  * **Styling:** Tailwind CSS v4.
  * **State:** TanStack Query & Store.
  * **Forms:** TanStack Form + Zod.
  * **I18n:** `i18next` + `react-i18next`.

### B. I18n Strategy (Path-Based)

  * **Structure:** `app/routes/$locale/index.tsx`
  * **Locales:** `de` (default), `en`, `hr`.
  * **Detection:** Middleware or Edge logic to redirect root `/` to `/$locale` based on `Accept-Language`.

### C. Dockerfile (Multi-Stage)

**Location:** `web/Dockerfile`
**Base:** Node 24 LTS (Alpine).

```dockerfile
# Stage 1: Build
FROM node:24-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# Public Key for Stripe Elements
ENV VITE_STRIPE_PUBLIC_KEY="pk_live_..."
RUN npm run build

# Stage 2: Production
FROM node:24-alpine
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/.output ./.output
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
```

-----

## 5\. Kubernetes & Gateway API

### A. Routing (`dysv-route.yaml`)

Split traffic at the Edge using Nginx Gateway Fabric.

  * `dysv.de/api/*` → Go Backend (Port 8080).
  * `dysv.de/*` → TanStack Frontend (Port 3000).

<!-- end list -->

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: dysv-route
  namespace: dysv
spec:
  parentRefs:
    - name: public
      namespace: nginx-gateway
  hostnames:
    - "dysv.de"
  rules:
    # 1. API Traffic
    - matches:
        - path:
            type: PathPrefix
            value: /api
      backendRefs:
        - name: dysv-api
          port: 8080
    # 2. Web Traffic
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
        - name: dysv-web
          port: 3000
```

-----

## 6\. Development Workflow (Instructions for Agent)

1.  **Backend Dev:**
      * Run: `go run cmd/api/main.go serve`
      * Listens on: `localhost:8080`
2.  **Frontend Dev:**
      * Run: `cd web && npm run dev`
      * Listens on: `localhost:3000`
      * **Proxy Note:** Configure Vite proxy in `web/app.config.ts` to forward `/api` requests to `localhost:8080` to avoid CORS during dev.

## 7\. Definition of Done

1.  **Repo:** Organized as Polyrepo (Root Go, /web Node).
2.  **Images:** Go built with `ko`, Web built with Docker.
3.  **Flow:** User visits `dysv.de/de` -\> Clicks "Kaufen" -\> Go creates Stripe Session -\> User pays -\> Success Page.
4.  **Content:** "Nürnberg & Falkenstein" mentioned in Hero and Features.
5.  **SEO:** `curl -I https://dysv.de/hr` shows correct Hreflang headers.
