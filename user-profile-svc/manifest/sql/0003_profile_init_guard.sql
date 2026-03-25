-- ============================================================================
-- user-profile-svc: consumer dedup table for register-init/event processing
-- ============================================================================

CREATE TABLE IF NOT EXISTS `consumer_event_dedup` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
  `consumer_name`  VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Consumer unique name',
  `event_id`       CHAR(36) NOT NULL DEFAULT '' COMMENT 'Event id for idempotency',
  `event_type`     VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Event type',
  `processed_at`   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT 'Processed timestamp',
  `created_at`     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT 'Created timestamp',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_consumer_event` (`consumer_name`, `event_id`),
  KEY `idx_processed_at` (`processed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Consumer idempotency dedup table';
