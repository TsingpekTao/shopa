USE `shopa_edge_gateway`;
SET NAMES utf8mb4;

INSERT INTO `edge_proxy_route`
(`route_code`, `method`, `path_pattern`, `upstream_service`, `upstream_path_template`, `auth_required`, `inject_user_context`, `body_mode`, `timeout_ms`, `rate_limit_rps`, `status`, `remark`)
VALUES
('EGW_AFTERSALE_BUYER_CREATE', 'POST', '/v1/aftersale/buyer/cases:create', 'aftersale', '/v1/aftersale/buyer/cases:create', 1, 1, 'NORMAL', 15000, 1000, 1, '买家创建售后'),
('EGW_AFTERSALE_BUYER_CANCEL', 'POST', '/v1/aftersale/buyer/cases:cancel', 'aftersale', '/v1/aftersale/buyer/cases:cancel', 1, 1, 'NORMAL', 15000, 1000, 1, '买家撤销售后'),
('EGW_AFTERSALE_BUYER_DETAIL', 'GET', '/v1/aftersale/buyer/cases/{after_sale_no}', 'aftersale', '/v1/aftersale/buyer/cases/{after_sale_no}', 1, 1, 'NORMAL', 15000, 1000, 1, '买家售后详情'),
('EGW_AFTERSALE_BUYER_LIST', 'GET', '/v1/aftersale/buyer/cases', 'aftersale', '/v1/aftersale/buyer/cases', 1, 1, 'NORMAL', 15000, 1000, 1, '买家售后列表'),
('EGW_AFTERSALE_SELLER_LIST', 'GET', '/v1/aftersale/seller/shops/{shop_no}/cases', 'aftersale', '/v1/aftersale/seller/shops/{shop_no}/cases', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家售后列表'),
('EGW_AFTERSALE_SELLER_DETAIL', 'GET', '/v1/aftersale/seller/cases/{after_sale_no}', 'aftersale', '/v1/aftersale/seller/cases/{after_sale_no}', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家售后详情'),
('EGW_AFTERSALE_SELLER_APPROVE', 'POST', '/v1/aftersale/seller/cases:approve', 'aftersale', '/v1/aftersale/seller/cases:approve', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家通过售后'),
('EGW_AFTERSALE_SELLER_REJECT', 'POST', '/v1/aftersale/seller/cases:reject', 'aftersale', '/v1/aftersale/seller/cases:reject', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家拒绝售后'),

('EGW_REVIEW_BUYER_CREATE', 'POST', '/v1/review/buyer/reviews:create', 'review', '/v1/review/buyer/reviews:create', 1, 1, 'NORMAL', 15000, 1000, 1, '买家发布评价'),
('EGW_REVIEW_BUYER_APPEND', 'POST', '/v1/review/buyer/reviews:append', 'review', '/v1/review/buyer/reviews:append', 1, 1, 'NORMAL', 15000, 1000, 1, '买家追评'),
('EGW_REVIEW_BUYER_LIST', 'GET', '/v1/review/buyer/reviews', 'review', '/v1/review/buyer/reviews', 1, 1, 'NORMAL', 15000, 1000, 1, '买家评价列表'),
('EGW_REVIEW_SELLER_REPLY', 'POST', '/v1/review/seller/reviews:reply', 'review', '/v1/review/seller/reviews:reply', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家回复评价'),
('EGW_REVIEW_PUBLIC_LIST', 'GET', '/v1/review/public/spu/{spu_no}/reviews', 'review', '/v1/review/public/spu/{spu_no}/reviews', 0, 0, 'NORMAL', 15000, 1000, 1, '商品评价列表'),
('EGW_REVIEW_PUBLIC_SUMMARY', 'GET', '/v1/review/public/spu/{spu_no}/summary', 'review', '/v1/review/public/spu/{spu_no}/summary', 0, 0, 'NORMAL', 15000, 1000, 1, '商品评分摘要'),

('EGW_FULFILLMENT_SELLER_CREATE', 'POST', '/v1/fulfillment/seller/shipments:create', 'fulfillment', '/v1/fulfillment/seller/shipments:create', 1, 1, 'NORMAL', 15000, 1000, 1, '创建履约单'),
('EGW_FULFILLMENT_SELLER_SHIP', 'POST', '/v1/fulfillment/seller/shipments:ship', 'fulfillment', '/v1/fulfillment/seller/shipments:ship', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家发货'),
('EGW_FULFILLMENT_SELLER_LIST', 'GET', '/v1/fulfillment/seller/shops/{shop_no}/shipments', 'fulfillment', '/v1/fulfillment/seller/shops/{shop_no}/shipments', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家履约列表'),
('EGW_FULFILLMENT_SELLER_DETAIL', 'GET', '/v1/fulfillment/seller/shipments/{shipment_no}', 'fulfillment', '/v1/fulfillment/seller/shipments/{shipment_no}', 1, 1, 'NORMAL', 15000, 1000, 1, '履约详情'),
('EGW_FULFILLMENT_BUYER_LOGISTICS', 'GET', '/v1/fulfillment/buyer/orders/{order_no}/logistics', 'fulfillment', '/v1/fulfillment/buyer/orders/{order_no}/logistics', 1, 1, 'NORMAL', 15000, 1000, 1, '买家查看物流'),

('EGW_CHAT_BUYER_CREATE_CONVERSATION', 'POST', '/v1/chat/buyer/conversations:get-or-create', 'chat', '/v1/chat/buyer/conversations:get-or-create', 1, 1, 'NORMAL', 15000, 1000, 1, '买家创建会话'),
('EGW_CHAT_BUYER_SEND', 'POST', '/v1/chat/buyer/messages:send', 'chat', '/v1/chat/buyer/messages:send', 1, 1, 'NORMAL', 15000, 1000, 1, '买家发送消息'),
('EGW_CHAT_BUYER_CONVERSATIONS', 'GET', '/v1/chat/buyer/conversations', 'chat', '/v1/chat/buyer/conversations', 1, 1, 'NORMAL', 15000, 1000, 1, '买家会话列表'),
('EGW_CHAT_BUYER_MESSAGES', 'GET', '/v1/chat/buyer/conversations/{conversation_no}/messages', 'chat', '/v1/chat/buyer/conversations/{conversation_no}/messages', 1, 1, 'NORMAL', 15000, 1000, 1, '买家消息列表'),
('EGW_CHAT_BUYER_MARK_READ', 'POST', '/v1/chat/buyer/conversations:mark-read', 'chat', '/v1/chat/buyer/conversations:mark-read', 1, 1, 'NORMAL', 15000, 1000, 1, '买家已读回执'),
('EGW_CHAT_BUYER_UNREAD_SUMMARY', 'GET', '/v1/chat/buyer/unread-summary', 'chat', '/v1/chat/buyer/unread-summary', 1, 1, 'NORMAL', 15000, 1000, 1, '买家未读汇总'),
('EGW_CHAT_SELLER_CONVERSATIONS', 'GET', '/v1/chat/seller/shops/{shop_no}/conversations', 'chat', '/v1/chat/seller/shops/{shop_no}/conversations', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家会话列表'),
('EGW_CHAT_SELLER_SEND', 'POST', '/v1/chat/seller/messages:send', 'chat', '/v1/chat/seller/messages:send', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家发送消息'),
('EGW_CHAT_SELLER_MARK_READ', 'POST', '/v1/chat/seller/conversations:mark-read', 'chat', '/v1/chat/seller/conversations:mark-read', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家已读回执'),
('EGW_CHAT_INTERNAL_SYSTEM_NOTICE', 'POST', '/v1/chat/internal/system-notices:publish', 'chat', '/v1/chat/internal/system-notices:publish', 1, 1, 'NORMAL', 15000, 1000, 1, '内部系统消息'),
('EGW_CHAT_INTERNAL_SNAPSHOT', 'GET', '/v1/chat/internal/conversations/{conversation_no}/snapshot', 'chat', '/v1/chat/internal/conversations/{conversation_no}/snapshot', 1, 1, 'NORMAL', 15000, 1000, 1, '内部会话快照')
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

