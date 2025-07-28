ALTER TABLE `ppm_prs_issue_object_type` DROP COLUMN `project_id`;
ALTER TABLE `ppm_prs_issue_source` DROP COLUMN `project_id`;


drop table if exists ppm_prs_issue_property;

ALTER TABLE `ppm_pri_issue` change `issue_object_type_id` `issus_object_type_id` bigint(20) NOT NULL DEFAULT '0';

