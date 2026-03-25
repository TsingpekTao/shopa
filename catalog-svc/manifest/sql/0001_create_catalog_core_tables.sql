-- ============================================================================
-- catalog-svc: core schema (generated from catalog v1 proto contract)
-- MySQL 8.0+
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1) Category dictionary
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_category` (
  `category_id`        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `category_name`      VARCHAR(128) NOT NULL,
  `parent_id`          BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `level`              TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `path`               VARCHAR(512) NOT NULL DEFAULT '',
  `sort_order`         INT NOT NULL DEFAULT 0,
  `is_leaf`            TINYINT(1) NOT NULL DEFAULT 1,
  `status`             TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 ENABLED, 2 DISABLED',
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`category_id`),
  KEY `idx_parent_sort` (`parent_id`, `sort_order`),
  KEY `idx_status_updated` (`status`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Category tree';

-- ----------------------------------------------------------------------------
-- 2) Brand dictionary
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_brand` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `brand_no`           VARCHAR(32) NOT NULL,
  `brand_name`         VARCHAR(128) NOT NULL,
  `status`             TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 ENABLED, 2 DISABLED',
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_brand_no` (`brand_no`),
  KEY `idx_status_updated` (`status`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Brand dictionary';

-- ----------------------------------------------------------------------------
-- 3) Category attribute template
-- attr_scope: 1 SPU, 2 SKU_SALE, 3 EXT
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_attribute_template` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `category_id`        BIGINT UNSIGNED NOT NULL,
  `attr_code`          VARCHAR(64) NOT NULL,
  `attr_name`          VARCHAR(128) NOT NULL,
  `attr_scope`         TINYINT UNSIGNED NOT NULL COMMENT '1 SPU, 2 SKU_SALE, 3 EXT',
  `value_type`         VARCHAR(32) NOT NULL DEFAULT 'TEXT' COMMENT 'TEXT/NUMBER/BOOL/ENUM',
  `required_flag`      TINYINT(1) NOT NULL DEFAULT 0,
  `searchable_flag`    TINYINT(1) NOT NULL DEFAULT 0,
  `filterable_flag`    TINYINT(1) NOT NULL DEFAULT 0,
  `options_json`       JSON NULL COMMENT 'Enum options if needed',
  `sort_order`         INT NOT NULL DEFAULT 0,
  `status`             TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 ENABLED, 2 DISABLED',
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_category_attr_code` (`category_id`, `attr_code`),
  KEY `idx_category_scope_status` (`category_id`, `attr_scope`, `status`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Attribute template by category';

-- ----------------------------------------------------------------------------
-- 4) SPU table
-- spu_status:
--   1 DRAFT, 2 REVIEWING, 3 APPROVED, 4 ON_SHELF, 5 OFF_SHELF,
--   6 REJECTED, 7 FROZEN, 8 DELETED
-- spu_stock_status: 1 IN_STOCK, 2 OUT_OF_STOCK
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_spu` (
  `id`                       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `spu_no`                   VARCHAR(32) NOT NULL,
  `shop_no`                  VARCHAR(32) NOT NULL,
  `title`                    VARCHAR(255) NOT NULL,
  `sub_title`                VARCHAR(255) NOT NULL DEFAULT '',
  `category_id`              BIGINT UNSIGNED NOT NULL,
  `brand_no`                 VARCHAR(32) NOT NULL DEFAULT '',
  `main_image_asset_ids_json`   JSON NULL,
  `detail_image_asset_ids_json` JSON NULL,
  `spu_status`               TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `spu_stock_status`         TINYINT UNSIGNED NOT NULL DEFAULT 2,
  `min_sale_price`           BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `max_sale_price`           BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `min_market_price`         BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `max_market_price`         BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `publish_time`             DATETIME(3) NULL,
  `version`                  INT UNSIGNED NOT NULL DEFAULT 1,
  `review_status`            TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 PENDING, 2 APPROVED, 3 REJECTED',
  `reject_reason_code`       VARCHAR(64) NOT NULL DEFAULT '',
  `reject_comment`           VARCHAR(500) NOT NULL DEFAULT '',
  `review_comment`           VARCHAR(500) NOT NULL DEFAULT '',
  `reviewer_id`              BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `submitted_at`             DATETIME(3) NULL,
  `reviewed_at`              DATETIME(3) NULL,
  `sold_count`               BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `deleted_at`               DATETIME(3) NULL,
  `created_at`               DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`               DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_spu_no` (`spu_no`),
  KEY `idx_shop_status_updated` (`shop_no`, `spu_status`, `updated_at`),
  KEY `idx_category_status_price` (`category_id`, `spu_status`, `min_sale_price`),
  KEY `idx_status_publish_time` (`spu_status`, `publish_time`),
  KEY `idx_spu_stock_status` (`spu_stock_status`, `updated_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SPU master table';

-- ----------------------------------------------------------------------------
-- 5) SPU attribute values
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_spu_attr_value` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `spu_no`             VARCHAR(32) NOT NULL,
  `attr_code`          VARCHAR(64) NOT NULL,
  `attr_name`          VARCHAR(128) NOT NULL DEFAULT '',
  `attr_scope`         TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `attr_value`         TEXT NOT NULL,
  `sort_order`         INT NOT NULL DEFAULT 0,
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_spu_attr` (`spu_no`, `attr_code`),
  KEY `idx_spu_sort` (`spu_no`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SPU attribute values';

-- ----------------------------------------------------------------------------
-- 6) SKU table
-- sku_status: 1 DRAFT, 2 ENABLED, 3 DISABLED, 4 DELETED
-- stock_status: 1 IN_STOCK, 2 OUT_OF_STOCK (projection only)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_sku` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sku_no`             VARCHAR(32) NOT NULL,
  `spu_no`             VARCHAR(32) NOT NULL,
  `shop_no`            VARCHAR(32) NOT NULL,
  `sku_name`           VARCHAR(255) NOT NULL,
  `sku_image_asset_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `sku_status`         TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `stock_status`       TINYINT UNSIGNED NOT NULL DEFAULT 2,
  `stock_version`      BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'From inventory event',
  `last_stock_event_id` CHAR(36) NOT NULL DEFAULT '',
  `sale_price`         BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `market_price`       BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `sale_specs_json`    JSON NULL COMMENT 'Canonical sale attrs json for display',
  `sale_specs_hash`    CHAR(64) NOT NULL DEFAULT '' COMMENT 'Canonical hash for unique combination',
  `active_specs_hash`  CHAR(64) GENERATED ALWAYS AS (CASE WHEN `deleted_at` IS NULL THEN `sale_specs_hash` ELSE NULL END) STORED,
  `sort_order`         INT NOT NULL DEFAULT 0,
  `version`            INT UNSIGNED NOT NULL DEFAULT 1,
  `deleted_at`         DATETIME(3) NULL,
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_no` (`sku_no`),
  UNIQUE KEY `uk_spu_active_specs` (`spu_no`, `active_specs_hash`),
  KEY `idx_spu_status_sort` (`spu_no`, `sku_status`, `sort_order`),
  KEY `idx_shop_status_updated` (`shop_no`, `sku_status`, `updated_at`),
  KEY `idx_stock_status_version` (`stock_status`, `stock_version`),
  KEY `idx_price` (`sale_price`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SKU table with stock projection';

-- ----------------------------------------------------------------------------
-- 7) SKU sale attrs
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_sku_sale_attr_value` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sku_no`             VARCHAR(32) NOT NULL,
  `spu_no`             VARCHAR(32) NOT NULL,
  `attr_code`          VARCHAR(64) NOT NULL,
  `attr_name`          VARCHAR(128) NOT NULL DEFAULT '',
  `attr_value`         VARCHAR(255) NOT NULL DEFAULT '',
  `sort_order`         INT NOT NULL DEFAULT 0,
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_attr` (`sku_no`, `attr_code`),
  KEY `idx_spu_attr_value` (`spu_no`, `attr_code`, `attr_value`),
  KEY `idx_sku_sort` (`sku_no`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SKU sale attributes';

-- ----------------------------------------------------------------------------
-- 8) Review task/history table
-- review_status: 1 PENDING, 2 APPROVED, 3 REJECTED
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_review_task` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_no`            VARCHAR(32) NOT NULL,
  `spu_no`             VARCHAR(32) NOT NULL,
  `shop_no`            VARCHAR(32) NOT NULL,
  `spu_version_at_submit` INT UNSIGNED NOT NULL DEFAULT 0,
  `submit_note`        VARCHAR(500) NOT NULL DEFAULT '',
  `review_status`      TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `reject_reason_code` VARCHAR(64) NOT NULL DEFAULT '',
  `reject_comment`     VARCHAR(500) NOT NULL DEFAULT '',
  `review_comment`     VARCHAR(500) NOT NULL DEFAULT '',
  `reviewer_id`        BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `submitted_at`       DATETIME(3) NULL,
  `review_started_at`  DATETIME(3) NULL,
  `reviewed_at`        DATETIME(3) NULL,
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_no` (`task_no`),
  KEY `idx_review_status_submitted` (`review_status`, `submitted_at`),
  KEY `idx_shop_submitted` (`shop_no`, `submitted_at`),
  KEY `idx_spu_no` (`spu_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Review tasks/history';

-- ----------------------------------------------------------------------------
-- 9) Idempotency table for write APIs
-- status: 1 PROCESSING, 2 SUCCEEDED, 3 FAILED
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_idempotency` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `actor_user_id`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `action_code`        VARCHAR(64) NOT NULL,
  `idempotency_key`    VARCHAR(128) NOT NULL,
  `request_hash`       CHAR(64) NOT NULL DEFAULT '',
  `shop_no`            VARCHAR(32) NOT NULL DEFAULT '',
  `spu_no`             VARCHAR(32) NOT NULL DEFAULT '',
  `status`             TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `response_code`      INT NOT NULL DEFAULT 0,
  `response_json`      JSON NULL,
  `expired_at`         DATETIME(3) NOT NULL,
  `created_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`         DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_actor_action_key` (`actor_user_id`, `action_code`, `idempotency_key`),
  KEY `idx_expired_at` (`expired_at`),
  KEY `idx_action_created` (`action_code`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Write idempotency';

-- ----------------------------------------------------------------------------
-- 10) Outbox hot table
-- status: 1 NEW, 2 PROCESSING, 3 SENT, 4 FAILED, 5 DLQ
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_outbox_event` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id`           CHAR(36) NOT NULL,
  `event_type`         VARCHAR(64) NOT NULL,
  `aggregate_type`     VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'SPU/SKU/REVIEW',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Catalog outbox hot table';

-- ----------------------------------------------------------------------------
-- 11) Outbox archive table
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_outbox_event_archive` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Catalog outbox archive';

-- ----------------------------------------------------------------------------
-- 12) Consumer dedup table
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `catalog_consumer_event_dedup` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `consumer_name`      VARCHAR(64) NOT NULL,
  `event_id`           CHAR(36) NOT NULL,
  `payload_hash`       CHAR(64) NOT NULL DEFAULT '',
  `first_seen_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_consumer_event` (`consumer_name`, `event_id`),
  KEY `idx_first_seen_at` (`first_seen_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Inbound event idempotency';

