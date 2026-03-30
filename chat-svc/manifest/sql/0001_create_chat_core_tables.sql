CREATE TABLE IF NOT EXISTS `conversation` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `conversation_no` VARCHAR(64) NOT NULL,
  `buyer_id` BIGINT UNSIGNED NOT NULL,
  `shop_no` VARCHAR(64) NOT NULL,
  `scene_code` VARCHAR(32) NOT NULL DEFAULT 'PRE_SALE',
  `order_no` VARCHAR(64) NOT NULL DEFAULT '',
  `sub_order_no` VARCHAR(64) NOT NULL DEFAULT '',
  `anchor_spu_no` VARCHAR(64) NOT NULL DEFAULT '',
  `anchor_sku_no` VARCHAR(64) NOT NULL DEFAULT '',
  `last_message_no` VARCHAR(64) NOT NULL DEFAULT '',
  `last_message_preview` VARCHAR(255) NOT NULL DEFAULT '',
  `last_message_at` DATETIME NULL,
  `conversation_status` TINYINT NOT NULL DEFAULT 1,
  `version` BIGINT UNSIGNED NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_conversation_no` (`conversation_no`),
  UNIQUE KEY `uk_scene_anchor` (`buyer_id`,`shop_no`,`scene_code`,`order_no`,`sub_order_no`,`anchor_spu_no`,`anchor_sku_no`),
  KEY `idx_shop_status_time` (`shop_no`,`conversation_status`,`last_message_at`),
  KEY `idx_buyer_status_time` (`buyer_id`,`conversation_status`,`last_message_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话主表';

CREATE TABLE IF NOT EXISTS `chat_message` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `message_no` VARCHAR(64) NOT NULL,
  `conversation_no` VARCHAR(64) NOT NULL,
  `sender_type` TINYINT NOT NULL,
  `sender_user_id` BIGINT UNSIGNED NOT NULL,
  `client_message_no` VARCHAR(64) NULL,
  `message_type` TINYINT NOT NULL,
  `content_text` TEXT NOT NULL,
  `media_asset_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `ext_json` JSON NULL,
  `sent_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_message_no` (`message_no`),
  UNIQUE KEY `uk_client_message` (`conversation_no`,`sender_type`,`sender_user_id`,`client_message_no`),
  KEY `idx_conversation_id` (`conversation_no`,`id`),
  KEY `idx_conversation_sent_at` (`conversation_no`,`sent_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';

CREATE TABLE IF NOT EXISTS `chat_read_offset` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `conversation_no` VARCHAR(64) NOT NULL,
  `reader_type` TINYINT NOT NULL,
  `reader_user_id` BIGINT UNSIGNED NOT NULL,
  `read_to_message_no` VARCHAR(64) NOT NULL DEFAULT '',
  `read_to_message_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `read_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reader_offset` (`conversation_no`,`reader_type`,`reader_user_id`),
  KEY `idx_reader` (`reader_type`,`reader_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话已读游标';

CREATE TABLE IF NOT EXISTS `chat_outbox_event` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id` VARCHAR(64) NOT NULL,
  `event_type` VARCHAR(64) NOT NULL,
  `aggregate_no` VARCHAR(64) NOT NULL,
  `payload_json` JSON NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `retry_count` INT NOT NULL DEFAULT 0,
  `next_retry_at` DATETIME NULL,
  `published_at` DATETIME NULL,
  `error_message` VARCHAR(255) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_status_next_retry` (`status`,`next_retry_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Outbox事件表';

CREATE TABLE IF NOT EXISTS `chat_idempotency` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `biz_code` VARCHAR(64) NOT NULL,
  `idempotency_key` VARCHAR(128) NOT NULL,
  `target_no` VARCHAR(64) NOT NULL DEFAULT '',
  `response_json` JSON NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `expired_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_biz_key` (`user_id`,`biz_code`,`idempotency_key`),
  KEY `idx_expired_at` (`expired_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='幂等记录表';
