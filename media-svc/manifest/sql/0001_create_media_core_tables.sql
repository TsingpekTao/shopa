-- ============================================================================
-- media-svc: core schema (from media v1 proto)
-- MySQL 8.0+
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1) Scene policy table
-- acl_type: 1 PUBLIC_READ, 2 PRIVATE
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `media_scene_policy` (
  `id`                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `scene_code`           VARCHAR(64) NOT NULL COMMENT 'Unique scene code, e.g. PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC',
  `acl_type`             TINYINT UNSIGNED NOT NULL COMMENT '1 PUBLIC_READ, 2 PRIVATE',
  `max_count`            INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Max assets per binding slot',
  `max_size_bytes`       BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Max file size in bytes, 0 means no limit',
  `allow_mime_json`      JSON NULL COMMENT 'Allowed mime list JSON',
  `allow_ext_json`       JSON NULL COMMENT 'Allowed extension list JSON',
  `retention_days`       INT UNSIGNED NOT NULL DEFAULT 30 COMMENT 'Retention days after unbound',
  `gc_grace_hours`       INT UNSIGNED NOT NULL DEFAULT 24 COMMENT 'Extra grace hours before physical delete',
  `risk_level`           TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 LOW, 2 MEDIUM, 3 HIGH',
  `risk_async_enabled`   TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Whether to enable async risk scanning',
  `process_async_enabled` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether to enable async document processing',
  `status`               TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 ENABLED, 2 DISABLED',
  `remark`               VARCHAR(255) NOT NULL DEFAULT '',
  `ext_json`             JSON NULL,
  `created_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_scene_code` (`scene_code`),
  KEY `idx_acl_status` (`acl_type`, `status`),
  KEY `idx_status_updated` (`status`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Media scene policy table';

-- ----------------------------------------------------------------------------
-- 2) Asset master table (supports asset tree)
-- asset_role: 1 ORIGINAL, 2 DERIVED
-- derived_kind: 0 UNSPECIFIED, 1 PREVIEW_IMAGE, 2 OCR_TEXT, 3 THUMBNAIL, 4 EMBEDDING
-- risk_status: 1 PENDING, 2 PASSED, 3 REJECTED
-- process_status: 1 PENDING, 2 PROCESSING, 3 SUCCESS, 4 FAILED, 5 PARTIAL_SUCCESS
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `media_asset` (
  `asset_id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_asset_id`      BIGINT UNSIGNED NULL COMMENT 'Parent asset id, null for root ORIGINAL asset',
  `root_asset_id`        BIGINT UNSIGNED NULL COMMENT 'Root asset id of the asset tree',
  `scene_code`           VARCHAR(64) NOT NULL,
  `acl_type`             TINYINT UNSIGNED NOT NULL COMMENT '1 PUBLIC_READ, 2 PRIVATE',
  `asset_role`           TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 ORIGINAL, 2 DERIVED',
  `derived_kind`         TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Derived asset subtype enum',
  `file_name`            VARCHAR(255) NOT NULL DEFAULT '',
  `mime_type`            VARCHAR(128) NOT NULL DEFAULT '',
  `size_bytes`           BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `checksum_sha256`      CHAR(64) NOT NULL DEFAULT '',
  `etag`                 VARCHAR(128) NOT NULL DEFAULT '',
  `storage_provider`     VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'Internal storage provider, e.g. minio/oss',
  `bucket`               VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Internal storage bucket',
  `object_key`           VARCHAR(1024) NOT NULL DEFAULT '' COMMENT 'Internal object key/path',
  `storage_object_hash`  CHAR(64) NOT NULL DEFAULT '' COMMENT 'SHA256 hash of storage_provider+bucket+object_key',
  `public_url`           VARCHAR(2048) NOT NULL DEFAULT '' COMMENT 'Static CDN URL for PUBLIC_READ scenes',
  `risk_status`          TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `process_status`       TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `process_progress`     TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `process_error_code`   VARCHAR(64) NOT NULL DEFAULT '',
  `process_error_message` VARCHAR(500) NOT NULL DEFAULT '',
  `processed_at`         DATETIME(3) NULL,
  `uploader_user_id`     BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Uploader user id from metadata',
  `trace_biz_type`       VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Optional trace-only biz type from InitUpload',
  `trace_biz_no`         VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Optional trace-only biz no from InitUpload',
  `deleted_at`           DATETIME(3) NULL COMMENT 'Soft delete mark, physical delete by GC',
  `created_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`asset_id`),
  UNIQUE KEY `uk_storage_object_hash` (`storage_object_hash`),
  KEY `idx_scene_created` (`scene_code`, `created_at`),
  KEY `idx_parent_asset_id` (`parent_asset_id`),
  KEY `idx_root_asset_id` (`root_asset_id`),
  KEY `idx_root_role_kind` (`root_asset_id`, `asset_role`, `derived_kind`),
  KEY `idx_process_status_updated` (`process_status`, `updated_at`),
  KEY `idx_risk_status_updated` (`risk_status`, `updated_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_checksum_sha256` (`checksum_sha256`),
  KEY `idx_uploader_created` (`uploader_user_id`, `created_at`),
  KEY `idx_trace_biz` (`trace_biz_type`, `trace_biz_no`, `created_at`),
  CONSTRAINT `chk_process_progress` CHECK (`process_progress` <= 100)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Media asset table with parent-child tree support';

-- ----------------------------------------------------------------------------
-- 3) Asset binding table (supports ordered/batch bind and unbind)
-- is_active: 1 active binding, 0 inactive history
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `media_asset_binding` (
  `id`                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `scene_code`           VARCHAR(64) NOT NULL,
  `biz_type`             VARCHAR(64) NOT NULL,
  `biz_no`               VARCHAR(64) NOT NULL,
  `binding_field`        VARCHAR(64) NOT NULL COMMENT 'Slot name, e.g. main_images/detail_images/carousel_images',
  `asset_id`             BIGINT UNSIGNED NOT NULL,
  `sort_order`           INT NOT NULL DEFAULT 0,
  `is_active`            TINYINT(1) NOT NULL DEFAULT 1,
  `active_asset_slot`    BIGINT UNSIGNED GENERATED ALWAYS AS (CASE WHEN `is_active` = 1 THEN `asset_id` ELSE NULL END) STORED,
  `active_sort_slot`     INT GENERATED ALWAYS AS (CASE WHEN `is_active` = 1 THEN `sort_order` ELSE NULL END) STORED,
  `operator_user_id`     BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `request_id`           VARCHAR(64) NOT NULL DEFAULT '',
  `unbound_at`           DATETIME(3) NULL,
  `created_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_active_asset` (`scene_code`, `biz_type`, `biz_no`, `binding_field`, `active_asset_slot`),
  UNIQUE KEY `uk_active_sort` (`scene_code`, `biz_type`, `biz_no`, `binding_field`, `active_sort_slot`),
  KEY `idx_biz_active_sort` (`scene_code`, `biz_type`, `biz_no`, `binding_field`, `is_active`, `sort_order`, `id`),
  KEY `idx_asset_active` (`asset_id`, `is_active`, `updated_at`),
  KEY `idx_unbound_at` (`unbound_at`),
  KEY `idx_request_id` (`request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Ordered business bindings for assets';

-- ----------------------------------------------------------------------------
-- 4) Idempotency table for write APIs
-- status: 1 PROCESSING, 2 SUCCEEDED, 3 FAILED
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `media_idempotency` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `actor_user_id`      BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'User id from metadata',
  `action_code`        VARCHAR(64) NOT NULL COMMENT 'InitUpload/BatchBindAssetsToBiz/...',
  `idempotency_key`    VARCHAR(128) NOT NULL COMMENT 'x-idempotency-key',
  `request_hash`       CHAR(64) NOT NULL DEFAULT '' COMMENT 'Request payload hash',
  `scene_code`         VARCHAR(64) NOT NULL DEFAULT '',
  `biz_type`           VARCHAR(64) NOT NULL DEFAULT '',
  `biz_no`             VARCHAR(64) NOT NULL DEFAULT '',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Write API idempotency table';

-- ----------------------------------------------------------------------------
-- 5) Outbox hot table for async event publishing
-- status: 1 NEW, 2 PROCESSING, 3 SENT, 4 FAILED, 5 DLQ
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `media_outbox_event` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id`          CHAR(36) NOT NULL COMMENT 'Global unique event id',
  `event_type`        VARCHAR(64) NOT NULL COMMENT 'AssetProcessingCompleted/AssetProcessingFailed/...',
  `aggregate_type`    VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'ASSET/BINDING/SCENE',
  `aggregate_key`     VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'asset_id/biz_key/scene_code',
  `request_id`        VARCHAR(64) NOT NULL DEFAULT '',
  `payload_json`      JSON NOT NULL,
  `status`            TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `available_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `sent_at`           DATETIME(3) NULL,
  `fail_count`        INT UNSIGNED NOT NULL DEFAULT 0,
  `last_error`        VARCHAR(512) NOT NULL DEFAULT '',
  `created_at`        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_status_available_id` (`status`, `available_at`, `id`),
  KEY `idx_type_created` (`event_type`, `created_at`),
  KEY `idx_aggregate_key` (`aggregate_type`, `aggregate_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Outbox hot table for media events';

-- ----------------------------------------------------------------------------
-- 6) Outbox archive table
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `media_outbox_event_archive` (
  `id`                BIGINT UNSIGNED NOT NULL,
  `event_id`          CHAR(36) NOT NULL,
  `event_type`        VARCHAR(64) NOT NULL,
  `aggregate_type`    VARCHAR(32) NOT NULL DEFAULT '',
  `aggregate_key`     VARCHAR(128) NOT NULL DEFAULT '',
  `request_id`        VARCHAR(64) NOT NULL DEFAULT '',
  `payload_json`      JSON NOT NULL,
  `status`            TINYINT UNSIGNED NOT NULL DEFAULT 3,
  `available_at`      DATETIME(3) NOT NULL,
  `sent_at`           DATETIME(3) NULL,
  `fail_count`        INT UNSIGNED NOT NULL DEFAULT 0,
  `last_error`        VARCHAR(512) NOT NULL DEFAULT '',
  `created_at`        DATETIME(3) NOT NULL,
  `updated_at`        DATETIME(3) NOT NULL,
  `archived_at`       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_archived_at` (`archived_at`),
  KEY `idx_status_sent_at` (`status`, `sent_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Outbox archive table for media events';

-- ----------------------------------------------------------------------------
-- 7) Consumer dedup table for inbound MQ events
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `media_consumer_event_dedup` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `consumer_name`     VARCHAR(64) NOT NULL,
  `event_id`          CHAR(36) NOT NULL,
  `payload_hash`      CHAR(64) NOT NULL DEFAULT '',
  `first_seen_at`     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_consumer_event` (`consumer_name`, `event_id`),
  KEY `idx_first_seen_at` (`first_seen_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Consumer idempotency table';
