ALTER TABLE `refund_task`
  ADD COLUMN `points_return_amount` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '本次退款返还的已消费积分，单位：分' AFTER `refund_amount`,
  ADD COLUMN `points_reverse_amount` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '本次退款冲回的已赠积分，单位：分' AFTER `points_return_amount`,
  ADD COLUMN `points_cash_offset_amount` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '赠分冲回余额不足时的现金抵扣金额，单位：分' AFTER `points_reverse_amount`,
  ADD COLUMN `final_cash_refund_amount` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '实际现金退款金额，单位：分' AFTER `points_cash_offset_amount`,
  ADD COLUMN `account_debt_after` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '退款执行后积分账户欠账绝对值，单位：分' AFTER `final_cash_refund_amount`;
