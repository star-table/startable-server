ALTER TABLE `ppm_prs_process_status` DROP COLUMN `project_id`;

ALTER TABLE `ppm_prs_process_process_status` DROP COLUMN `project_id`;

delete from ppm_prs_process_status where id between 28 and 30;

Update ppm_prs_process_status set name = '已确认缺陷' where `id` = 10;
