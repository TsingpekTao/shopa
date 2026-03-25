# Shopa Web

Monorepo frontend workspace for Shopa seller-side UI.

## Run

```bash
npm install
npm run dev:seller
```

## Structure

- `apps/seller-console`: seller BFF frontend app.
- `packages/api-client`: HTTP client with auth/degraded interceptors.
- `packages/types`: shared TypeScript contracts.
- `packages/ui`: shared UI shell and status components.
