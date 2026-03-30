USE `shopa_edge_gateway`;
SET NAMES utf8mb4;

INSERT INTO `edge_proxy_route`
(
  `route_code`,
  `method`,
  `path_pattern`,
  `upstream_service`,
  `upstream_path_template`,
  `auth_required`,
  `inject_user_context`,
  `body_mode`,
  `timeout_ms`,
  `rate_limit_rps`,
  `status`,
  `remark`
)
VALUES
  ('EGW_PAYMENT_BUYER_CREATE_INTENT', 'POST', '/v1/payment/buyer/payment-intents:create', 'payment', '/v1/payment/buyer/payment-intents:create', 1, 1, 'NORMAL', 15000, 1000, 1, 'buyer create payment intent'),
  ('EGW_PAYMENT_BUYER_QUERY_INTENT', 'GET', '/v1/payment/buyer/payment-intents:query', 'payment', '/v1/payment/buyer/payment-intents:query', 1, 1, 'NORMAL', 15000, 1000, 1, 'buyer query payment intent'),
  ('EGW_PAYMENT_INTERNAL_CALLBACK', 'POST', '/v1/payment/internal/gateway:callback', 'payment', '/v1/payment/internal/gateway:callback', 1, 1, 'NORMAL', 15000, 1000, 1, 'payment callback'),
  ('EGW_PAYMENT_INTERNAL_CLOSE_INTENT', 'POST', '/v1/payment/internal/payment-intents:close', 'payment', '/v1/payment/internal/payment-intents:close', 1, 1, 'NORMAL', 15000, 1000, 1, 'close payment intent'),
  ('EGW_PAYMENT_INTERNAL_CREATE_REFUND_TASK', 'POST', '/v1/payment/internal/refunds:create-task', 'payment', '/v1/payment/internal/refunds:create-task', 1, 1, 'NORMAL', 15000, 1000, 1, 'create refund task'),
  ('EGW_PAYMENT_INTERNAL_EXECUTE_REFUND_TASK', 'POST', '/v1/payment/internal/refunds:execute-task', 'payment', '/v1/payment/internal/refunds:execute-task', 1, 1, 'NORMAL', 15000, 1000, 1, 'execute refund task'),
  ('EGW_PAYMENT_INTERNAL_RUN_RECON', 'POST', '/v1/payment/internal/reconciliation:run-daily', 'payment', '/v1/payment/internal/reconciliation:run-daily', 1, 1, 'NORMAL', 15000, 1000, 1, 'run daily reconciliation'),
  ('EGW_PAYMENT_INTERNAL_LIST_DIFFS', 'GET', '/v1/payment/internal/reconciliation/diffs', 'payment', '/v1/payment/internal/reconciliation/diffs', 1, 1, 'NORMAL', 15000, 1000, 1, 'list reconciliation diffs'),
  ('EGW_PAYMENT_INTERNAL_RESOLVE_DIFF', 'POST', '/v1/payment/internal/reconciliation:resolve-diff', 'payment', '/v1/payment/internal/reconciliation:resolve-diff', 1, 1, 'NORMAL', 15000, 1000, 1, 'resolve reconciliation diff'),
  ('EGW_PAYMENT_INTERNAL_SNAPSHOT', 'GET', '/v1/payment/internal/payment-intents:snapshot', 'payment', '/v1/payment/internal/payment-intents:snapshot', 1, 1, 'NORMAL', 15000, 1000, 1, 'query payment snapshot'),

  ('EGW_RISK_PRE_CREATE_ORDER', 'POST', '/v1/risk/check:pre-create-order', 'risk', '/v1/risk/check:pre-create-order', 1, 1, 'NORMAL', 15000, 1000, 1, 'risk check before create order'),
  ('EGW_RISK_PRE_PAY', 'POST', '/v1/risk/check:pre-pay', 'risk', '/v1/risk/check:pre-pay', 1, 1, 'NORMAL', 15000, 1000, 1, 'risk check before pay'),
  ('EGW_RISK_PRE_REFUND', 'POST', '/v1/risk/check:pre-refund', 'risk', '/v1/risk/check:pre-refund', 1, 1, 'NORMAL', 15000, 1000, 1, 'risk check before refund'),
  ('EGW_RISK_INTERNAL_INGEST_EVENT', 'POST', '/v1/risk/internal/events:ingest', 'risk', '/v1/risk/internal/events:ingest', 1, 1, 'NORMAL', 15000, 1000, 1, 'ingest risk event'),
  ('EGW_RISK_INTERNAL_REBUILD_FEATURES', 'POST', '/v1/risk/internal/features:rebuild', 'risk', '/v1/risk/internal/features:rebuild', 1, 1, 'NORMAL', 15000, 1000, 1, 'rebuild risk features'),
  ('EGW_RISK_ADMIN_UPSERT_RULE', 'POST', '/v1/risk/admin/rules:upsert', 'risk', '/v1/risk/admin/rules:upsert', 1, 1, 'NORMAL', 15000, 1000, 1, 'upsert risk rule'),
  ('EGW_RISK_ADMIN_ENABLE_RULE', 'POST', '/v1/risk/admin/rules:enable', 'risk', '/v1/risk/admin/rules:enable', 1, 1, 'NORMAL', 15000, 1000, 1, 'enable or disable risk rule'),
  ('EGW_RISK_ADMIN_LIST_HITS', 'GET', '/v1/risk/admin/hits', 'risk', '/v1/risk/admin/hits', 1, 1, 'NORMAL', 15000, 1000, 1, 'list risk hits'),

  ('EGW_NOTIFICATION_SEND', 'POST', '/v1/notification/messages:send', 'notification', '/v1/notification/messages:send', 1, 1, 'NORMAL', 15000, 1000, 1, 'send template message'),
  ('EGW_NOTIFICATION_BATCH_SEND', 'POST', '/v1/notification/messages:batch-send', 'notification', '/v1/notification/messages:batch-send', 1, 1, 'NORMAL', 15000, 1000, 1, 'batch send template message'),
  ('EGW_NOTIFICATION_GET_STATUS', 'GET', '/v1/notification/messages/{notification_no}:status', 'notification', '/v1/notification/messages/{notification_no}:status', 1, 1, 'NORMAL', 15000, 1000, 1, 'query delivery status'),
  ('EGW_NOTIFICATION_GET_PREFERENCE', 'GET', '/v1/notification/me/preference', 'notification', '/v1/notification/me/preference', 1, 1, 'NORMAL', 15000, 1000, 1, 'get notification preference'),
  ('EGW_NOTIFICATION_UPDATE_PREFERENCE', 'POST', '/v1/notification/me/preference:update', 'notification', '/v1/notification/me/preference:update', 1, 1, 'NORMAL', 15000, 1000, 1, 'update notification preference'),
  ('EGW_NOTIFICATION_INTERNAL_RETRY', 'POST', '/v1/notification/internal/messages:retry', 'notification', '/v1/notification/internal/messages:retry', 1, 1, 'NORMAL', 15000, 1000, 1, 'retry failed notification'),

  ('EGW_PROMOTION_LIST_AVAILABLE_COUPONS', 'GET', '/v1/promotion/buyer/coupons:available', 'promotion', '/v1/promotion/buyer/coupons:available', 1, 1, 'NORMAL', 15000, 1000, 1, 'list available coupons'),
  ('EGW_PROMOTION_CALCULATE_DISCOUNT', 'POST', '/v1/promotion/buyer/discount:calculate', 'promotion', '/v1/promotion/buyer/discount:calculate', 1, 1, 'NORMAL', 15000, 1000, 1, 'calculate order discount'),
  ('EGW_PROMOTION_INTERNAL_LOCK_COUPON', 'POST', '/v1/promotion/internal/coupons:lock', 'promotion', '/v1/promotion/internal/coupons:lock', 1, 1, 'NORMAL', 15000, 1000, 1, 'lock coupon'),
  ('EGW_PROMOTION_INTERNAL_CONFIRM_USAGE', 'POST', '/v1/promotion/internal/coupons:confirm', 'promotion', '/v1/promotion/internal/coupons:confirm', 1, 1, 'NORMAL', 15000, 1000, 1, 'confirm coupon usage'),
  ('EGW_PROMOTION_INTERNAL_RELEASE_LOCK', 'POST', '/v1/promotion/internal/coupons:release', 'promotion', '/v1/promotion/internal/coupons:release', 1, 1, 'NORMAL', 15000, 1000, 1, 'release coupon lock'),
  ('EGW_PROMOTION_ADMIN_CREATE_CAMPAIGN', 'POST', '/v1/promotion/admin/campaigns:create', 'promotion', '/v1/promotion/admin/campaigns:create', 1, 1, 'NORMAL', 15000, 1000, 1, 'create campaign'),

  ('EGW_SEARCH_PUBLIC_PRODUCTS', 'GET', '/v1/search/public/products', 'search', '/v1/search/public/products', 0, 0, 'NORMAL', 15000, 1000, 1, 'public product search'),
  ('EGW_SEARCH_PUBLIC_SUGGEST', 'GET', '/v1/search/public/keywords:suggest', 'search', '/v1/search/public/keywords:suggest', 0, 0, 'NORMAL', 15000, 1000, 1, 'public keyword suggest'),
  ('EGW_SEARCH_PUBLIC_BATCH_GET_CARDS', 'POST', '/v1/search/public/spu-cards:batch-get', 'search', '/v1/search/public/spu-cards:batch-get', 0, 0, 'NORMAL', 15000, 1000, 1, 'public batch get spu cards'),
  ('EGW_SEARCH_INTERNAL_UPSERT_DOC', 'POST', '/v1/search/internal/docs:upsert', 'search', '/v1/search/internal/docs:upsert', 1, 1, 'NORMAL', 15000, 1000, 1, 'upsert search doc'),
  ('EGW_SEARCH_INTERNAL_DELETE_DOC', 'POST', '/v1/search/internal/docs:delete', 'search', '/v1/search/internal/docs:delete', 1, 1, 'NORMAL', 15000, 1000, 1, 'delete search doc'),
  ('EGW_SEARCH_INTERNAL_REBUILD_INDEX', 'POST', '/v1/search/internal/index:rebuild', 'search', '/v1/search/internal/index:rebuild', 1, 1, 'NORMAL', 15000, 1000, 1, 'rebuild search index')
ON DUPLICATE KEY UPDATE
  `method` = VALUES(`method`),
  `path_pattern` = VALUES(`path_pattern`),
  `upstream_service` = VALUES(`upstream_service`),
  `upstream_path_template` = VALUES(`upstream_path_template`),
  `auth_required` = VALUES(`auth_required`),
  `inject_user_context` = VALUES(`inject_user_context`),
  `body_mode` = VALUES(`body_mode`),
  `timeout_ms` = VALUES(`timeout_ms`),
  `rate_limit_rps` = VALUES(`rate_limit_rps`),
  `status` = VALUES(`status`),
  `remark` = VALUES(`remark`);
