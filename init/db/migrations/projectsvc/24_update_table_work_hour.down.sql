-- down
ALTER TABLE `ppm_pri_issue_work_hours`
DROP COLUMN `project_id`,
DROP INDEX `org_proj_issue_type_id_idx`,
ADD INDEX `org_issue_id_idx`(`org_id`, `issue_id`) USING BTREE;

-- 原始状态是：开启。
UPDATE `ppm_pro_project_detail` SET `is_enable_work_hours`=1 WHERE `create_time`<"2020-11-11 23:58:00";

-- 原始状态是：项目工时功能默认开启。
ALTER TABLE `ppm_pro_project_detail`
MODIFY COLUMN `is_enable_work_hours` tinyint(0) NOT NULL DEFAULT '1' AFTER `notice`;
