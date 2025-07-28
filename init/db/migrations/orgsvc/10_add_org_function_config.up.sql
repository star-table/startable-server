CREATE TABLE `ppm_orc_function_config` (
  `id` bigint NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `function_code` varchar(50) NOT NULL DEFAULT '',
  `remark` varchar(250) NOT NULL DEFAULT '',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_orc_function_config_org_id` (`org_id`),
  KEY `index_ppm_orc_function_config_function_code` (`function_code`),
  KEY `index_ppm_orc_function_config_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
