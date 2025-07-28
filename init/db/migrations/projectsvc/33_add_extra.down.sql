ALTER TABLE `ppm_pri_issue` DROP COLUMN `table_id`;

ALTER TABLE `ppm_pri_issue` DROP INDEX `index_ppm_pri_issue_project_id_table_id_sort`;

DELETE from `ppm_prs_project_type` where id=47;
