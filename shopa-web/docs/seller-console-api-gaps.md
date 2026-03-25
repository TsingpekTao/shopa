# Seller Console API Gap List

The following routes are currently consumed with mock fallback in frontend and should be completed in edge-gateway or downstream services.

- `POST /v1/auth/login/password`
- `POST /v1/auth/token/refresh`
- `POST /v1/auth/logout`
- `GET /v1/me/profile`
- `PUT /v1/me/profile`
- `GET /v1/me/addresses`
- `POST /v1/me/addresses`
- `PUT /v1/me/addresses/{address_id}`
- `PUT /v1/me/addresses/{address_id}/default`
- `GET /v1/seller/workbench`
- `GET /v1/seller/shops/{shop_no}/dashboard`
- `GET /v1/seller/media/library`
- `POST /v1/seller/media/upload/init`
- `GET /v1/seller/media/upload/{upload_id}/status`
- `GET /v1/seller/products`
- `POST /v1/seller/products/draft`
- `GET /v1/seller/inventory/stocks`
- `GET /v1/seller/inventory/summary`
- `POST /v1/seller/inventory/adjust`
