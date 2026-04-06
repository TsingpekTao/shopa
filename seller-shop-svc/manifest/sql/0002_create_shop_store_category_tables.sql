CREATE TABLE IF NOT EXISTS `shop_store_category` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `shop_no` VARCHAR(40) NOT NULL,
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `name` VARCHAR(64) NOT NULL,
  `level` TINYINT UNSIGNED NOT NULL,
  `sort_order` INT NOT NULL DEFAULT 0,
  `is_visible` TINYINT(1) NOT NULL DEFAULT 1,
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0,
  `product_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_shop_parent_visible_sort` (`shop_no`, `parent_id`, `is_deleted`, `is_visible`, `sort_order`, `id`),
  KEY `idx_shop_level_deleted` (`shop_no`, `level`, `is_deleted`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Seller-defined in-shop category tree';

CREATE TABLE IF NOT EXISTS `shop_store_category_product` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `shop_no` VARCHAR(40) NOT NULL,
  `spu_no` VARCHAR(64) NOT NULL,
  `store_category_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shop_spu` (`shop_no`, `spu_no`),
  KEY `idx_shop_store_category` (`shop_no`, `store_category_id`),
  KEY `idx_store_category_id` (`store_category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SPU to seller-defined in-shop category binding';
