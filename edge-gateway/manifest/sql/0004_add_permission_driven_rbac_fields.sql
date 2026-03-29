USE `shopa_edge_gateway`;
SET NAMES utf8mb4;

ALTER TABLE `edge_proxy_route`
  ADD COLUMN IF NOT EXISTS `required_permission_key` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Required permission key' AFTER `remark`,
  ADD COLUMN IF NOT EXISTS `action` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Audit action key' AFTER `required_permission_key`,
  ADD COLUMN IF NOT EXISTS `resource_id_path_key` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'Resource id path param key' AFTER `action`;

ALTER TABLE `edge_request_audit`
  ADD COLUMN IF NOT EXISTS `permission_key` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Permission key used for this request' AFTER `user_agent`,
  ADD COLUMN IF NOT EXISTS `action` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Audit action' AFTER `permission_key`,
  ADD COLUMN IF NOT EXISTS `resource_id` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'Resource id' AFTER `action`;

INSERT INTO `edge_proxy_route`
(`route_code`, `method`, `path_pattern`, `upstream_service`, `upstream_path_template`, `auth_required`, `inject_user_context`, `body_mode`, `timeout_ms`, `rate_limit_rps`, `status`, `remark`, `required_permission_key`, `action`, `resource_id_path_key`)
VALUES
('EGW_ADMIN_SELLER_APPLICATION_LIST', 'GET', '/v1/admin/seller/applications', 'seller_shop', '/v1/admin/seller/applications', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin merchant application list', 'merchant:review:view', 'merchant.review.list', ''),
('EGW_ADMIN_SELLER_APPLICATION_DETAIL', 'GET', '/v1/admin/seller/applications/{applicationNo}', 'seller_shop', '/v1/admin/seller/applications/{applicationNo}', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin merchant application detail', 'merchant:review:view', 'merchant.review.detail', 'applicationNo'),
('EGW_ADMIN_SELLER_APPLICATION_APPROVE', 'POST', '/v1/admin/seller/applications/{applicationNo}/approve', 'seller_shop', '/v1/admin/seller/applications/{applicationNo}/approve', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin approve merchant application', 'merchant:review:approve', 'merchant.review.approve', 'applicationNo'),
('EGW_ADMIN_SELLER_APPLICATION_REJECT', 'POST', '/v1/admin/seller/applications/{applicationNo}/reject', 'seller_shop', '/v1/admin/seller/applications/{applicationNo}/reject', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin reject merchant application', 'merchant:review:reject', 'merchant.review.reject', 'applicationNo'),
('EGW_ADMIN_SHOP_FREEZE', 'POST', '/v1/admin/seller/shops/{shopNo}/freeze', 'seller_shop', '/v1/admin/seller/shops/{shopNo}/freeze', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin freeze shop', 'shop:manage:freeze', 'shop.manage.freeze', 'shopNo'),
('EGW_ADMIN_SHOP_CLOSE', 'POST', '/v1/admin/seller/shops/{shopNo}/close', 'seller_shop', '/v1/admin/seller/shops/{shopNo}/close', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin close shop', 'shop:manage:close', 'shop.manage.close', 'shopNo'),
('EGW_ADMIN_PRODUCT_REVIEW_TASKS', 'GET', '/v1/catalog/admin/review/tasks', 'catalog', '/v1/catalog/admin/review/tasks', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin product review task list', 'product:review:view', 'product.review.list', ''),
('EGW_ADMIN_PRODUCT_REVIEW_DETAIL', 'GET', '/v1/catalog/admin/review/{spu_no}', 'catalog', '/v1/catalog/admin/review/{spu_no}', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin product review detail', 'product:review:view', 'product.review.detail', 'spu_no'),
('EGW_ADMIN_PRODUCT_APPROVE', 'POST', '/v1/catalog/admin/review/approve', 'catalog', '/v1/catalog/admin/review/approve', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin approve product', 'product:review:approve', 'product.review.approve', ''),
('EGW_ADMIN_PRODUCT_REJECT', 'POST', '/v1/catalog/admin/review/reject', 'catalog', '/v1/catalog/admin/review/reject', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin reject product', 'product:review:reject', 'product.review.reject', ''),
('EGW_ADMIN_PRODUCT_FREEZE', 'POST', '/v1/catalog/admin/review/freeze', 'catalog', '/v1/catalog/admin/review/freeze', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin freeze product', 'product:review:freeze', 'product.review.freeze', ''),
('EGW_ADMIN_PRODUCT_UNFREEZE', 'POST', '/v1/catalog/admin/review/unfreeze', 'catalog', '/v1/catalog/admin/review/unfreeze', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin unfreeze product', 'product:review:unfreeze', 'product.review.unfreeze', ''),
('EGW_ADMIN_PRODUCT_FORCE_OFF_SHELF', 'POST', '/v1/catalog/admin/review/force-off-shelf', 'catalog', '/v1/catalog/admin/review/force-off-shelf', 1, 1, 'NORMAL', 15000, 1000, 1, 'Admin force off shelf', 'product:review:force_off_shelf', 'product.review.force_off_shelf', '')
ON DUPLICATE KEY UPDATE
`method` = VALUES(`method`),
`path_pattern` = VALUES(`path_pattern`),
`upstream_service` = VALUES(`upstream_service`),
`upstream_path_template` = VALUES(`upstream_path_template`),
`auth_required` = VALUES(`auth_required`),
`inject_user_context` = VALUES(`inject_user_context`),
`body_mode` = VALUES(`body_mode`),
`timeout_ms` = VALUES(`timeout_ms`),
`rate_limit_rps` = VALUES(`rate_limit_rps`),
`status` = VALUES(`status`),
`remark` = VALUES(`remark`),
`required_permission_key` = VALUES(`required_permission_key`),
`action` = VALUES(`action`),
`resource_id_path_key` = VALUES(`resource_id_path_key`);
