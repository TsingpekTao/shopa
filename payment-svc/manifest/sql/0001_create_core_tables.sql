CREATE TABLE IF NOT EXISTS payment_intent (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  payment_no VARCHAR(64) NOT NULL,
  order_no VARCHAR(64) NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  pay_channel TINYINT UNSIGNED NOT NULL DEFAULT 0,
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  payable_amount BIGINT UNSIGNED NOT NULL DEFAULT 0,
  refunded_amount BIGINT UNSIGNED NOT NULL DEFAULT 0,
  currency_code VARCHAR(16) NOT NULL DEFAULT 'CNY',
  order_expire_at DATETIME NULL,
  gateway_expire_at DATETIME NULL,
  external_trade_no VARCHAR(128) NOT NULL DEFAULT '',
  paid_at DATETIME NULL,
  closed_at DATETIME NULL,
  version BIGINT UNSIGNED NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uniq_payment_no (payment_no),
  UNIQUE KEY uniq_order_no (order_no),
  KEY idx_status_updated (status, updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS payment_callback_log (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  callback_event_id VARCHAR(128) NOT NULL,
  payment_no VARCHAR(64) NOT NULL,
  order_no VARCHAR(64) NOT NULL,
  gateway_status_code VARCHAR(64) NOT NULL DEFAULT '',
  raw_payload MEDIUMTEXT,
  processed TINYINT UNSIGNED NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uniq_callback_event_id (callback_event_id),
  KEY idx_payment_no (payment_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS payment_idempotency (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  idem_key VARCHAR(128) NOT NULL,
  biz_type VARCHAR(64) NOT NULL,
  biz_no VARCHAR(64) NOT NULL,
  response_json MEDIUMTEXT,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uniq_idem_biz (idem_key, biz_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS refund_task (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  refund_task_no VARCHAR(64) NOT NULL,
  after_sale_no VARCHAR(64) NOT NULL,
  order_no VARCHAR(64) NOT NULL,
  payment_no VARCHAR(64) NOT NULL,
  refund_amount BIGINT UNSIGNED NOT NULL DEFAULT 0,
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  retry_count INT UNSIGNED NOT NULL DEFAULT 0,
  next_retry_at DATETIME NULL,
  last_error_code VARCHAR(64) NOT NULL DEFAULT '',
  last_error_message VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uniq_refund_task_no (refund_task_no),
  UNIQUE KEY uniq_after_sale_no (after_sale_no),
  KEY idx_status_next_retry (status, next_retry_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS payment_reconciliation_task (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  recon_task_no VARCHAR(64) NOT NULL,
  recon_date DATE NOT NULL,
  pay_channel TINYINT UNSIGNED NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
  total_records BIGINT UNSIGNED NOT NULL DEFAULT 0,
  diff_records BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uniq_recon_task_no (recon_task_no),
  UNIQUE KEY uniq_recon_date_channel (recon_date, pay_channel)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS payment_reconciliation_record (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  diff_no VARCHAR(64) NOT NULL,
  recon_task_no VARCHAR(64) NOT NULL,
  diff_type TINYINT UNSIGNED NOT NULL DEFAULT 0,
  payment_no VARCHAR(64) NOT NULL DEFAULT '',
  order_no VARCHAR(64) NOT NULL DEFAULT '',
  local_amount BIGINT UNSIGNED NOT NULL DEFAULT 0,
  gateway_amount BIGINT UNSIGNED NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'OPEN',
  detail_json MEDIUMTEXT,
  resolved_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uniq_diff_no (diff_no),
  KEY idx_recon_task_no (recon_task_no),
  KEY idx_status_created (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS payment_outbox_event (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  event_no VARCHAR(64) NOT NULL,
  event_type VARCHAR(64) NOT NULL,
  biz_no VARCHAR(64) NOT NULL,
  payload_json MEDIUMTEXT,
  status TINYINT UNSIGNED NOT NULL DEFAULT 0,
  retry_count INT UNSIGNED NOT NULL DEFAULT 0,
  next_retry_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uniq_event_no (event_no),
  KEY idx_status_retry (status, next_retry_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;