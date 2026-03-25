-- ============================================================================
-- iam-svc: outbox lifecycle upgrade
-- Goal:
--   1) keep hot outbox table lean and dispatch-friendly
--   2) support multi-worker safe dispatch with status machine
--   3) provide archive table for long-retention events
-- ============================================================================

ALTER TABLE `iam_outbox_event`
  ADD COLUMN `event_id` CHAR(36) NOT NULL DEFAULT '' COMMENT 'Global unique event id' AFTER `id`,
  ADD COLUMN `available_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT 'Next dispatch schedule time' AFTER `status`,
  ADD COLUMN `sent_at` DATETIME(3) NULL COMMENT 'Successful dispatch time' AFTER `available_at`,
  ADD COLUMN `fail_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Dispatch failure count' AFTER `sent_at`,
  ADD COLUMN `last_error` VARCHAR(512) NOT NULL DEFAULT '' COMMENT 'Last dispatch error text' AFTER `fail_count`;

-- Backfill existing rows so dispatcher can work on historical data safely.
UPDATE `iam_outbox_event`
SET
  `event_id` = IF(`event_id` = '', UUID(), `event_id`),
  `available_at` = IFNULL(`available_at`, `created_at`),
  `fail_count` = IFNULL(`fail_count`, 0)
WHERE `event_id` = '' OR `available_at` IS NULL OR `fail_count` IS NULL;

-- Reuse status column:
-- 1 NEW / 2 PROCESSING / 3 SENT / 4 FAILED / 5 DLQ
ALTER TABLE `iam_outbox_event`
  MODIFY COLUMN `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 NEW,2 PROCESSING,3 SENT,4 FAILED,5 DLQ';

ALTER TABLE `iam_outbox_event`
  ADD UNIQUE KEY `uk_event_id` (`event_id`),
  ADD KEY `idx_status_available_id` (`status`, `available_at`, `id`);

CREATE TABLE IF NOT EXISTS `iam_outbox_event_archive` (
  `id`            BIGINT UNSIGNED NOT NULL,
  `event_id`      CHAR(36) NOT NULL DEFAULT '',
  `event_type`    VARCHAR(64) NOT NULL DEFAULT '',
  `payload_json`  JSON NOT NULL,
  `status`        TINYINT UNSIGNED NOT NULL DEFAULT 3 COMMENT 'Archive rows are usually SENT/FAILED/DLQ snapshots',
  `available_at`  DATETIME(3) NOT NULL,
  `sent_at`       DATETIME(3) NULL,
  `fail_count`    INT UNSIGNED NOT NULL DEFAULT 0,
  `last_error`    VARCHAR(512) NOT NULL DEFAULT '',
  `retry_count`   INT UNSIGNED NOT NULL DEFAULT 0,
  `next_retry_at` DATETIME(3) NULL,
  `created_at`    DATETIME(3) NOT NULL,
  `updated_at`    DATETIME(3) NOT NULL,
  `archived_at`   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_archived_at` (`archived_at`),
  KEY `idx_status_sent_at` (`status`, `sent_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM outbox archive table';
