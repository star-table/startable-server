ALTER TABLE `ppm_pri_issue`
MODIFY COLUMN `sort` int NOT NULL DEFAULT 0 AFTER `end_time`;

ALTER TABLE `ppm_pro_project`
DROP COLUMN `template_flag`;