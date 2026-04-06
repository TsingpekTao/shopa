-- ============================================================================
-- media-svc: default scene policy seeds
-- keep this script idempotent so local/dev databases can be repaired safely
-- ============================================================================

INSERT INTO `media_scene_policy` (
  `scene_code`,
  `acl_type`,
  `max_count`,
  `max_size_bytes`,
  `allow_mime_json`,
  `allow_ext_json`,
  `retention_days`,
  `gc_grace_hours`,
  `risk_level`,
  `risk_async_enabled`,
  `process_async_enabled`,
  `status`,
  `remark`,
  `ext_json`
) VALUES
(
  'seller_onboarding_cert',
  2,
  3,
  10485760,
  JSON_ARRAY('image/jpeg', 'image/png', 'image/webp'),
  JSON_ARRAY('.jpg', '.jpeg', '.png', '.webp'),
  30,
  24,
  2,
  1,
  0,
  1,
  'seller onboarding cert images',
  NULL
),
(
  'seller_media',
  1,
  20,
  10485760,
  JSON_ARRAY('image/jpeg', 'image/png', 'image/webp'),
  JSON_ARRAY('.jpg', '.jpeg', '.png', '.webp'),
  180,
  24,
  1,
  1,
  0,
  1,
  'seller product publish images',
  JSON_OBJECT('biz_types', JSON_ARRAY('seller_product_publish'))
)
ON DUPLICATE KEY UPDATE
  `acl_type` = VALUES(`acl_type`),
  `max_count` = VALUES(`max_count`),
  `max_size_bytes` = VALUES(`max_size_bytes`),
  `allow_mime_json` = VALUES(`allow_mime_json`),
  `allow_ext_json` = VALUES(`allow_ext_json`),
  `retention_days` = VALUES(`retention_days`),
  `gc_grace_hours` = VALUES(`gc_grace_hours`),
  `risk_level` = VALUES(`risk_level`),
  `risk_async_enabled` = VALUES(`risk_async_enabled`),
  `process_async_enabled` = VALUES(`process_async_enabled`),
  `status` = VALUES(`status`),
  `remark` = VALUES(`remark`),
  `ext_json` = VALUES(`ext_json`);
