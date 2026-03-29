CREATE TABLE IF NOT EXISTS `catalog_spu_image` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `spu_no` VARCHAR(32) NOT NULL,
  `image_url` VARCHAR(2048) NOT NULL,
  `source` VARCHAR(32) NOT NULL DEFAULT 'oss',
  `is_primary` TINYINT(1) NOT NULL DEFAULT 1,
  `sort_order` INT NOT NULL DEFAULT 1,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '1 ENABLED, 2 DISABLED',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_spu_sort` (`spu_no`, `sort_order`),
  KEY `idx_spu_primary` (`spu_no`, `is_primary`, `status`),
  KEY `idx_status_updated` (`status`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SPU image url mapping';

SET NAMES utf8mb4;
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2001','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2001.jpg?Expires=1932449596&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=poGpNh8l8i3QjbCs%2B9SnXEVfgTk%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2002','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2002.jpg?Expires=1932449596&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=C30vZOLWA6Q6JSJGnNY5rCwYHyE%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2003','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2003.jpg?Expires=1932449597&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=Clb9FzOKM9NOUR7TKUWPBGBHH8E%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2004','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2004.jpg?Expires=1932449597&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=RCNjDQG0WzkolGZ3aLbmS%2F%2Fw22g%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2005','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2005.jpg?Expires=1932449597&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=S4s7c9ADHfKdo5fgCSMMb8EdpdM%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2006','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2006.jpg?Expires=1932449600&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=xrzwkILPRWIlcFLpH9PZFwAubA8%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2007','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2007.jpg?Expires=1932449601&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=tsEJ0FUZgGie9vecXv1AssWZKyw%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2008','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2008.jpg?Expires=1932449602&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=CNNmyOw35pbDdQsppcO0QZBOkwk%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2009','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2009.jpg?Expires=1932449602&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=bH5Z39AI36qbCNefGTb5Dx7kU4c%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2010','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2010.jpg?Expires=1932449603&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=5XNaoDM4oo8yBut3dC87YuxDqb0%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2011','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2011.jpg?Expires=1932449603&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=A5Koe6SICPIoi6hcgyTXKONUrOc%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2012','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2012.jpg?Expires=1932449603&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=%2FfiU0uGah1V0OPwBnz7ZBPP5VF4%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2013','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2013.jpg?Expires=1932449603&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=HVWzF%2Bb6YIl%2FjO3zCYIAO5QdkHs%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2014','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2014.jpg?Expires=1932449603&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=4343ldnupU66slZCEaecFmGlYJs%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);
INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)
VALUES ('P2015','https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/mall%2Fproducts%2F20260329%2Fp2015.jpg?Expires=1932449603&OSSAccessKeyId=LTAI5tPNdFg9D1vBc8bCBTQ8&Signature=QTngvMgOd9dmqldDF0Dy78XBUok%3D','oss',1,1,1,NOW(3),NOW(3))
ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);

