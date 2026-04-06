USE shopa_search_svc;

SET @stmt = IF(
  EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'search_spu_doc' AND COLUMN_NAME = 'source_version'
  ),
  'SELECT 1',
  'ALTER TABLE search_spu_doc ADD COLUMN source_version BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER attrs_json'
);
PREPARE add_source_version FROM @stmt;
EXECUTE add_source_version;
DEALLOCATE PREPARE add_source_version;

SET @stmt = IF(
  EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'search_spu_doc' AND COLUMN_NAME = 'source_updated_at'
  ),
  'SELECT 1',
  'ALTER TABLE search_spu_doc ADD COLUMN source_updated_at DATETIME(3) NULL AFTER source_version'
);
PREPARE add_source_updated_at FROM @stmt;
EXECUTE add_source_updated_at;
DEALLOCATE PREPARE add_source_updated_at;

SET @stmt = IF(
  EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'search_spu_doc' AND COLUMN_NAME = 'deleted'
  ),
  'SELECT 1',
  'ALTER TABLE search_spu_doc ADD COLUMN deleted TINYINT UNSIGNED NOT NULL DEFAULT 0 AFTER source_updated_at'
);
PREPARE add_deleted FROM @stmt;
EXECUTE add_deleted;
DEALLOCATE PREPARE add_deleted;

SET @stmt = IF(
  EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'search_spu_doc' AND COLUMN_NAME = 'deleted_at'
  ),
  'SELECT 1',
  'ALTER TABLE search_spu_doc ADD COLUMN deleted_at DATETIME(3) NULL AFTER deleted'
);
PREPARE add_deleted_at FROM @stmt;
EXECUTE add_deleted_at;
DEALLOCATE PREPARE add_deleted_at;

SET @stmt = IF(
  EXISTS (
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'search_spu_doc' AND INDEX_NAME = 'idx_deleted_at'
  ),
  'SELECT 1',
  'ALTER TABLE search_spu_doc ADD INDEX idx_deleted_at (deleted, deleted_at)'
);
PREPARE add_idx_deleted_at FROM @stmt;
EXECUTE add_idx_deleted_at;
DEALLOCATE PREPARE add_idx_deleted_at;

ALTER TABLE search_spu_doc
  MODIFY COLUMN attrs_json JSON NULL,
  MODIFY COLUMN updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  MODIFY COLUMN created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

ALTER TABLE search_keyword_stat
  MODIFY COLUMN last_searched_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  MODIFY COLUMN created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  MODIFY COLUMN updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);

UPDATE search_spu_doc
SET
  source_version = CASE
    WHEN source_version > 0 THEN source_version
    WHEN updated_at IS NOT NULL THEN UNIX_TIMESTAMP(updated_at) * 1000
    ELSE 0
  END,
  source_updated_at = COALESCE(source_updated_at, updated_at),
  deleted = COALESCE(deleted, 0)
WHERE 1 = 1;
