-- ============================================================================
-- seller-shop-svc: core schema (from sellershop v1 proto + event proto)
-- MySQL 8.0+
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1) Master entity table (approved latest legal subject data)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_entity` (
  `id`                           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Internal PK',
  `entity_no`                    VARCHAR(40) NOT NULL COMMENT 'External business id',
  `owner_user_id`                BIGINT UNSIGNED NOT NULL COMMENT 'IAM user id',
  `merchant_type_code`           VARCHAR(32) NOT NULL COMMENT 'PERSONAL/ENTERPRISE/...',
  `entity_name`                  VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Display name',
  `contact_name`                 VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Contact person',
  `contact_phone`                VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'Contact phone',
  `contact_email`                VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Contact email',
  `subject_type_code`            VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'Legal subject type',
  `subject_name`                 VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Legal subject name',
  `subject_cert_no`              VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Legal certificate no',
  `legal_representative_name`    VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Legal representative',
  `legal_representative_cert_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Legal rep certificate no',
  `cert_valid_from`              DATETIME(3) NULL COMMENT 'Legal cert valid from',
  `cert_valid_until`             DATETIME(3) NULL COMMENT 'Legal cert valid until',
  `legal_subject_ext_json`       JSON NULL COMMENT 'Ext fields for legal subject',
  `ext_json`                     JSON NULL COMMENT 'Entity ext fields',
  `version`                      INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Optimistic version',
  `created_at`                   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`                   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_entity_no` (`entity_no`),
  KEY `idx_owner_user_id` (`owner_user_id`),
  KEY `idx_owner_updated` (`owner_user_id`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Seller legal entity master table';

-- ----------------------------------------------------------------------------
-- 2) Entity qualification docs (approved latest docs)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_entity_doc` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `entity_id`     BIGINT UNSIGNED NOT NULL COMMENT 'FK seller_entity.id',
  `entity_no`     VARCHAR(40) NOT NULL COMMENT 'Redundant external entity no',
  `doc_type_code` VARCHAR(32) NOT NULL COMMENT 'BUSINESS_LICENSE/ID_FRONT/...',
  `asset_id`      BIGINT UNSIGNED NOT NULL COMMENT 'media-svc asset id',
  `valid_from`    DATETIME(3) NULL,
  `valid_until`   DATETIME(3) NULL,
  `issuer`        VARCHAR(128) NOT NULL DEFAULT '',
  `ext_json`      JSON NULL,
  `status`        TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 ACTIVE,2 INACTIVE',
  `created_at`    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_entity_doc_asset` (`entity_id`, `doc_type_code`, `asset_id`),
  KEY `idx_entity_doc_type` (`entity_id`, `doc_type_code`),
  KEY `idx_entity_no` (`entity_no`),
  CONSTRAINT `fk_seller_entity_doc_entity_id`
    FOREIGN KEY (`entity_id`) REFERENCES `seller_entity` (`id`)
    ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Seller legal docs table';

-- ----------------------------------------------------------------------------
-- 3) Shop master table
-- status: 1 PROVISIONING,2 ACTIVE,3 FREEZING,4 FROZEN,5 CLOSING,6 CLOSED,7 PROVISION_FAILED
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_shop` (
  `id`                     BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Internal shop id (scope_id)',
  `shop_no`                VARCHAR(40) NOT NULL COMMENT 'External business shop id',
  `entity_id`              BIGINT UNSIGNED NOT NULL COMMENT 'FK seller_entity.id',
  `entity_no`              VARCHAR(40) NOT NULL COMMENT 'External entity id',
  `application_no`         VARCHAR(40) NOT NULL DEFAULT '' COMMENT 'Source application no',
  `owner_user_id`          BIGINT UNSIGNED NOT NULL COMMENT 'IAM user id',
  `shop_name`              VARCHAR(128) NOT NULL COMMENT 'Requested unique shop name',
  `shop_name_norm`         VARCHAR(128) NOT NULL COMMENT 'Normalized lower name',
  `shop_display_name`      VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Display name',
  `shop_type_code`         VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'Shop type code',
  `main_category_ids_json` JSON NULL COMMENT 'Main category id list',
  `logo_asset_id`          BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `banner_asset_id`        BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `service_phone`          VARCHAR(32) NOT NULL DEFAULT '',
  `service_email`          VARCHAR(128) NOT NULL DEFAULT '',
  `description`            VARCHAR(1024) NOT NULL DEFAULT '',
  `ext_json`               JSON NULL,
  `status`                 TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'ShopStatus enum',
  `buyer_visible`          TINYINT(1) NOT NULL DEFAULT 0 COMMENT '1 visible to buyers',
  `version`                INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Optimistic version',
  `freeze_reason_code`     VARCHAR(32) NOT NULL DEFAULT '',
  `freeze_reason`          VARCHAR(255) NOT NULL DEFAULT '',
  `close_reason_code`      VARCHAR(32) NOT NULL DEFAULT '',
  `close_reason`           VARCHAR(255) NOT NULL DEFAULT '',
  `provision_started_at`   DATETIME(3) NULL,
  `provision_finished_at`  DATETIME(3) NULL,
  `activated_at`           DATETIME(3) NULL,
  `frozen_at`              DATETIME(3) NULL,
  `closed_at`              DATETIME(3) NULL,
  `created_at`             DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`             DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shop_no` (`shop_no`),
  UNIQUE KEY `uk_shop_name_norm` (`shop_name_norm`),
  KEY `idx_owner_status_updated` (`owner_user_id`, `status`, `updated_at`),
  KEY `idx_status_visible_updated` (`status`, `buyer_visible`, `updated_at`),
  KEY `idx_application_no` (`application_no`),
  KEY `idx_entity_no` (`entity_no`),
  CONSTRAINT `fk_seller_shop_entity_id`
    FOREIGN KEY (`entity_id`) REFERENCES `seller_entity` (`id`)
    ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Seller shop master table';

-- ----------------------------------------------------------------------------
-- 4) Application table
-- status: 1 DRAFT,2 SUBMITTED,3 REVIEWING,4 APPROVED,5 REJECTED,6 CANCELLED,7 SYSTEM_REJECTED
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_application` (
  `id`                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `application_no`       VARCHAR(40) NOT NULL COMMENT 'External application id',
  `version`              INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Optimistic version',
  `previous_application_no` VARCHAR(40) NULL COMMENT 'Chain link to previous rejected application',
  `next_application_no`  VARCHAR(40) NULL COMMENT 'Chain link to next resubmitted application',
  `owner_user_id`        BIGINT UNSIGNED NOT NULL,
  `status`               TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'ApplicationStatus enum',
  `entity_no`            VARCHAR(40) NOT NULL DEFAULT '',
  `shop_no`              VARCHAR(40) NOT NULL DEFAULT '',
  `entity_name`          VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Denormalized for list/search',
  `shop_name`            VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Denormalized for list/search',
  `shop_name_norm`       VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Denormalized normalized shop name',
  `entity_draft_json`    JSON NOT NULL COMMENT 'Current draft entity payload',
  `shop_draft_json`      JSON NOT NULL COMMENT 'Current draft shop payload',
  `entity_submitted_json` JSON NULL COMMENT 'Snapshot for review at submit time',
  `shop_submitted_json`  JSON NULL COMMENT 'Snapshot for review at submit time',
  `latest_reject_json`   JSON NULL COMMENT 'Latest reject info',
  `submitted_at`         DATETIME(3) NULL,
  `review_started_at`    DATETIME(3) NULL,
  `reviewed_at`          DATETIME(3) NULL,
  `reviewer_id`          BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `review_comment`       VARCHAR(500) NOT NULL DEFAULT '',
  `created_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_application_no` (`application_no`),
  UNIQUE KEY `uk_previous_application_no` (`previous_application_no`),
  UNIQUE KEY `uk_next_application_no` (`next_application_no`),
  KEY `idx_owner_status_created` (`owner_user_id`, `status`, `created_at`),
  KEY `idx_status_created` (`status`, `created_at`),
  KEY `idx_owner_created` (`owner_user_id`, `created_at`),
  KEY `idx_shop_name_norm` (`shop_name_norm`),
  KEY `idx_keyword` (`shop_name`, `entity_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Seller application table with draft/snapshot';

-- ----------------------------------------------------------------------------
-- 5) Shop name registry (reserve/bind/release)
-- status: 1 RESERVED,2 BOUND,3 RELEASED,4 EXPIRED
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_shop_name_registry` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `shop_name_norm`  VARCHAR(128) NOT NULL COMMENT 'Normalized unique shop name',
  `shop_name`       VARCHAR(128) NOT NULL COMMENT 'Original input shop name',
  `owner_user_id`   BIGINT UNSIGNED NOT NULL,
  `application_no`  VARCHAR(40) NOT NULL DEFAULT '',
  `status`          TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `expire_at`       DATETIME(3) NOT NULL COMMENT 'Reservation expiration time',
  `bound_shop_no`   VARCHAR(40) NOT NULL DEFAULT '' COMMENT 'Final bound shop_no when approved',
  `created_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shop_name_norm` (`shop_name_norm`),
  UNIQUE KEY `uk_application_no` (`application_no`),
  KEY `idx_owner_status_expire` (`owner_user_id`, `status`, `expire_at`),
  KEY `idx_status_expire` (`status`, `expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Global shop-name reservation table';

-- ----------------------------------------------------------------------------
-- 6) Idempotency table for write APIs (x-idempotency-key via metadata/header)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_idempotency` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `actor_user_id`     BIGINT UNSIGNED NOT NULL COMMENT 'User id or operator id',
  `action_code`       VARCHAR(64) NOT NULL COMMENT 'CreateApplicationDraft/ApproveApplication/...',
  `idempotency_key`   VARCHAR(128) NOT NULL COMMENT 'x-idempotency-key',
  `request_hash`      CHAR(64) NOT NULL DEFAULT '' COMMENT 'Payload hash for mismatch detection',
  `resource_type`     VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'application/shop/entity',
  `resource_no`       VARCHAR(40) NOT NULL DEFAULT '' COMMENT 'Business id',
  `status`            TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 PROCESSING,2 SUCCEEDED,3 FAILED',
  `response_code`     INT NOT NULL DEFAULT 0,
  `response_json`     JSON NULL,
  `expired_at`        DATETIME(3) NOT NULL,
  `created_at`        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_actor_action_key` (`actor_user_id`, `action_code`, `idempotency_key`),
  KEY `idx_expired_at` (`expired_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Write API idempotency registry';

-- ----------------------------------------------------------------------------
-- 7) Outbox hot table
-- status: 1 NEW,2 PROCESSING,3 SENT,4 FAILED,5 DLQ
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_outbox_event` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id`      CHAR(36) NOT NULL COMMENT 'Global unique event id',
  `event_type`    VARCHAR(64) NOT NULL COMMENT 'SellerRoleAssignRequested/...',
  `aggregate_type` VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'APPLICATION/SHOP/ENTITY',
  `aggregate_no`  VARCHAR(40) NOT NULL DEFAULT '' COMMENT 'application_no/shop_no/entity_no',
  `request_id`    VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Trace request id',
  `payload_json`  JSON NOT NULL,
  `status`        TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `available_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `sent_at`       DATETIME(3) NULL,
  `fail_count`    INT UNSIGNED NOT NULL DEFAULT 0,
  `last_error`    VARCHAR(512) NOT NULL DEFAULT '',
  `created_at`    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_status_available_id` (`status`, `available_at`, `id`),
  KEY `idx_type_created` (`event_type`, `created_at`),
  KEY `idx_aggregate_no` (`aggregate_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Seller outbox hot table';

-- ----------------------------------------------------------------------------
-- 8) Outbox archive table
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_outbox_event_archive` (
  `id`            BIGINT UNSIGNED NOT NULL,
  `event_id`      CHAR(36) NOT NULL,
  `event_type`    VARCHAR(64) NOT NULL,
  `aggregate_type` VARCHAR(32) NOT NULL DEFAULT '',
  `aggregate_no`  VARCHAR(40) NOT NULL DEFAULT '',
  `request_id`    VARCHAR(64) NOT NULL DEFAULT '',
  `payload_json`  JSON NOT NULL,
  `status`        TINYINT UNSIGNED NOT NULL DEFAULT 3,
  `available_at`  DATETIME(3) NOT NULL,
  `sent_at`       DATETIME(3) NULL,
  `fail_count`    INT UNSIGNED NOT NULL DEFAULT 0,
  `last_error`    VARCHAR(512) NOT NULL DEFAULT '',
  `created_at`    DATETIME(3) NOT NULL,
  `updated_at`    DATETIME(3) NOT NULL,
  `archived_at`   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_id` (`event_id`),
  KEY `idx_archived_at` (`archived_at`),
  KEY `idx_status_sent_at` (`status`, `sent_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Seller outbox archive table';

-- ----------------------------------------------------------------------------
-- 9) Consumer dedup table for inbound MQ events
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_consumer_event_dedup` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `consumer_name` VARCHAR(64) NOT NULL COMMENT 'seller-role-assigned-consumer/...',
  `event_id`      CHAR(36) NOT NULL COMMENT 'Incoming event id',
  `payload_hash`  CHAR(64) NOT NULL DEFAULT '',
  `first_seen_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_consumer_event` (`consumer_name`, `event_id`),
  KEY `idx_first_seen_at` (`first_seen_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='MQ consumer idempotency table';

-- ----------------------------------------------------------------------------
-- 10) Application operation log
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_application_op_log` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `application_no`   VARCHAR(40) NOT NULL,
  `owner_user_id`    BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `operator_user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Admin or seller operator',
  `action_code`      VARCHAR(32) NOT NULL COMMENT 'CREATE_DRAFT/SUBMIT/APPROVE/REJECT/RESUBMIT',
  `from_status`      TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `to_status`        TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `comment`          VARCHAR(512) NOT NULL DEFAULT '',
  `extra_json`       JSON NULL,
  `created_at`       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_application_created` (`application_no`, `created_at`),
  KEY `idx_operator_created` (`operator_user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Application lifecycle operation log';

-- ----------------------------------------------------------------------------
-- 11) Shop status operation log
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `seller_shop_op_log` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `shop_no`          VARCHAR(40) NOT NULL,
  `owner_user_id`    BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `operator_user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `action_code`      VARCHAR(32) NOT NULL COMMENT 'PROVISION_START/ACTIVATE/FREEZE/CLOSE/REOPEN',
  `from_status`      TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `to_status`        TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `reason_code`      VARCHAR(32) NOT NULL DEFAULT '',
  `reason`           VARCHAR(512) NOT NULL DEFAULT '',
  `request_id`       VARCHAR(64) NOT NULL DEFAULT '',
  `event_id`         CHAR(36) NOT NULL DEFAULT '',
  `created_at`       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_shop_created` (`shop_no`, `created_at`),
  KEY `idx_operator_created` (`operator_user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Shop lifecycle operation log';

