ALTER TABLE `ppm_prs_process_status`
ADD COLUMN `project_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '项目id' AFTER `org_id`,
ADD INDEX `index_ppm_prs_process_status_project_id`(`project_id`) USING BTREE;

ALTER TABLE `ppm_prs_process_process_status`
ADD COLUMN `project_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '项目id' AFTER `org_id`,
ADD INDEX `index_ppm_prs_process_status_project_id`(`project_id`) USING BTREE;

INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 28, 0, 'ProcessStatus.Issue.Pending', '待处理', 1, '#DBDBDB', '#FFFFFF', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);

Update ppm_prs_process_status set name = '已确认' where `id` = 10;

INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 29, 0, 'ProcessStatus.Issue.RetestThrough', '复测通过', 16, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 30, 0, 'ProcessStatus.Issue.Resolved', '已解决', 22, '#FFCD1C', '#FFFFFF', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
