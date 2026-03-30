CREATE TABLE IF NOT EXISTS `cart_item_backup` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `sku_no` VARCHAR(64) NOT NULL,
  `spu_no` VARCHAR(64) NOT NULL DEFAULT '',
  `shop_no` VARCHAR(64) NOT NULL DEFAULT '',
  `qty` INT UNSIGNED NOT NULL DEFAULT 0,
  `checked` TINYINT(1) NOT NULL DEFAULT 0,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1=ACTIVE,2=INVALID',
  `invalid_reason_code` VARCHAR(64) NOT NULL DEFAULT '',
  `spu_title` VARCHAR(255) NOT NULL DEFAULT '',
  `sku_name` VARCHAR(255) NOT NULL DEFAULT '',
  `sku_image_asset_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `sale_price` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `market_price` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `sale_attrs_json` JSON NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_user_sku` (`user_id`, `sku_no`),
  KEY `idx_user_checked_status_updated` (`user_id`, `checked`, `status`, `updated_at`),
  KEY `idx_user_updated` (`user_id`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车持久化备份表';

CREATE TABLE IF NOT EXISTS `cart_sync_checkpoint` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `last_synced_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_sync_version` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_user` (`user_id`),
  KEY `idx_last_synced_at` (`last_synced_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车同步检查点';

CREATE TABLE IF NOT EXISTS `cart_outbox_event` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id` VARCHAR(64) NOT NULL,
  `event_type` VARCHAR(64) NOT NULL,
  `aggregate_type` VARCHAR(64) NOT NULL DEFAULT 'CART',
  `aggregate_id` VARCHAR(64) NOT NULL DEFAULT '',
  `payload_json` JSON NOT NULL,
  `headers_json` JSON NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=NEW,1=PROCESSING,2=SENT,3=FAILED,4=DLQ',
  `available_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `retry_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `last_error` VARCHAR(1024) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_event_id` (`event_id`),
  KEY `idx_status_available_id` (`status`, `available_at`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车 outbox 事件表';
