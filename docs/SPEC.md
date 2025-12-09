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
* **Hero Video:** Drone shot flying over a digital map of Germany -> Landing on a glowing data center hub -> Rack close-up.
* **Social Banner:** Futuristic datacenter, depth of field, blue LEDs.

### C. Copywriting Key Points
* **Geo-Redundancy:** "Unkillable Uptime. Your site runs across multiple physical availability zones."
* **Cost Certainty:** "No metering. No surprise bills. Just flat pricing."
* **Data Sovereignty:** "**Serverstandort Deutschland**. ISO 27001 Certified."

### D. Copywriting Tone (Strict Guidelines)
* **Do:** Be direct, professional, confident. Use technical terms correctly (Kubernetes, NVMe, SSR).
* **Don't:** Use marketing fluff or mention specific provider names (Hetzner/Netcup). Focus on the value: "German Location".

### E. Legal & SEO (German Market)
* **Impressum:** Mandatory legal notice.
* **Datenschutz:** Privacy Policy.
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

### B. Architecture: Service/Repository Pattern
**Strict Rule:** Handlers are thin. Services hold logic. Repos wrap DB calls.

**Directory Structure:**
```text
internal/
├── handler/       # HTTP/JSON parsing, Validation, calls Service
├── service/       # Business Logic (Calculations, Stripe calls)
└── repo/          # Database Access (Mongo Driver wrappers)
````

**1. Repository Layer (`internal/repo`)**
Only implement methods that are actually used.

  * `repo.CartRepo`: `FindOne`, `ReplaceOne`
  * `repo.OrderRepo`: `InsertOne`
  * *Signature Example:* `func (r *CartRepo) FindOne(ctx context.Context, filter interface{}) (*model.Cart, error)`

**2. Service Layer (`internal/service`)**
Holds the use cases.

  * `service.CartService`: `AddItemToCart`, `RemoveItem`, `GetCartTotal`
  * `service.OrderService`: `CreateOrderFromCart`
  * *Example:* `AddItemToCart` calls `repo.CartRepo.FindOne`, updates struct, calls `repo.CartRepo.ReplaceOne`.

**3. Handler Layer (`internal/handler`)**

  * Parses JSON body.
  * Calls `service.CartService.AddItemToCart`.
  * Returns JSON response.

### C. Configuration

Load from Environment Variables (injected by K8s Secrets).

```go
type Config struct {
    Port         string `mapstructure:"PORT"`
    MongoURI     string `mapstructure:"MONGODB_URI"`
    StripeSecret string `mapstructure:"STRIPE_SECRET"`
}
```

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

### B. User Flow

1.  **Pricing Page:** User clicks "Add to Cart" on a plan.
2.  **Cart Drawer/Page:** User reviews selection, toggles "Monthly/Yearly", adds ".de Domain".
3.  **Checkout:** User clicks "Checkout" -\> Redirects to Stripe.

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

## 6\. Definition of Done

1.  **Architecture:** Go code follows `handler -> service -> repo` pattern cleanly.
2.  **Shopping Flow:** Functional "Add to Cart" -\> "View Cart" -\> "Checkout" flow.
3.  **Images:** Go built with `ko`, Web built with Docker.
4.  **Content:** Hero section says "**Serverstandort Deutschland**" (No Hetzner mentions).
5.  **SEO:** `curl -I https://dysv.de/hr` shows correct Hreflang headers.
