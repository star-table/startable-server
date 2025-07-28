CREATE TABLE `ppm_org_user_login_record` (
  `id` bigint NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `user_id` bigint NOT NULL DEFAULT '0',
  `login_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `login_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_org_user_org_id` (`org_id`),
  KEY `index_ppm_org_user_user_id` (`user_id`),
  KEY `index_ppm_org_user_login_time` (`login_time`),
  KEY `index_ppm_org_user_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
