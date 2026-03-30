SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `fulfillment_shipment` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `shipment_no` VARCHAR(64) NOT NULL,
  `order_no` VARCHAR(64) NOT NULL,
  `sub_order_no` VARCHAR(64) NOT NULL,
  `shop_no` VARCHAR(64) NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `logistics_company_code` VARCHAR(64) NOT NULL DEFAULT '',
  `logistics_company_name` VARCHAR(128) NOT NULL DEFAULT '',
  `logistics_no` VARCHAR(128) NOT NULL DEFAULT '',
  `shipment_status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 WAIT_SHIP 2 SHIPPED 3 IN_TRANSIT 4 DELIVERED 5 EXCEPTION 6 CLOSED',
  `shipped_at` DATETIME NULL,
  `delivered_at` DATETIME NULL,
  `receiver_name` VARCHAR(64) NOT NULL DEFAULT '',
  `receiver_phone` VARCHAR(32) NOT NULL DEFAULT '',
  `receiver_address` VARCHAR(512) NOT NULL DEFAULT '' COMMENT 'deprecated flattened address',
  `receiver_country_code` VARCHAR(16) NOT NULL DEFAULT '',
  `receiver_province_code` VARCHAR(32) NOT NULL DEFAULT '',
  `receiver_province_name` VARCHAR(64) NOT NULL DEFAULT '',
  `receiver_city_code` VARCHAR(32) NOT NULL DEFAULT '',
  `receiver_city_name` VARCHAR(64) NOT NULL DEFAULT '',
  `receiver_district_code` VARCHAR(32) NOT NULL DEFAULT '',
  `receiver_district_name` VARCHAR(64) NOT NULL DEFAULT '',
  `receiver_street` VARCHAR(128) NOT NULL DEFAULT '',
  `receiver_detail` VARCHAR(255) NOT NULL DEFAULT '',
  `receiver_postal_code` VARCHAR(32) NOT NULL DEFAULT '',
  `version` BIGINT UNSIGNED NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shipment_no` (`shipment_no`),
  UNIQUE KEY `uk_sub_order_no` (`sub_order_no`),
  KEY `idx_shop_status_created` (`shop_no`,`shipment_status`,`created_at`),
  KEY `idx_user_created` (`user_id`,`created_at`),
  KEY `idx_order_no` (`order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `fulfillment_tracking_node` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `node_no` VARCHAR(64) NOT NULL,
  `shipment_no` VARCHAR(64) NOT NULL,
  `status_code` VARCHAR(64) NOT NULL DEFAULT '',
  `content` VARCHAR(512) NOT NULL DEFAULT '',
  `location` VARCHAR(255) NOT NULL DEFAULT '',
  `event_time` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_node_no` (`node_no`),
  KEY `idx_shipment_event_time` (`shipment_no`,`event_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `fulfillment_operate_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `shipment_no` VARCHAR(64) NOT NULL,
  `operator_type` VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'SYSTEM/SELLER/INTERNAL',
  `operator_id` VARCHAR(64) NOT NULL DEFAULT '',
  `action` VARCHAR(64) NOT NULL DEFAULT '',
  `from_status` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `to_status` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `remark` VARCHAR(512) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_shipment_created` (`shipment_no`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `fulfillment_outbox_event` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id` VARCHAR(64) NOT NULL,
  `event_type` VARCHAR(64) NOT NULL,
  `aggregate_type` VARCHAR(32) NOT NULL DEFAULT 'SHIPMENT',
  `aggregate_id` VARCHAR(64) NOT NULL,
  `payload_json` JSON NOT NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0 PENDING 1 PUBLISHED 2 FAILED',
  `retry_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `next_retry_at` DATETIME NULL,
  `last_error` VARCHAR(512) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_status_next_retry` (`status`,`next_retry_at`),
  KEY `idx_aggregate` (`aggregate_type`,`aggregate_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `fulfillment_idempotency` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `idempotency_key` VARCHAR(128) NOT NULL,
  `scope` VARCHAR(64) NOT NULL,
  `request_hash` VARCHAR(128) NOT NULL DEFAULT '',
  `resource_id` VARCHAR(64) NOT NULL DEFAULT '',
  `response_json` JSON NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 SUCCEEDED 2 PROCESSING 3 FAILED',
  `expire_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_scope_key` (`scope`,`idempotency_key`),
  KEY `idx_expire_at` (`expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
