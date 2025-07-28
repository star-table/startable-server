ALTER TABLE `ppm_org_department`
ADD COLUMN `path`  varchar(255) NOT NULL DEFAULT "0" COMMENT '父级路径拼接（“,”隔开）' AFTER `parent_id`;

CREATE TABLE `ppm_org_user_invite` (
  `id` bigint(20) NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '姓名',
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '邮箱',
  `invite_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '邀请人id',
  `is_register` int NOT NULL DEFAULT '2' COMMENT '是否注册（1已注册2未注册）',
  `last_invite_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '上次邀请时间',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_org_user_invite_org_id` (`org_id`),
  KEY `index_ppm_org_user_invite_create_time` (`create_time`),
  KEY `index_ppm_org_user_invite_invite_user_id` (`invite_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `ppm_orc_config` ADD COLUMN `db_id` bigint NOT NULL DEFAULT "0" COMMENT '数据库id' AFTER `time_difference`;
ALTER TABLE `ppm_orc_config` ADD COLUMN `ds_id` bigint NOT NULL DEFAULT "0" COMMENT '数据源id' AFTER `time_difference`;
ALTER TABLE `ppm_orc_config` ADD COLUMN `dc_id` bigint NOT NULL DEFAULT "0" COMMENT '数据中心id' AFTER `time_difference`;

ALTER TABLE `ppm_org_user_organization` ADD COLUMN `invite_id` bigint NOT NULL DEFAULT "0" COMMENT '邀请人id' AFTER `user_id`;
