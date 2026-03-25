-- ============================================================================
-- iam-svc: core account/auth schema
-- ============================================================================

CREATE TABLE IF NOT EXISTS `iam_user_auth` (
  `user_id`                 BIGINT UNSIGNED NOT NULL COMMENT 'Snowflake user id',
  `phone`                   VARCHAR(20) NULL COMMENT 'Phone number',
  `email`                   VARCHAR(128) NULL COMMENT 'Email address',
  `password_hash`           VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Password hash',
  `password_salt`           VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Password salt',
  `password_algo`           VARCHAR(32) NOT NULL DEFAULT 'argon2id' COMMENT 'Password hash algorithm',
  `password_ver`            INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Password algorithm version',
  `account_status`          TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 active,2 locked,3 disabled,4 review',
  `failed_login_count`      INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Fail count for audit',
  `locked_until`            DATETIME(3) NULL COMMENT 'Account lock expiry time',
  `token_version`           INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Token version for global token invalidation',
  `last_login_at`           DATETIME(3) NULL COMMENT 'Last login time',
  `last_login_ip`           VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Last login ip',
  `last_login_geo`          VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Last login geo',
  `last_login_ua`           VARCHAR(512) NOT NULL DEFAULT '' COMMENT 'Last login user-agent',
  `last_login_fingerprint`  VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Last login fingerprint',
  `created_at`              DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`              DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uk_phone` (`phone`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status` (`account_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM account table';

CREATE TABLE IF NOT EXISTS `iam_user_role` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`     BIGINT UNSIGNED NOT NULL,
  `role_code`   TINYINT UNSIGNED NOT NULL COMMENT '1 customer,2 seller,3 admin,4 cs',
  `scope_type`  TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 global,2 shop',
  `scope_id`    BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `status`      TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 active,2 disabled',
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role_scope` (`user_id`, `role_code`, `scope_type`, `scope_id`),
  KEY `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM user roles';

CREATE TABLE IF NOT EXISTS `iam_user_oauth` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`       BIGINT UNSIGNED NOT NULL,
  `provider`      TINYINT UNSIGNED NOT NULL COMMENT '1 wechat,2 alipay',
  `provider_uid`  VARCHAR(128) NOT NULL DEFAULT '',
  `union_id`      VARCHAR(128) NOT NULL DEFAULT '',
  `status`        TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 active,2 unbound',
  `created_at`    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_provider_uid` (`provider`, `provider_uid`),
  KEY `idx_user_provider` (`user_id`, `provider`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM third-party oauth links';

CREATE TABLE IF NOT EXISTS `iam_membership` (
  `user_id`     BIGINT UNSIGNED NOT NULL,
  `level_code`  VARCHAR(32) NOT NULL DEFAULT 'BASIC',
  `points`      BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `expire_at`   DATETIME(3) NULL,
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM membership summary';

CREATE TABLE IF NOT EXISTS `iam_refresh_session` (
  `sid`                 VARCHAR(64) NOT NULL COMMENT 'Session id',
  `user_id`             BIGINT UNSIGNED NOT NULL,
  `refresh_token_hash`  CHAR(64) NOT NULL DEFAULT '',
  `ua_hash`             CHAR(64) NOT NULL DEFAULT '',
  `ip`                  VARCHAR(64) NOT NULL DEFAULT '',
  `geo`                 VARCHAR(128) NOT NULL DEFAULT '',
  `fingerprint`         VARCHAR(128) NOT NULL DEFAULT '',
  `created_at`          DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `expires_at`          DATETIME(3) NOT NULL,
  `revoked_at`          DATETIME(3) NULL,
  `replaced_by_sid`     VARCHAR(64) NOT NULL DEFAULT '',
  PRIMARY KEY (`sid`),
  KEY `idx_user_revoked_exp` (`user_id`, `revoked_at`, `expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM refresh token sessions';

CREATE TABLE IF NOT EXISTS `iam_login_log` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`      BIGINT UNSIGNED NULL,
  `identifier`   VARCHAR(128) NOT NULL DEFAULT '',
  `channel`      TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '1 password,2 sms,3 oauth',
  `success`      TINYINT(1) NOT NULL DEFAULT 0,
  `fail_reason`  VARCHAR(64) NOT NULL DEFAULT '',
  `ip`           VARCHAR(64) NOT NULL DEFAULT '',
  `geo`          VARCHAR(128) NOT NULL DEFAULT '',
  `ua`           VARCHAR(512) NOT NULL DEFAULT '',
  `fingerprint`  VARCHAR(128) NOT NULL DEFAULT '',
  `created_at`   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_user_created` (`user_id`, `created_at`),
  KEY `idx_identifier_created` (`identifier`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM login audit logs';

CREATE TABLE IF NOT EXISTS `iam_sms_log` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `scene`        TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '1 register,2 login,3 mfa,4 reset_password',
  `target`       VARCHAR(32) NOT NULL DEFAULT '',
  `provider`     VARCHAR(32) NOT NULL DEFAULT 'mock',
  `biz_id`       VARCHAR(128) NOT NULL DEFAULT '',
  `ip`           VARCHAR(64) NOT NULL DEFAULT '',
  `ua`           VARCHAR(512) NOT NULL DEFAULT '',
  `fingerprint`  VARCHAR(128) NOT NULL DEFAULT '',
  `success`      TINYINT(1) NOT NULL DEFAULT 0,
  `error_code`   INT NOT NULL DEFAULT 0,
  `created_at`   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_target_created` (`target`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM sms sending logs';

CREATE TABLE IF NOT EXISTS `iam_outbox_event` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_type`    VARCHAR(64) NOT NULL DEFAULT '',
  `payload_json`  JSON NOT NULL,
  `status`        TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 new,2 sent,3 failed',
  `retry_count`   INT UNSIGNED NOT NULL DEFAULT 0,
  `next_retry_at` DATETIME(3) NULL,
  `created_at`    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_status_retry` (`status`, `next_retry_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='IAM outbox events';

