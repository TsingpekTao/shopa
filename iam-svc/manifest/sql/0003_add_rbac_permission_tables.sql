-- ============================================================================
-- iam-svc: RBAC permission model (permission-point driven)
-- ============================================================================

ALTER TABLE `iam_user_role`
  MODIFY COLUMN `role_code` TINYINT UNSIGNED NOT NULL COMMENT '1 customer,2 seller,3 admin(legacy),4 cs,5 super_admin,6 auditor,7 ops_analyst';

CREATE TABLE IF NOT EXISTS `iam_role` (
  `role_code`   TINYINT UNSIGNED NOT NULL COMMENT 'role code from enum',
  `role_key`    VARCHAR(64) NOT NULL COMMENT 'stable role key',
  `role_name`   VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'display role name',
  `status`      TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 active,2 disabled',
  `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`role_code`),
  UNIQUE KEY `uk_role_key` (`role_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='RBAC role dictionary';

CREATE TABLE IF NOT EXISTS `iam_permission` (
  `permission_key`  VARCHAR(128) NOT NULL COMMENT 'resource:domain:action',
  `permission_name` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'display permission name',
  `status`          TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 active,2 disabled',
  `created_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`permission_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='RBAC permission dictionary';

CREATE TABLE IF NOT EXISTS `iam_role_permission` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `role_code`       TINYINT UNSIGNED NOT NULL,
  `permission_key`  VARCHAR(128) NOT NULL,
  `status`          TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 active,2 disabled',
  `created_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at`      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_code`, `permission_key`),
  KEY `idx_permission_status` (`permission_key`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='RBAC role-permission mapping';

INSERT INTO `iam_role` (`role_code`, `role_key`, `role_name`, `status`) VALUES
  (3, 'ADMIN', 'Legacy Admin', 1),
  (4, 'CUSTOMER_SERVICE', 'Customer Service', 1),
  (5, 'SUPER_ADMIN', 'Super Admin', 1),
  (6, 'AUDITOR', 'Auditor', 1),
  (7, 'OPS_ANALYST', 'Operations Analyst', 1)
ON DUPLICATE KEY UPDATE
  `role_key` = VALUES(`role_key`),
  `role_name` = VALUES(`role_name`),
  `status` = VALUES(`status`);

INSERT INTO `iam_permission` (`permission_key`, `permission_name`, `status`) VALUES
  ('*:*:*', 'All Permissions', 1),
  ('admin:me:overview', 'Admin Me Overview', 1),
  ('merchant:review:view', 'View Merchant Review Tasks', 1),
  ('merchant:review:approve', 'Approve Merchant Applications', 1),
  ('merchant:review:reject', 'Reject Merchant Applications', 1),
  ('shop:manage:freeze', 'Freeze Shop', 1),
  ('shop:manage:close', 'Close Shop', 1),
  ('product:review:view', 'View Product Review Tasks', 1),
  ('product:review:approve', 'Approve Product', 1),
  ('product:review:reject', 'Reject Product', 1),
  ('product:review:freeze', 'Freeze Product', 1),
  ('product:review:unfreeze', 'Unfreeze Product', 1),
  ('product:review:force_off_shelf', 'Force Off Shelf Product', 1),
  ('dashboard:mall:view', 'View Mall Dashboard', 1),
  ('dashboard:shop:view', 'View Shop Insights', 1),
  ('cs:conversation:view', 'View Customer Service Conversation', 1),
  ('cs:conversation:reply', 'Reply Customer Service Conversation', 1)
ON DUPLICATE KEY UPDATE
  `permission_name` = VALUES(`permission_name`),
  `status` = VALUES(`status`);

INSERT INTO `iam_role_permission` (`role_code`, `permission_key`, `status`) VALUES
  (3, '*:*:*', 1),
  (4, 'admin:me:overview', 1),
  (4, 'cs:conversation:view', 1),
  (4, 'cs:conversation:reply', 1),
  (5, '*:*:*', 1),
  (6, 'admin:me:overview', 1),
  (6, 'merchant:review:view', 1),
  (6, 'merchant:review:approve', 1),
  (6, 'merchant:review:reject', 1),
  (6, 'product:review:view', 1),
  (6, 'product:review:approve', 1),
  (6, 'product:review:reject', 1),
  (7, 'admin:me:overview', 1),
  (7, 'dashboard:mall:view', 1),
  (7, 'dashboard:shop:view', 1),
  (7, 'merchant:review:view', 1),
  (7, 'product:review:view', 1)
ON DUPLICATE KEY UPDATE
  `status` = VALUES(`status`);
