ALTER TABLE `ppm_org_user_login_record` ADD COLUMN `source_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `login_ip`;
ALTER TABLE `ppm_org_user_login_record` ADD COLUMN `source_platform` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `login_ip`;
