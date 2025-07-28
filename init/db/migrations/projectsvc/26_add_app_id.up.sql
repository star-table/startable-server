-- 为项目表增加关联的 app_id 字段。
ALTER TABLE `ppm_pro_project`
    ADD COLUMN `app_id` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '融合无码系统，关联的应用id' AFTER `org_id`;


-- 因为融合版增加了些sql导致版本号增加，这会导致和主线版本极星冲突，因此将两个文件合并到一个版本中，以解决冲突。
-- for 27_add_issue_audit_status.up.sql start
ALTER TABLE `ppm_pri_issue`
    ADD COLUMN `audit_status` tinyint NOT NULL DEFAULT '1' AFTER `is_filing`;

ALTER TABLE `ppm_pri_issue_relation`
    ADD COLUMN `status` tinyint NOT NULL DEFAULT '1' AFTER `relation_type`;
-- for 27_add_issue_audit_status.up.sql end
