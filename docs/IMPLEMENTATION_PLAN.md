# dysv.de MVP Implementation Plan

> High-performance, multilingual hosting platform for Next.js/Nuxt developers.

---

## Current State

| Layer | Status | Notes |
|-------|--------|-------|
| Go Backend | Scaffolded | Cobra CLI only, no handlers/services/repos |
| Frontend | Scaffolded | TanStack Start with demo routes, no product pages |
| Infrastructure | Not started | K8s manifests, Gateway API routing TBD |

---

## Phase 1: Frontend — Landing & Pricing

### Goal
Create the public-facing marketing site with hero, pricing, and cart functionality.

---

#### [NEW] [pricing-data.ts](file:///home/darko/go/src/github.com/deicod/dysv/web/src/data/pricing-data.ts)

Static pricing data for the 3 tiers and `.de` domain upsell:

```typescript
export const plans = [
  { id: 'static-micro', name: 'Static Micro', monthlyPrice: 3.90, limits: '1GB Storage' },
  { id: 'node-starter', name: 'Node Starter', monthlyPrice: 9.90, limits: '0.5 vCPU, 512MB RAM, 5GB' },
  { id: 'node-pro', name: 'Node Pro', monthlyPrice: 39.90, limits: '2 vCPU, 4GB RAM, 20GB' },
]
export const domainAddon = { id: 'de-domain', name: '.de Domain', monthlyPrice: 1.00 }
```

---

#### [MODIFY] [index.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/routes/index.tsx)

Replace TanStack demo content with dysv.de landing page:

- **Hero Section**: "Serverstandort Deutschland", tagline about flat pricing & geo-redundancy
- **Features Grid**: NVMe storage, K8s-powered, ISO 27001, no metering
- **Call-to-Action**: Link to `/pricing`

---

#### [NEW] [pricing.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/routes/pricing.tsx)

Pricing page with:

- Monthly/Yearly toggle (Yearly = 2 months free)
- 3 pricing cards from `pricing-data.ts`
- "Add to Cart" buttons using TanStack Store
- `.de Domain` upsell checkbox

---

#### [NEW] [cart.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/routes/cart.tsx)

Cart review page:

- Display selected plan and add-ons
- Show subtotal with Monthly/Yearly breakdown
- "Proceed to Checkout" button (→ Stripe in Phase 3)

---

#### [NEW] [cart-store.ts](file:///home/darko/go/src/github.com/deicod/dysv/web/src/lib/cart-store.ts)

TanStack Store for cart state:

```typescript
interface CartState {
  planId: string | null
  billingCycle: 'monthly' | 'yearly'
  addons: string[]
}
```

---

#### [NEW] [PricingCard.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/components/PricingCard.tsx)

Reusable pricing card component with glassmorphism styling.

---

#### [MODIFY] [Header.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/components/Header.tsx)

Update navigation to include: Home, Pricing, Cart (with item count badge).

---

## Phase 2: Backend — Cart & Order API

### Goal
Implement Go API with `handler → service → repo` pattern.

---

#### [NEW] [cmd/api/main.go](file:///home/darko/go/src/github.com/deicod/dysv/cmd/api/main.go)

Cobra subcommand `dysv api` to start HTTP server on port 8080.

---

#### [NEW] [internal/config/config.go](file:///home/darko/go/src/github.com/deicod/dysv/internal/config/config.go)

Viper-based config loading:

```go
type Config struct {
    Port         string `mapstructure:"PORT"`
    MongoURI     string `mapstructure:"MONGODB_URI"`
    StripeSecret string `mapstructure:"STRIPE_SECRET"`
}
```

---

#### [NEW] [internal/model/](file:///home/darko/go/src/github.com/deicod/dysv/internal/model/)

Domain models: `Cart`, `Order`, `LineItem`.

---

#### [NEW] [internal/repo/](file:///home/darko/go/src/github.com/deicod/dysv/internal/repo/)

Repository layer with MongoDB driver:

- `cart_repo.go`: `FindOne`, `ReplaceOne`
- `order_repo.go`: `InsertOne`

---

#### [NEW] [internal/service/](file:///home/darko/go/src/github.com/deicod/dysv/internal/service/)

Business logic:

- `cart_service.go`: `AddItemToCart`, `RemoveItem`, `GetCartTotal`
- `order_service.go`: `CreateOrderFromCart`

---

#### [NEW] [internal/handler/](file:///home/darko/go/src/github.com/deicod/dysv/internal/handler/)

HTTP handlers (JSON parsing, validation, response):

- `POST /api/cart` — Add item
- `GET /api/cart` — Get cart
- `DELETE /api/cart/:itemId` — Remove item
- `POST /api/checkout` — Create Stripe session

---

## Phase 3: Stripe Integration

### Goal
Connect checkout flow to Stripe.

---

#### [MODIFY] [internal/handler/checkout.go](file:///home/darko/go/src/github.com/deicod/dysv/internal/handler/checkout.go)

- Create Stripe Checkout Session with line items
- Return session URL for redirect

---

#### [MODIFY] [cart.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/routes/cart.tsx)

- On "Checkout" click, call backend, redirect to Stripe URL

---

## Phase 4: Legal & SEO

### Goal
German market compliance.

---

#### [NEW] [impressum.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/routes/impressum.tsx)
#### [NEW] [datenschutz.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/routes/datenschutz.tsx)
#### [NEW] [agb.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/routes/agb.tsx)

Static legal pages (content placeholder, user to provide text).

---

#### [MODIFY] [__root.tsx](file:///home/darko/go/src/github.com/deicod/dysv/web/src/routes/__root.tsx)

Add hreflang meta tags for `/de`, `/en`, `/hr`.

---

## Phase 5: Infrastructure

### Goal
Kubernetes deployment with Gateway API routing.

---

#### [NEW] [k8s/dysv-api.yaml](file:///home/darko/go/src/github.com/deicod/dysv/k8s/dysv-api.yaml)

Deployment + Service for Go backend.

---

#### [NEW] [k8s/dysv-web.yaml](file:///home/darko/go/src/github.com/deicod/dysv/k8s/dysv-web.yaml)

Deployment + Service for TanStack frontend.

---

#### [NEW] [k8s/dysv-route.yaml](file:///home/darko/go/src/github.com/deicod/dysv/k8s/dysv-route.yaml)

HTTPRoute splitting `/api/*` → backend, `/*` → frontend.

---

#### [NEW] [.ko.yaml](file:///home/darko/go/src/github.com/deicod/dysv/.ko.yaml)

Ko build config for Go backend.

---

#### [NEW] [web/Dockerfile](file:///home/darko/go/src/github.com/deicod/dysv/web/Dockerfile)

Multi-stage Node 24 Dockerfile for frontend.

---

## Verification Plan

### Automated

| Test | Command |
|------|---------|
| Frontend builds | `cd web && npm run build` |
| Go compiles | `go build ./...` |
| Frontend tests | `cd web && npm test` |
| Biome lint | `cd web && npm run check` |

### Manual (User)

1. **Pricing Page**: Verify 3 tiers display correctly, Monthly/Yearly toggle works
2. **Cart Flow**: Add plan → View cart → Totals correct
3. **SEO Headers**: `curl -I https://dysv.de/hr` returns hreflang
4. **Stripe**: Test checkout in Stripe test mode

---

## Definition of Done

- [ ] Go code follows `handler → service → repo` pattern
- [ ] Shopping flow: Add to Cart → View Cart → Checkout
- [ ] Go built with `ko`, Web built with Docker
- [ ] Hero says "**Serverstandort Deutschland**"
- [ ] Hreflang headers for `/de`, `/en`, `/hr`
