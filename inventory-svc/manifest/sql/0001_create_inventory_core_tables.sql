-- ============================================================================
-- inventory-svc: core schema (generated from inventory v1 proto contract)
-- MySQL 8.0+
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1) SKU context table (source of ownership and aggregate relations)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_sku_context` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sku_no`             VARCHAR(32) NOT NULL,
  `spu_no`             VARCHAR(32) NOT NULL,
  `shop_no`            VARCHAR(32) NOT NULL,
  `enabled`            TINYINT(1) NOT NULL DEFAULT 1,
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_no` (`sku_no`),
  KEY `idx_spu_enabled` (`spu_no`, `enabled`),
  KEY `idx_shop_enabled` (`shop_no`, `enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SKU ownership/context cache';

-- ----------------------------------------------------------------------------
-- 2) Inventory stock source-of-truth
-- stock_status: 1 IN_STOCK, 2 OUT_OF_STOCK
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_stock` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sku_no`             VARCHAR(32) NOT NULL,
  `spu_no`             VARCHAR(32) NOT NULL,
  `shop_no`            VARCHAR(32) NOT NULL,
  `total_qty`          BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `locked_qty`         BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `available_qty`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `stock_version`      BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Monotonic version for projection events',
  `stock_status`       TINYINT UNSIGNED NOT NULL DEFAULT 2,
  `is_hot`             TINYINT(1) NOT NULL DEFAULT 0,
  `row_version`        BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Internal optimistic lock version',
  `last_event_id`      CHAR(36) NOT NULL DEFAULT '',
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_no` (`sku_no`),
  KEY `idx_shop_status` (`shop_no`, `stock_status`, `updated_at`),
  KEY `idx_spu_status` (`spu_no`, `stock_status`, `updated_at`),
  KEY `idx_hot_sku` (`is_hot`, `updated_at`),
  KEY `idx_stock_version` (`stock_version`),
  CONSTRAINT `chk_locked_le_total` CHECK (`locked_qty` <= `total_qty`),
  CONSTRAINT `chk_available_formula` CHECK (`available_qty` + `locked_qty` = `total_qty`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Inventory source of truth';

-- ----------------------------------------------------------------------------
-- 3) Reservation header
-- reservation_status: 1 RESERVED, 2 CONFIRMED, 3 CANCELED, 4 EXPIRED, 5 FAILED
-- reserve_mode: 1 ALL_OR_NOTHING
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_reservation` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `reservation_no`     VARCHAR(64) NOT NULL,
  `order_no`           VARCHAR(64) NOT NULL,
  `user_id`            BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `reservation_status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `reserve_mode`       TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `expired_at`         DATETIME(3) NOT NULL,
  `confirmed_at`       DATETIME(3) NULL,
  `canceled_at`        DATETIME(3) NULL,
  `cancel_reason_code` VARCHAR(64) NOT NULL DEFAULT '',
  `idempotency_key`    VARCHAR(128) NULL,
  `request_id`         VARCHAR(64) NOT NULL DEFAULT '',
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reservation_no` (`reservation_no`),
  UNIQUE KEY `uk_idempotency_key` (`idempotency_key`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_user_created` (`user_id`, `created_at`),
  KEY `idx_status_expired` (`reservation_status`, `expired_at`),
  KEY `idx_request_id` (`request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Reservation header';

-- ----------------------------------------------------------------------------
-- 4) Reservation items
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_reservation_item` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `reservation_no`     VARCHAR(64) NOT NULL,
  `order_no`           VARCHAR(64) NOT NULL,
  `sku_no`             VARCHAR(32) NOT NULL,
  `spu_no`             VARCHAR(32) NOT NULL,
  `shop_no`            VARCHAR(32) NOT NULL,
  `qty`                INT UNSIGNED NOT NULL,
  `stock_version`      BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Version after reserve success',
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reserve_sku` (`reservation_no`, `sku_no`),
  KEY `idx_order_sku` (`order_no`, `sku_no`),
  KEY `idx_shop_created` (`shop_no`, `created_at`),
  KEY `idx_sku` (`sku_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Reservation item details';

-- ----------------------------------------------------------------------------
-- 5) Inventory operation ledger
-- action_code examples: RESERVE/CONFIRM/CANCEL/ADJUST
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_ledger` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `txn_no`             CHAR(36) NOT NULL,
  `biz_type`           VARCHAR(64) NOT NULL DEFAULT '',
  `biz_no`             VARCHAR(64) NOT NULL DEFAULT '',
  `action_code`        VARCHAR(32) NOT NULL,
  `reservation_no`     VARCHAR(64) NOT NULL DEFAULT '',
  `order_no`           VARCHAR(64) NOT NULL DEFAULT '',
  `sku_no`             VARCHAR(32) NOT NULL,
  `spu_no`             VARCHAR(32) NOT NULL,
  `shop_no`            VARCHAR(32) NOT NULL,
  `delta_total_qty`    BIGINT NOT NULL DEFAULT 0,
  `delta_locked_qty`   BIGINT NOT NULL DEFAULT 0,
  `delta_available_qty` BIGINT NOT NULL DEFAULT 0,
  `after_total_qty`    BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `after_locked_qty`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `after_available_qty` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `stock_version`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `operator_type`      VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'SELLER/ORDER/ADMIN/SYSTEM',
  `operator_user_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `request_id`         VARCHAR(64) NOT NULL DEFAULT '',
  `remark`             VARCHAR(500) NOT NULL DEFAULT '',
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_txn_no` (`txn_no`),
  UNIQUE KEY `uk_biz_action_sku` (`biz_type`, `biz_no`, `action_code`, `sku_no`),
  KEY `idx_sku_created` (`sku_no`, `created_at`),
  KEY `idx_order_created` (`order_no`, `created_at`),
  KEY `idx_reservation_created` (`reservation_no`, `created_at`),
  KEY `idx_shop_created` (`shop_no`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Inventory ledger table';

-- ----------------------------------------------------------------------------
-- 6) Idempotency table for write APIs
-- status: 1 PROCESSING, 2 SUCCEEDED, 3 FAILED
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_idempotency` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `actor_scope`        VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'SELLER/ORDER/ADMIN/SYSTEM',
  `actor_id`           BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `action_code`        VARCHAR(64) NOT NULL,
  `idempotency_key`    VARCHAR(128) NOT NULL,
  `request_hash`       CHAR(64) NOT NULL DEFAULT '',
  `status`             TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `response_code`      INT NOT NULL DEFAULT 0,
  `response_json`      JSON NULL,
  `expired_at`         DATETIME(3) NOT NULL,
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_actor_action_key` (`actor_scope`, `actor_id`, `action_code`, `idempotency_key`),
  KEY `idx_expired_at` (`expired_at`),
  KEY `idx_action_created` (`action_code`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Write idempotency';

-- ----------------------------------------------------------------------------
-- 7) Outbox hot table
-- status: 1 NEW, 2 PROCESSING, 3 SENT, 4 FAILED, 5 DLQ
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_outbox_event` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id`           CHAR(36) NOT NULL,
  `event_type`         VARCHAR(64) NOT NULL,
  `aggregate_type`     VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'STOCK/RESERVATION',
  `aggregate_key`      VARCHAR(128) NOT NULL DEFAULT '',
  `request_id`         VARCHAR(64) NOT NULL DEFAULT '',
  `payload_json`       JSON NOT NULL,
  `status`             TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `available_at`       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `sent_at`            DATETIME(3) NULL,
  `fail_count`         INT UNSIGNED NOT NULL DEFAULT 0,
  `last_error`         VARCHAR(512) NOT NULL DEFAULT '',
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_status_available_id` (`status`, `available_at`, `id`),
  KEY `idx_type_created` (`event_type`, `created_at`),
  KEY `idx_aggregate_key` (`aggregate_type`, `aggregate_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Inventory outbox hot table';

-- ----------------------------------------------------------------------------
-- 8) Outbox archive table
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_outbox_event_archive` (
  `id`                 BIGINT UNSIGNED NOT NULL,
  `event_id`           CHAR(36) NOT NULL,
  `event_type`         VARCHAR(64) NOT NULL,
  `aggregate_type`     VARCHAR(32) NOT NULL DEFAULT '',
  `aggregate_key`      VARCHAR(128) NOT NULL DEFAULT '',
  `request_id`         VARCHAR(64) NOT NULL DEFAULT '',
  `payload_json`       JSON NOT NULL,
  `status`             TINYINT UNSIGNED NOT NULL DEFAULT 3,
  `available_at`       DATETIME(3) NOT NULL,
  `sent_at`            DATETIME(3) NULL,
  `fail_count`         INT UNSIGNED NOT NULL DEFAULT 0,
  `last_error`         VARCHAR(512) NOT NULL DEFAULT '',
  `created_at`         DATETIME(3) NOT NULL,
  `updated_at`         DATETIME(3) NOT NULL,
  `archived_at`        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_archived_at` (`archived_at`),
  KEY `idx_status_sent_at` (`status`, `sent_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Inventory outbox archive';

-- ----------------------------------------------------------------------------
-- 9) Consumer dedup table
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `inventory_consumer_event_dedup` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `consumer_name`      VARCHAR(64) NOT NULL,
  `event_id`           CHAR(36) NOT NULL,
  `payload_hash`       CHAR(64) NOT NULL DEFAULT '',
  `first_seen_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_consumer_event` (`consumer_name`, `event_id`),
  KEY `idx_first_seen_at` (`first_seen_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Inbound event idempotency';
