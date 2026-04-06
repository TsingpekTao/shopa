ALTER TABLE `search_spu_doc`
  ADD COLUMN `store_category_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `category_no`,
  ADD COLUMN `store_category_l1` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `store_category_id`,
  ADD COLUMN `store_category_l2` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `store_category_l1`,
  ADD COLUMN `store_category_path_json` JSON NULL AFTER `store_category_l2`,
  ADD KEY `idx_store_category_id` (`store_category_id`),
  ADD KEY `idx_store_category_l1` (`store_category_l1`),
  ADD KEY `idx_store_category_l2` (`store_category_l2`);
