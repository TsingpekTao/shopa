-- ============================================================================
-- points-svc: V2.2-B upgrade for reservation / bucket / refund / admin flows
-- Execute after 0001_create_points_tables.sql
-- ============================================================================

ALTER TABLE `points_account`
  CHANGE COLUMN `balance` `available_balance` BIGINT NOT NULL DEFAULT 0 COMMENT 'Current available points balance, may be negative when debt exists',
  CHANGE COLUMN `status` `status_code` VARCHAR(32) NOT NULL DEFAULT 'ACTIVE' COMMENT 'ACTIVE/FROZEN/DISABLED',
  ADD COLUMN `frozen_balance` BIGINT NOT NULL DEFAULT 0 COMMENT 'Currently frozen points balance' AFTER `available_balance`,
  ADD COLUMN `total_earned_points` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Lifetime granted points' AFTER `status_code`,
  ADD COLUMN `total_used_points` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Lifetime confirmed used points' AFTER `total_earned_points`,
  ADD COLUMN `total_expired_points` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Lifetime expired points' AFTER `total_used_points`,
  ADD COLUMN `total_adjusted_points` BIGINT NOT NULL DEFAULT 0 COMMENT 'Lifetime manual adjustment points' AFTER `total_expired_points`;

ALTER TABLE `points_ledger`
  ADD COLUMN `ledger_no` VARCHAR(64) NOT NULL COMMENT 'Business ledger identifier' AFTER `id`,
  ADD COLUMN `entry_type_code` VARCHAR(32) NOT NULL DEFAULT 'INIT' COMMENT 'INIT/LOCK/CONFIRM/CANCEL/GRANT/RETURN/REVERSE/EXPIRE/ADJUST/FREEZE/UNFREEZE' AFTER `user_id`,
  CHANGE COLUMN `biz_id` `biz_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Business identifier for idempotency',
  CHANGE COLUMN `delta` `points_delta` BIGINT NOT NULL DEFAULT 0 COMMENT 'Points delta for this ledger entry',
  CHANGE COLUMN `balance_after` `available_after` BIGINT NOT NULL DEFAULT 0 COMMENT 'Available balance snapshot after apply',
  ADD COLUMN `frozen_after` BIGINT NOT NULL DEFAULT 0 COMMENT 'Frozen balance snapshot after apply' AFTER `available_after`,
  ADD COLUMN `debt_after` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Debt snapshot after apply' AFTER `frozen_after`,
  ADD COLUMN `reservation_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Reservation identifier if relevant' AFTER `biz_no`,
  ADD COLUMN `related_bucket_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Bucket identifier if relevant' AFTER `reservation_no`,
  ADD COLUMN `cash_amount_cent` BIGINT NOT NULL DEFAULT 0 COMMENT 'Related cash amount in cents' AFTER `debt_after`,
  ADD COLUMN `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Operator remark or domain explanation' AFTER `cash_amount_cent`,
  ADD COLUMN `extra_json` JSON NULL COMMENT 'Extended metadata snapshot' AFTER `remark`;

ALTER TABLE `points_ledger`
  ADD UNIQUE KEY `uk_ledger_no` (`ledger_no`),
  ADD KEY `idx_user_entry_created` (`user_id`, `entry_type_code`, `created_at`),
  ADD KEY `idx_biz_lookup` (`biz_type`, `biz_no`),
  DROP INDEX `uk_user_biz`;

CREATE TABLE `points_reservation` (
  `reservation_no` VARCHAR(64) NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `order_no` VARCHAR(64) NOT NULL,
  `shop_no` VARCHAR(64) NOT NULL DEFAULT '',
  `reservation_status_code` VARCHAR(32) NOT NULL DEFAULT 'LOCKED' COMMENT 'LOCKED/CONFIRMED/CANCELED/EXPIRED',
  `requested_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `locked_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `locked_cash_amount_cent` BIGINT NOT NULL DEFAULT 0,
  `deduction_digest` VARCHAR(128) NOT NULL DEFAULT '',
  `rule_snapshot_json` JSON NULL,
  `idempotency_key` VARCHAR(128) NOT NULL DEFAULT '',
  `expire_at` DATETIME(3) NULL,
  `confirmed_at` DATETIME(3) NULL,
  `canceled_at` DATETIME(3) NULL,
  `cancel_reason_code` VARCHAR(64) NOT NULL DEFAULT '',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`reservation_no`),
  UNIQUE KEY `uk_order_idempotency` (`order_no`, `idempotency_key`),
  KEY `idx_user_status_created` (`user_id`, `reservation_status_code`, `created_at`),
  KEY `idx_expire_at_status` (`expire_at`, `reservation_status_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Order points reservation';

CREATE TABLE `points_expire_bucket` (
  `bucket_no` VARCHAR(64) NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `source_type` VARCHAR(32) NOT NULL DEFAULT 'GRANT' COMMENT 'GRANT/RETURN/REFUND_GRACE/ADJUST',
  `source_no` VARCHAR(64) NOT NULL DEFAULT '',
  `source_version` BIGINT UNSIGNED NOT NULL DEFAULT 1,
  `bucket_period` VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'YYYYMM or custom window code',
  `total_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `remaining_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `used_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `expired_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `returned_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `expire_at` DATETIME(3) NOT NULL,
  `bucket_status` VARCHAR(32) NOT NULL DEFAULT 'ACTIVE' COMMENT 'ACTIVE/EXPIRED/CLOSED',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`bucket_no`),
  UNIQUE KEY `uk_user_period_source` (`user_id`, `bucket_period`, `source_type`, `source_no`),
  KEY `idx_user_expire_status` (`user_id`, `expire_at`, `bucket_status`),
  KEY `idx_expire_sweep` (`bucket_status`, `expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Points expire buckets, aggregated by period';

CREATE TABLE `points_consumption_detail` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `detail_no` VARCHAR(64) NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `reservation_no` VARCHAR(64) NOT NULL,
  `order_no` VARCHAR(64) NOT NULL,
  `refund_no` VARCHAR(64) NOT NULL DEFAULT '',
  `sub_order_no` VARCHAR(64) NOT NULL DEFAULT '',
  `bucket_no` VARCHAR(64) NOT NULL,
  `consumed_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `returned_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `cash_amount_cent` BIGINT NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_detail_no` (`detail_no`),
  UNIQUE KEY `uk_reservation_bucket_sub_order` (`reservation_no`, `bucket_no`, `sub_order_no`),
  KEY `idx_user_order` (`user_id`, `order_no`),
  KEY `idx_refund_lookup` (`refund_no`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Detailed bucket consumption for later refund return';

CREATE TABLE `points_grant_detail` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `grant_detail_no` VARCHAR(64) NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `order_no` VARCHAR(64) NOT NULL,
  `sub_order_no` VARCHAR(64) NOT NULL,
  `shop_no` VARCHAR(64) NOT NULL DEFAULT '',
  `granted_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `reversed_points` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `rule_code` VARCHAR(64) NOT NULL DEFAULT '',
  `grant_status_code` VARCHAR(32) NOT NULL DEFAULT 'GRANTED' COMMENT 'GRANTED/PARTIAL_REVERSED/REVERSED',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_grant_detail_no` (`grant_detail_no`),
  UNIQUE KEY `uk_order_sub_order` (`order_no`, `sub_order_no`),
  KEY `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Granted points progress for refund reverse';

CREATE TABLE `points_rule_config` (
  `rule_code` VARCHAR(64) NOT NULL,
  `rule_name` VARCHAR(128) NOT NULL,
  `status_code` VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  `min_order_amount_cent` BIGINT NOT NULL DEFAULT 0,
  `max_deduction_rate_bps` INT UNSIGNED NOT NULL DEFAULT 3000 COMMENT 'Max deduction rate in basis points',
  `deduct_points_per_cent` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'How many points are required for one cent discount',
  `grant_points_per_cent` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'How many points are granted for one cent paid amount',
  `refund_grace_days` INT UNSIGNED NOT NULL DEFAULT 7,
  `effective_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `expire_at` DATETIME(3) NULL,
  `rule_snapshot_json` JSON NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`rule_code`),
  KEY `idx_status_effective` (`status_code`, `effective_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Points rule configuration';

CREATE TABLE `points_outbox_event` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_no` VARCHAR(64) NOT NULL,
  `aggregate_type` VARCHAR(32) NOT NULL,
  `aggregate_no` VARCHAR(64) NOT NULL,
  `event_type` VARCHAR(64) NOT NULL,
  `payload_json` JSON NULL,
  `publish_status_code` VARCHAR(32) NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING/SENT/FAILED/DEAD',
  `retry_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `next_retry_at` DATETIME(3) NULL,
  `last_error` VARCHAR(255) NOT NULL DEFAULT '',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_no` (`event_no`),
  KEY `idx_publish_retry` (`publish_status_code`, `next_retry_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Outbox events for points domain';

INSERT INTO `points_rule_config` (
  `rule_code`,
  `rule_name`,
  `status_code`,
  `min_order_amount_cent`,
  `max_deduction_rate_bps`,
  `deduct_points_per_cent`,
  `grant_points_per_cent`,
  `refund_grace_days`,
  `rule_snapshot_json`
) VALUES (
  'DEFAULT_RULE',
  'Default points rule',
  'ACTIVE',
  100,
  3000,
  1,
  1,
  7,
  JSON_OBJECT(
    'rule_code', 'DEFAULT_RULE',
    'min_order_amount_cent', 100,
    'max_deduction_rate_bps', 3000,
    'deduct_points_per_cent', 1,
    'grant_points_per_cent', 1,
    'refund_grace_days', 7
  )
)
ON DUPLICATE KEY UPDATE
  `rule_name` = VALUES(`rule_name`),
  `status_code` = VALUES(`status_code`),
  `min_order_amount_cent` = VALUES(`min_order_amount_cent`),
  `max_deduction_rate_bps` = VALUES(`max_deduction_rate_bps`),
  `deduct_points_per_cent` = VALUES(`deduct_points_per_cent`),
  `grant_points_per_cent` = VALUES(`grant_points_per_cent`),
  `refund_grace_days` = VALUES(`refund_grace_days`),
  `rule_snapshot_json` = VALUES(`rule_snapshot_json`);
