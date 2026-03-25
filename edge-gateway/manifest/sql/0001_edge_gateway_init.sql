CREATE DATABASE IF NOT EXISTS `shopa_edge_gateway`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE `shopa_edge_gateway`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `edge_proxy_route` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `route_code` VARCHAR(64) NOT NULL COMMENT '路由编码（全局唯一）',
  `method` VARCHAR(16) NOT NULL COMMENT 'HTTP方法',
  `path_pattern` VARCHAR(255) NOT NULL COMMENT '网关路由匹配路径',
  `upstream_service` VARCHAR(64) NOT NULL COMMENT '目标上游服务标识',
  `upstream_path_template` VARCHAR(255) NOT NULL COMMENT '目标上游路径模板',
  `auth_required` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否需要鉴权',
  `inject_user_context` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否注入用户上下文头',
  `body_mode` ENUM('STREAM','NORMAL') NOT NULL DEFAULT 'NORMAL' COMMENT '请求体处理模式',
  `timeout_ms` INT UNSIGNED NOT NULL DEFAULT 15000 COMMENT '上游调用超时时间（毫秒）',
  `rate_limit_rps` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '每秒限流阈值，0表示不限流',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：1启用 0停用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_route_code` (`route_code`),
  KEY `idx_method_path_status` (`method`, `path_pattern`, `status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='网关透传路由策略表';

CREATE TABLE IF NOT EXISTS `edge_request_audit` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `request_id` VARCHAR(64) NOT NULL COMMENT '请求ID（全局唯一）',
  `route_code` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '命中路由编码',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '调用用户ID',
  `method` VARCHAR(16) NOT NULL COMMENT 'HTTP方法',
  `path` VARCHAR(255) NOT NULL COMMENT '请求路径',
  `upstream_service` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '上游服务标识',
  `status_code` INT NOT NULL DEFAULT 0 COMMENT 'HTTP状态码',
  `latency_ms` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '请求耗时毫秒',
  `partial` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否降级返回',
  `degraded_fields_json` JSON NULL COMMENT '降级字段列表',
  `error_code` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '业务错误码',
  `client_ip` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '客户端IP',
  `user_agent` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '客户端UA',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_request_id` (`request_id`),
  KEY `idx_route_created_at` (`route_code`, `created_at`),
  KEY `idx_user_created_at` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='网关请求审计表';

INSERT INTO `edge_proxy_route`
(`route_code`, `method`, `path_pattern`, `upstream_service`, `upstream_path_template`, `auth_required`, `inject_user_context`, `body_mode`, `timeout_ms`, `rate_limit_rps`, `status`, `remark`)
VALUES
('EGW_CATALOG_CREATE_DRAFT', 'POST', '/v1/catalog/seller/products/draft', 'catalog', '/v1/catalog/seller/products/draft', 1, 1, 'NORMAL', 15000, 1000, 1, '商家创建商品草稿'),
('EGW_CATALOG_UPSERT_SKU', 'POST', '/v1/catalog/seller/products/skus:upsert', 'catalog', '/v1/catalog/seller/products/skus:upsert', 1, 1, 'NORMAL', 15000, 1000, 1, '商家批量更新SKU草稿'),
('EGW_MEDIA_UPLOAD', 'POST', '/v1/media/upload', 'media', '/v1/media/upload', 1, 1, 'STREAM', 600000, 200, 1, '媒体上传流式透传')
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
`remark` = VALUES(`remark`);
