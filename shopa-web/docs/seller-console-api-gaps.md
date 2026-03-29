# Seller Console API Alignment (GoFrame)

## Aligned With Current Backend Contracts

### IAM (`iam-svc`)
- `POST /v1/auth/login/password`
- `POST /v1/auth/token/refresh`
- `POST /v1/auth/logout`

### User Profile (`user-profile-svc`)
- `GET /v1/me/profile`
- `PATCH /v1/me/profile`
- `GET /v1/me/addresses`
- `POST /v1/me/addresses`
- `PATCH /v1/me/addresses/{addressId}`
- `POST /v1/me/addresses/default`

### Seller Aggregation (`edge-gateway`)
- `GET /v1/seller/workbench`
- `GET /v1/seller/shops/{shopNo}/dashboard`

### Catalog (`catalog-svc`)
- `GET /v1/catalog/seller/products`
- `POST /v1/catalog/seller/products/draft`

### Inventory (`inventory-svc`)
- `POST /v1/inventory/sku:batch-get`
- `POST /v1/inventory/seller/stock:adjust`

### Media (`media-svc`)
- `GET /v1/media/biz-assets`
- `POST /v1/media/upload/init`
- `GET /v1/media/assets/{assetId}/process-status`

## Still Missing / Not Yet Fully Productized

- `GET /v1/seller/media/library`
  - No dedicated route in current services; frontend now composes from `GET /v1/media/biz-assets`.
- `GET /v1/seller/media/upload/{uploadId}/status`
  - No dedicated route; frontend maps `uploadId -> assetId` and calls `GET /v1/media/assets/{assetId}/process-status`.
- `GET /v1/seller/inventory/summary`
  - No direct summary route; frontend computes summary from batch stock query results.
- `GET /v1/seller/inventory/stocks`
  - No direct seller stock list route; frontend uses `POST /v1/inventory/sku:batch-get` with known sku list.

## Call-Path Recommendation

- Seller console should call `edge-gateway` first for aggregated read endpoints (`/v1/seller/*`).
- Domain write/read endpoints without gateway aggregation are currently called on domain services directly (`iam-svc`, `user-profile-svc`, `catalog-svc`, `inventory-svc`, `media-svc`).
- If all traffic must go through gateway in later phases, add reverse-proxy routes in `edge-gateway` for:
  - `/v1/me/profile*`
  - `/v1/catalog/seller/*`
  - `/v1/inventory/*`
  - `/v1/media/*`
