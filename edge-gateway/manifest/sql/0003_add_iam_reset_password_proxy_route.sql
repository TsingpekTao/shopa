USE `shopa_edge_gateway`;
SET NAMES utf8mb4;

INSERT INTO `edge_proxy_route`
(`route_code`, `method`, `path_pattern`, `upstream_service`, `upstream_path_template`, `auth_required`, `inject_user_context`, `body_mode`, `timeout_ms`, `rate_limit_rps`, `status`, `remark`)
VALUES
('EGW_IAM_RESET_PASSWORD_SMS', 'POST', '/v1/auth/password/reset/sms', 'iam', '/v1/auth/password/reset/sms', 0, 0, 'NORMAL', 15000, 1000, 1, '未登录短信重置密码')
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
