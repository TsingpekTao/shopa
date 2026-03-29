-- ============================================================================
-- iam-svc: batch grant admin roles by phone (template, idempotent)
-- ============================================================================
USE `shopa_iam_svc`;
SET NAMES utf8mb4;

-- Put phones you want to authorize here.
DROP TEMPORARY TABLE IF EXISTS `tmp_admin_grant_phone`;
CREATE TEMPORARY TABLE `tmp_admin_grant_phone` (
  `phone` VARCHAR(20) NOT NULL PRIMARY KEY
);

INSERT INTO `tmp_admin_grant_phone` (`phone`) VALUES
  ('13364027679');
-- ,('13800138000')
-- ,('13900139000');

-- Grant role:
-- 5 = SUPER_ADMIN
-- 6 = AUDITOR
-- 7 = OPS_ANALYST
-- 4 = CUSTOMER_SERVICE
SET @grant_role_code := 6;

INSERT INTO `iam_user_role` (`user_id`, `role_code`, `scope_type`, `scope_id`, `status`)
SELECT a.`user_id`, @grant_role_code, 1, 0, 1
FROM `iam_user_auth` a
JOIN `tmp_admin_grant_phone` t ON t.`phone` = a.`phone`
ON DUPLICATE KEY UPDATE
  `status` = VALUES(`status`);

SELECT
  a.`user_id`,
  a.`phone`,
  GROUP_CONCAT(r.`role_code` ORDER BY r.`role_code`) AS role_codes
FROM `iam_user_auth` a
JOIN `tmp_admin_grant_phone` t ON t.`phone` = a.`phone`
LEFT JOIN `iam_user_role` r
  ON r.`user_id` = a.`user_id`
 AND r.`status` = 1
GROUP BY a.`user_id`, a.`phone`;
