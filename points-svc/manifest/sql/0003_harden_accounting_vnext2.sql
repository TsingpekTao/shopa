-- ============================================================================
-- points-svc: VNext.2 accounting hardening for locked buckets and reconciliation
-- Execute after 0002_upgrade_points_v22b.sql
-- ============================================================================

ALTER TABLE `points_expire_bucket`
  ADD COLUMN `locked_points` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Currently locked points reserved by active reservations' AFTER `remaining_points`;

ALTER TABLE `points_consumption_detail`
  ADD COLUMN `detail_status_code` VARCHAR(32) NOT NULL DEFAULT 'LOCKED' COMMENT 'LOCKED/CONFIRMED/CANCELED' AFTER `bucket_no`;

ALTER TABLE `points_grant_detail`
  ADD COLUMN `rule_snapshot_json` JSON NULL COMMENT 'Grant rule snapshot used when points were granted' AFTER `rule_code`;

CREATE TABLE IF NOT EXISTS `points_refund_action` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `action_no` VARCHAR(64) NOT NULL,
  `action_type` VARCHAR(32) NOT NULL COMMENT 'RETURN/REVERSE',
  `refund_no` VARCHAR(64) NOT NULL,
  `order_no` VARCHAR(64) NOT NULL,
  `sub_order_no` VARCHAR(64) NOT NULL DEFAULT '',
  `shop_no` VARCHAR(64) NOT NULL DEFAULT '',
  `user_id` BIGINT UNSIGNED NOT NULL,
  `idempotency_key` VARCHAR(128) NOT NULL,
  `requested_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `effective_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `cash_amount_cent` BIGINT NOT NULL DEFAULT 0,
  `action_status` VARCHAR(32) NOT NULL DEFAULT 'SUCCESS' COMMENT 'SUCCESS',
  `result_payload_json` JSON NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_action_no` (`action_no`),
  UNIQUE KEY `uk_action_dedupe` (`action_type`, `refund_no`, `sub_order_no`, `idempotency_key`),
  KEY `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Refund action idempotency and replay store';
