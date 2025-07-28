ALTER TABLE `ppm_org_department` DROP COLUMN `path`;

drop table if exists ppm_org_user_invite;

ALTER TABLE `ppm_orc_config` DROP COLUMN `db_id`;
ALTER TABLE `ppm_orc_config` DROP COLUMN `ds_id`;
ALTER TABLE `ppm_orc_config` DROP COLUMN `dc_id`;

ALTER TABLE `ppm_org_user_organization` DROP COLUMN `invite_id`;
