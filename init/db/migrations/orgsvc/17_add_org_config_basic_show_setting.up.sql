ALTER TABLE `ppm_orc_config`
ADD COLUMN `basic_show_setting` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `pay_level`;

ALTER TABLE `ppm_org_organization_out_info`
ADD COLUMN `tenant_code` varchar(64) NOT NULL COMMENT '企业编号' DEFAULT '' AFTER `source_platform`;