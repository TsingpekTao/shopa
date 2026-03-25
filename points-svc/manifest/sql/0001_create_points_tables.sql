-- ============================================================================
-- points-svc: minimal points source-of-truth schema
-- ============================================================================

CREATE TABLE IF NOT EXISTS `points_account` (
  `user_id`     BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
  `balance`     BIGINT NOT NULL DEFAULT 0 COMMENT 'Current points balance',
  `status`      TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 active,2 disabled',
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`),
  KEY `idx_status_updated` (`status`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Points account';

CREATE TABLE IF NOT EXISTS `points_ledger` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`        BIGINT UNSIGNED NOT NULL,
  `biz_type`       VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'REGISTER_INIT/ORDER_PAY/REFUND/etc',
  `biz_id`         VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Business id for idempotency',
  `delta`          BIGINT NOT NULL DEFAULT 0 COMMENT 'Points delta',
  `balance_after`  BIGINT NOT NULL DEFAULT 0 COMMENT 'Balance snapshot after apply',
  `created_at`     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_biz` (`user_id`, `biz_type`, `biz_id`),
  KEY `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Points ledger';
