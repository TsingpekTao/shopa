USE `shopa_edge_gateway`;
SET NAMES utf8mb4;

INSERT INTO `edge_proxy_route`
(`route_code`, `method`, `path_pattern`, `upstream_service`, `upstream_path_template`, `auth_required`, `inject_user_context`, `body_mode`, `timeout_ms`, `rate_limit_rps`, `status`, `remark`)
VALUES
('EGW_CART_ADD_ITEM', 'POST', '/v1/cart/items:add', 'cart', '/v1/cart/items:add', 1, 1, 'NORMAL', 15000, 1000, 1, '购物车加购'),
('EGW_CART_UPDATE_QTY', 'POST', '/v1/cart/items:qty', 'cart', '/v1/cart/items:qty', 1, 1, 'NORMAL', 15000, 1000, 1, '购物车改数量'),
('EGW_CART_TOGGLE_CHECKED', 'POST', '/v1/cart/items:check', 'cart', '/v1/cart/items:check', 1, 1, 'NORMAL', 15000, 1000, 1, '购物车勾选切换'),
('EGW_CART_BATCH_CHECKED', 'POST', '/v1/cart/items:batch-check', 'cart', '/v1/cart/items:batch-check', 1, 1, 'NORMAL', 15000, 1000, 1, '购物车批量勾选切换'),
('EGW_CART_REMOVE_ITEMS', 'POST', '/v1/cart/items:remove', 'cart', '/v1/cart/items:remove', 1, 1, 'NORMAL', 15000, 1000, 1, '购物车移除项目'),
('EGW_CART_CLEAR_INVALID', 'POST', '/v1/cart/items:clear-invalid', 'cart', '/v1/cart/items:clear-invalid', 1, 1, 'NORMAL', 15000, 1000, 1, '清理失效购物车项目'),
('EGW_CART_GET_MY_CART', 'GET', '/v1/cart/me', 'cart', '/v1/cart/me', 1, 1, 'NORMAL', 15000, 1000, 1, '查询我的购物车'),
('EGW_CART_PREPARE_CHECKOUT', 'POST', '/v1/cart/checkout:prepare', 'cart', '/v1/cart/checkout:prepare', 1, 1, 'NORMAL', 15000, 1000, 1, '准备结算快照'),
('EGW_CART_INTERNAL_CONSUME_CHECKOUT', 'POST', '/v1/cart/internal/checkout:consume', 'cart', '/v1/cart/internal/checkout:consume', 1, 1, 'NORMAL', 15000, 1000, 1, '内部消费结算token'),
('EGW_CART_INTERNAL_MARK_ORDERED', 'POST', '/v1/cart/internal/items:ordered', 'cart', '/v1/cart/internal/items:ordered', 1, 1, 'NORMAL', 15000, 1000, 1, '内部下单后清理购物车'),
('EGW_CART_INTERNAL_UPSERT_SKU_PROJECTION', 'POST', '/v1/cart/internal/sku-projection:batch-upsert', 'cart', '/v1/cart/internal/sku-projection:batch-upsert', 1, 1, 'NORMAL', 15000, 1000, 1, '内部批量刷新SKU投影'),
('EGW_ORDER_CREATE_FROM_CART', 'POST', '/v1/order/buyer/orders:create-from-cart', 'order', '/v1/order/buyer/orders:create-from-cart', 1, 1, 'NORMAL', 15000, 1000, 1, '买家购物车下单'),
('EGW_ORDER_CREATE_BUY_NOW', 'POST', '/v1/order/buyer/orders:create-buy-now', 'order', '/v1/order/buyer/orders:create-buy-now', 1, 1, 'NORMAL', 15000, 1000, 1, '买家立即购买下单'),
('EGW_ORDER_REQUEST_PAY', 'POST', '/v1/order/buyer/orders:request-pay', 'order', '/v1/order/buyer/orders:request-pay', 1, 1, 'NORMAL', 15000, 1000, 1, '买家发起支付'),
('EGW_ORDER_CANCEL_MY_ORDER', 'POST', '/v1/order/buyer/orders:cancel', 'order', '/v1/order/buyer/orders:cancel', 1, 1, 'NORMAL', 15000, 1000, 1, '买家取消订单'),
('EGW_ORDER_GET_MY_ORDER_DETAIL', 'GET', '/v1/order/buyer/orders/{order_no}', 'order', '/v1/order/buyer/orders/{order_no}', 1, 1, 'NORMAL', 15000, 1000, 1, '买家查询订单详情'),
('EGW_ORDER_LIST_MY_ORDERS', 'GET', '/v1/order/buyer/orders', 'order', '/v1/order/buyer/orders', 1, 1, 'NORMAL', 15000, 1000, 1, '买家订单列表'),
('EGW_ORDER_LIST_SHOP_ORDERS', 'GET', '/v1/order/seller/shops/{shop_no}/orders', 'order', '/v1/order/seller/shops/{shop_no}/orders', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家店铺订单列表'),
('EGW_ORDER_GET_SHOP_ORDER_DETAIL', 'GET', '/v1/order/seller/orders/{sub_order_no}', 'order', '/v1/order/seller/orders/{sub_order_no}', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家子单详情'),
('EGW_ORDER_MARK_SUB_ORDER_SHIPPED', 'POST', '/v1/order/seller/orders:ship', 'order', '/v1/order/seller/orders:ship', 1, 1, 'NORMAL', 15000, 1000, 1, '卖家发货'),
('EGW_ORDER_INTERNAL_CLOSE_UNPAID', 'POST', '/v1/order/internal/orders:close-if-unpaid', 'order', '/v1/order/internal/orders:close-if-unpaid', 1, 1, 'NORMAL', 15000, 1000, 1, '内部关闭超时未支付订单'),
('EGW_ORDER_INTERNAL_PAY_CALLBACK', 'POST', '/v1/order/internal/payments:callback', 'order', '/v1/order/internal/payments:callback', 1, 1, 'NORMAL', 15000, 1000, 1, '内部处理支付回调'),
('EGW_ORDER_INTERNAL_SNAPSHOT', 'GET', '/v1/order/internal/orders/{order_no}/snapshot', 'order', '/v1/order/internal/orders/{order_no}/snapshot', 1, 1, 'NORMAL', 15000, 1000, 1, '内部查询订单快照')
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