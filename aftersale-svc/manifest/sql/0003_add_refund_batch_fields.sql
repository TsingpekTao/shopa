ALTER TABLE `after_sale_case`
  ADD COLUMN `refund_batch_no` VARCHAR(64) NOT NULL DEFAULT '' AFTER `cancel_reason_code`,
  ADD COLUMN `scope_code` VARCHAR(32) NOT NULL DEFAULT '' AFTER `refund_batch_no`,
  ADD COLUMN `review_deadline_at` DATETIME(3) NULL AFTER `scope_code`,
  ADD COLUMN `auto_approved_at` DATETIME(3) NULL AFTER `review_deadline_at`,
  ADD COLUMN `selected_item_nos_json` JSON NULL AFTER `auto_approved_at`,
  ADD COLUMN `payment_no` VARCHAR(64) NOT NULL DEFAULT '' AFTER `selected_item_nos_json`,
  ADD KEY `idx_refund_batch_no` (`refund_batch_no`),
  ADD KEY `idx_status_review_deadline` (`after_sale_status`,`review_deadline_at`);
