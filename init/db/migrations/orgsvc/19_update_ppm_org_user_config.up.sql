ALTER TABLE ppm_org_user_config ADD user_view_location_config json NULL COMMENT '用户上一次浏览的位置';

ALTER TABLE `ppm_org_user_config` ADD COLUMN `collaborate_message_status` tinyint not null default 0 COMMENT '用户推送配置-我协作的' AFTER `pc_project_message_status`;
