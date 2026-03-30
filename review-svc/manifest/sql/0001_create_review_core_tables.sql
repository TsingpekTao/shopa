SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `review_record` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `review_no` VARCHAR(64) NOT NULL,
  `order_no` VARCHAR(64) NOT NULL,
  `sub_order_no` VARCHAR(64) NOT NULL,
  `item_no` VARCHAR(64) NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `shop_no` VARCHAR(64) NOT NULL,
  `spu_no` VARCHAR(64) NOT NULL,
  `sku_no` VARCHAR(64) NOT NULL,
  `score` TINYINT UNSIGNED NOT NULL,
  `content` TEXT NOT NULL,
  `medias_json` JSON NULL,
  `anonymous` TINYINT(1) NOT NULL DEFAULT 0,
  `append_content` TEXT NULL,
  `append_medias_json` JSON NULL,
  `append_at` DATETIME NULL,
  `seller_reply` TEXT NULL,
  `seller_reply_at` DATETIME NULL,
  `review_status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `like_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `version` BIGINT UNSIGNED NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_review_no` (`review_no`),
  UNIQUE KEY `uk_order_item` (`order_no`, `sub_order_no`, `item_no`),
  KEY `idx_user_id_id` (`user_id`, `id`),
  KEY `idx_spu_status_id` (`spu_no`, `review_status`, `id`),
  KEY `idx_shop_status_id` (`shop_no`, `review_status`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `review_spu_summary` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `spu_no` VARCHAR(64) NOT NULL,
  `total_reviews` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `score_1_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `score_2_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `score_3_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `score_4_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `score_5_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `avg_score_x100` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `positive_rate_x100` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_spu_no` (`spu_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `review_operate_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `review_no` VARCHAR(64) NOT NULL,
  `operator_type` VARCHAR(32) NOT NULL,
  `operator_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `action_code` VARCHAR(64) NOT NULL,
  `reason_code` VARCHAR(64) NOT NULL DEFAULT '',
  `remark` VARCHAR(500) NOT NULL DEFAULT '',
  `snapshot_json` JSON NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_review_no_created_at` (`review_no`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `review_outbox_event` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id` VARCHAR(64) NOT NULL,
  `aggregate_type` VARCHAR(64) NOT NULL,
  `aggregate_id` VARCHAR(64) NOT NULL,
  `event_type` VARCHAR(128) NOT NULL,
  `payload_json` JSON NOT NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `retry_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `next_retry_at` DATETIME NULL,
  `published_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_status_next_retry_id` (`status`, `next_retry_at`, `id`),
  KEY `idx_aggregate` (`aggregate_type`, `aggregate_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `review_idempotency` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `action` VARCHAR(64) NOT NULL,
  `idempotency_key` VARCHAR(128) NOT NULL,
  `resource_no` VARCHAR(64) NOT NULL DEFAULT '',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `response_json` JSON NULL,
  `error_code` VARCHAR(64) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_action_key` (`user_id`, `action`, `idempotency_key`),
  KEY `idx_resource_no` (`resource_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
