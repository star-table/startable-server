CREATE TABLE `ppm_rol_role_department` (
  `id` bigint NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `project_id` bigint NOT NULL DEFAULT '0',
  `role_id` bigint NOT NULL DEFAULT '0',
  `department_id` bigint NOT NULL DEFAULT '0',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_rol_role_department_org_id` (`org_id`),
  KEY `index_ppm_rol_role_department_role_id` (`role_id`),
  KEY `index_ppm_rol_role_department_department_id_project_id` (`department_id`,`project_id`),
  KEY `index_ppm_rol_role_department_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
