ALTER TABLE `agent_conversation`
  ADD COLUMN `last_queue_notice_at` DATETIME NULL AFTER `last_message_at`;

ALTER TABLE `agent_run`
  ADD COLUMN `subject_user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `shop_no`,
  ADD COLUMN `subject_shop_no` VARCHAR(64) NOT NULL DEFAULT '' AFTER `subject_user_id`,
  ADD COLUMN `current_node_code` VARCHAR(64) NOT NULL DEFAULT '' AFTER `run_status_code`,
  ADD COLUMN `graph_state_json` JSON NULL AFTER `current_node_code`,
  ADD COLUMN `checkpoint_version` BIGINT UNSIGNED NOT NULL DEFAULT 1 AFTER `graph_state_json`,
  ADD COLUMN `tool_wait_timeout_at` DATETIME NULL AFTER `checkpoint_version`,
  ADD COLUMN `tool_result_status` VARCHAR(32) NOT NULL DEFAULT 'UNSPECIFIED' AFTER `tool_wait_timeout_at`,
  ADD COLUMN `degraded_reason_code` VARCHAR(64) NOT NULL DEFAULT '' AFTER `tool_result_status`,
  ADD COLUMN `queue_blocked` TINYINT(1) NOT NULL DEFAULT 0 AFTER `degraded_reason_code`,
  ADD COLUMN `queue_hint_message` VARCHAR(512) NOT NULL DEFAULT '' AFTER `queue_blocked`;

ALTER TABLE `agent_run`
  ADD KEY `idx_run_status_timeout` (`run_status_code`, `tool_wait_timeout_at`, `id`);

ALTER TABLE `agent_escalation_ticket`
  ADD COLUMN `queue_appendix_json` JSON NULL AFTER `remark`,
  ADD COLUMN `accepted_at` DATETIME NULL AFTER `handoff_generated_at`,
  ADD COLUMN `closed_at` DATETIME NULL AFTER `accepted_at`;

ALTER TABLE `agent_tool_call_log`
  ADD COLUMN `adapter_code` VARCHAR(64) NOT NULL DEFAULT '' AFTER `tool_name`,
  ADD COLUMN `tool_scope_code` VARCHAR(32) NOT NULL DEFAULT 'UNSPECIFIED' AFTER `adapter_code`,
  ADD COLUMN `subject_user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER `tool_scope_code`,
  ADD COLUMN `subject_shop_no` VARCHAR(64) NOT NULL DEFAULT '' AFTER `subject_user_id`,
  ADD COLUMN `request_payload_json` JSON NULL AFTER `subject_shop_no`,
  ADD COLUMN `response_payload_json` JSON NULL AFTER `request_payload_json`,
  ADD COLUMN `tool_result_status` VARCHAR(32) NOT NULL DEFAULT 'UNSPECIFIED' AFTER `response_payload_json`,
  ADD COLUMN `degraded_reason_code` VARCHAR(64) NOT NULL DEFAULT '' AFTER `tool_result_status`,
  ADD COLUMN `error_code` VARCHAR(64) NOT NULL DEFAULT '' AFTER `degraded_reason_code`,
  ADD COLUMN `error_message` VARCHAR(512) NOT NULL DEFAULT '' AFTER `error_code`,
  ADD COLUMN `duration_ms` INT UNSIGNED NOT NULL DEFAULT 0 AFTER `error_message`,
  ADD COLUMN `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER `created_at`;

ALTER TABLE `agent_tool_call_log`
  ADD KEY `idx_subject_scope` (`subject_user_id`, `subject_shop_no`, `tool_scope_code`, `id`);

UPDATE `agent_run`
SET `accepted_status_code` = 'PENDING'
WHERE `accepted_status_code` = '' OR `accepted_status_code` = 'RECEIVED';

UPDATE `agent_run`
SET `run_status_code` = 'PENDING'
WHERE `run_status_code` = '' OR `run_status_code` = 'RECEIVED';

UPDATE `agent_run`
SET `tool_result_status` = 'UNSPECIFIED'
WHERE `tool_result_status` = '';

UPDATE `agent_conversation`
SET `last_run_status_code` = 'PENDING'
WHERE `last_run_status_code` = '' OR `last_run_status_code` = 'RECEIVED';

UPDATE `agent_escalation_ticket`
SET `status_code` = 'ESCALATION_PENDING'
WHERE `status_code` = '' OR `status_code` = 'OPEN';
