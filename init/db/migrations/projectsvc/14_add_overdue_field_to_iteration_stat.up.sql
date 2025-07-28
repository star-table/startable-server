ALTER TABLE `ppm_sta_iteration_stat` ADD COLUMN `issue_overdue_count` int NOT NULL DEFAULT '0' AFTER `issue_running_count`;
ALTER TABLE `ppm_sta_iteration_stat` ADD COLUMN `demand_overdue_count` int NOT NULL DEFAULT '0' AFTER `demand_running_count`;
ALTER TABLE `ppm_sta_iteration_stat` ADD COLUMN `story_point_overdue_count` int NOT NULL DEFAULT '0' AFTER `story_point_running_count`;
ALTER TABLE `ppm_sta_iteration_stat` ADD COLUMN `task_overdue_count` int NOT NULL DEFAULT '0' AFTER `task_running_count`;
ALTER TABLE `ppm_sta_iteration_stat` ADD COLUMN `bug_overdue_count` int NOT NULL DEFAULT '0' AFTER `bug_running_count`;
ALTER TABLE `ppm_sta_iteration_stat` ADD COLUMN `testtask_overdue_count` int NOT NULL DEFAULT '0' AFTER `testtask_running_count`;
