ALTER TABLE `order_main`
  ADD COLUMN `points_reservation_no` VARCHAR(64) NOT NULL DEFAULT '' AFTER `reservation_no`,
  ADD COLUMN `points_used` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `paid_amount`,
  ADD COLUMN `points_discount_amount` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `points_used`,
  ADD COLUMN `points_rule_snapshot_json` JSON NULL AFTER `points_discount_amount`,
  ADD COLUMN `points_rule_snapshot_digest` VARCHAR(128) NOT NULL DEFAULT '' AFTER `points_rule_snapshot_json`,
  ADD KEY `idx_points_reservation_no` (`points_reservation_no`);

ALTER TABLE `order_sub`
  ADD COLUMN `points_used` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `paid_amount`,
  ADD COLUMN `points_discount_amount` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `points_used`;

CREATE TABLE IF NOT EXISTS `order_points_compensation_task` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_no` VARCHAR(64) NOT NULL,
  `order_no` VARCHAR(64) NOT NULL DEFAULT '',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `points_reservation_no` VARCHAR(64) NOT NULL DEFAULT '',
  `action_code` VARCHAR(64) NOT NULL DEFAULT '',
  `task_status` VARCHAR(32) NOT NULL DEFAULT 'PENDING',
  `retry_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `next_retry_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `last_error` VARCHAR(1000) NOT NULL DEFAULT '',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_no` (`task_no`),
  KEY `idx_status_retry` (`task_status`,`next_retry_at`),
  KEY `idx_order_action` (`order_no`,`action_code`),
  KEY `idx_reservation_action` (`points_reservation_no`,`action_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
