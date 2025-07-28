ALTER TABLE `ppm_pri_issue`
MODIFY COLUMN `sort` bigint NOT NULL DEFAULT 0 AFTER `end_time`;

alter table ppm_pro_project add column `template_flag` tinyint NOT NULL DEFAULT '2' COMMENT '模板标识：1：是，2：否' after public_status;