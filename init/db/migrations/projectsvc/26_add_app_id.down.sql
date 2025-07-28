ALTER TABLE `ppm_pro_project` DROP COLUMN `app_id`;


-- for 27_add_issue_audit_status.down.sql start
ALTER TABLE `ppm_pri_issue`
DROP COLUMN `audit_status`;

ALTER TABLE `ppm_pri_issue_relation`
DROP COLUMN `status`;
-- for 27_add_issue_audit_status.down.sql end
