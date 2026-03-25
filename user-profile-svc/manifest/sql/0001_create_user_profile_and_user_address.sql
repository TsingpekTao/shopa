-- ============================================================================
-- user-profile-svc: schema rebuilt from current userprofile.proto (v1)
-- Notes:
-- 1) default_address_id is kept as a compatibility projection field.
-- 2) address-book default uniqueness is enforced by generated column default_slot.
-- 3) address list business cap (<=20) is enforced in application logic, not SQL.
-- ============================================================================

CREATE TABLE IF NOT EXISTS `user_profile` (
  `user_id`              BIGINT UNSIGNED NOT NULL COMMENT 'User ID from iam-svc',
  `display_name`         VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Display name',
  `avatar_asset_id`      BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Avatar asset id from media-svc',
  `avatar_url`           VARCHAR(512) NOT NULL DEFAULT '' COMMENT 'Avatar URL cache',
  `gender`               TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=unspecified,1=male,2=female,3=other',
  `birthday`             DATE NULL COMMENT 'Birthday date',
  `locale`               VARCHAR(16) NOT NULL DEFAULT 'zh-CN' COMMENT 'Locale',
  `timezone`             VARCHAR(64) NOT NULL DEFAULT 'Asia/Shanghai' COMMENT 'Timezone',
  `marketing_opt_in`     TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Marketing opt in',
  `default_address_id`   BIGINT UNSIGNED NULL COMMENT 'Deprecated compatibility field',
  `profile_version`      BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Optimistic version of profile aggregate',
  `address_book_version` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'CAS version of address-book aggregate',
  `display_name_source`  TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 SYSTEM_INIT,2 USER_SET',
  `last_init_event_at`   DATETIME(3) NULL COMMENT 'Last register-init event time',
  `last_init_event_id`   CHAR(36) NOT NULL DEFAULT '' COMMENT 'Last register-init event id',
  `ext`                  JSON NULL COMMENT 'Extension payload',
  `created_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`           DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`),
  KEY `idx_default_address_id` (`default_address_id`),
  KEY `idx_display_name_source` (`display_name_source`),
  KEY `idx_last_init_event_at` (`last_init_event_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='User profile aggregate';

CREATE TABLE IF NOT EXISTS `user_address` (
  `address_id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Address ID',
  `user_id`                   BIGINT UNSIGNED NOT NULL COMMENT 'Owner user ID',
  `status`                    TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1=active,2=deleted,3=replaced',
  `address_version`           BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Optimistic version of address row',
  `replaced_from_address_id`  BIGINT UNSIGNED NULL COMMENT 'Address replaced chain source',
  `label`                     VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'Address label',
  `receiver_name`             VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Receiver name',
  `receiver_phone`            VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'Receiver phone',
  `country_code`              CHAR(2) NOT NULL DEFAULT 'CN' COMMENT 'ISO country code',
  `province_code`             VARCHAR(20) NOT NULL DEFAULT '' COMMENT 'Province code',
  `province_name`             VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Province name',
  `city_code`                 VARCHAR(20) NOT NULL DEFAULT '' COMMENT 'City code',
  `city_name`                 VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'City name',
  `district_code`             VARCHAR(20) NOT NULL DEFAULT '' COMMENT 'District code',
  `district_name`             VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'District name',
  `street`                    VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Street/town',
  `detail`                    VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Detailed address',
  `postal_code`               VARCHAR(20) NOT NULL DEFAULT '' COMMENT 'Postal code',
  `is_default`                TINYINT(1) NULL DEFAULT NULL COMMENT '1 default, NULL non-default',
  `default_slot`              TINYINT GENERATED ALWAYS AS (CASE WHEN (`status` = 1 AND `is_default` = 1) THEN 1 ELSE NULL END) STORED COMMENT 'Default uniqueness slot',
  `latitude`                  DECIMAL(10,7) NULL COMMENT 'Latitude',
  `longitude`                 DECIMAL(10,7) NULL COMMENT 'Longitude',
  `ext`                       JSON NULL COMMENT 'Extension payload',
  `created_at`                DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`                DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`address_id`),
  UNIQUE KEY `uk_user_default_slot` (`user_id`, `default_slot`),
  KEY `idx_user_status` (`user_id`, `status`),
  KEY `idx_user_default` (`user_id`, `is_default`),
  KEY `idx_receiver_phone` (`receiver_phone`),
  KEY `idx_replaced_from_address_id` (`replaced_from_address_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='User shipping addresses';
