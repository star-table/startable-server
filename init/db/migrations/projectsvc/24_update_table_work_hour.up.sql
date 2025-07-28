-- up
ALTER TABLE `ppm_pri_issue_work_hours`
ADD COLUMN `project_id` bigint(0) NOT NULL COMMENT '项目id' AFTER `org_id`,
DROP INDEX `org_issue_id_idx`,
ADD INDEX `org_proj_issue_type_id_idx`(`org_id`, `project_id`, `issue_id`, `type`) USING BTREE;


-- 将历史数据改为不开启，需要用户手动开启工时功能。
UPDATE `ppm_pro_project_detail` SET `is_enable_work_hours`=2 WHERE `create_time`<"2020-11-23 23:58:00";

-- 项目工时功能默认不开启。
ALTER TABLE `ppm_pro_project_detail`
MODIFY COLUMN `is_enable_work_hours` tinyint(0) NOT NULL DEFAULT 2 AFTER `notice`;
