CREATE TABLE `ppm_pro_custom_field` (
  `id` bigint NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `field_type` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `field_value` varchar(10000) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `remark` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `is_org_field` tinyint NOT NULL DEFAULT '2',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_pro_custom_field_org_id` (`org_id`),
  KEY `index_ppm_pro_custom_field_name` (`name`),
  KEY `index_ppm_pro_custom_field_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `ppm_pro_project_view` (
  `id` bigint NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `project_id` bigint NOT NULL DEFAULT '0',
  `project_object_type_id` bigint NOT NULL DEFAULT '0',
  `view_type` tinyint NOT NULL DEFAULT '0',
  `closed_default_field` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `closed_custom_field` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_pro_project_view_org_id` (`org_id`),
  KEY `index_ppm_pro_project_view_project_id` (`project_id`),
  KEY `index_ppm_pro_project_view_project_object_type_id` (`project_object_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

ALTER TABLE `ppm_pri_issue`
ADD COLUMN `custom_field`  json AFTER `parent_id`;

INSERT INTO `ppm_pro_custom_field` (`id`, `org_id`, `name`, `field_type`, `field_value`, `remark`, `is_org_field`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`) VALUES ('900', '0', '进度', '7', '[{\"fieldName\":\"字段格式\",\"value\":\"percentage\",\"type\":2,\"id\":\"5\"},{\"fieldName\":\"小数点位数\",\"value\":\"0\",\"type\":4,\"id\":\"1\"}]', '', '3', '0', '2020-11-09 10:36:40', '0', '2020-11-09 10:36:40', '1', '2');
INSERT INTO `ppm_pro_custom_field` (`id`, `org_id`, `name`, `field_type`, `field_value`, `remark`, `is_org_field`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`) VALUES ('901', '0', 'Story Points', '2', '[{\"fieldName\":\"选项值\",\"value\":\"无\",\"type\":1,\"id\":\"6e696cc1-a4e5-198b-a647-4278b932f1e6\"},{\"value\":\"0\",\"type\":1,\"id\":\"682ce558-4dfb-9725-5a41-ddb8df4f81f3\",\"fieldName\":\"选项值\"},{\"id\":\"65b07ef9-b72f-3864-0293-a48df46e5bf9\",\"fieldName\":\"选项值\",\"value\":\"0.5\",\"type\":1},{\"id\":\"10bd5f98-bbb8-a1ac-4f79-d8af4f1d77ee\",\"fieldName\":\"选项值\",\"value\":\"1\",\"type\":1},{\"fieldName\":\"选项值\",\"value\":\"2\",\"type\":1,\"id\":\"751f46a6-9e0d-d70a-d398-4ba020ba4841\"},{\"value\":\"3\",\"type\":1,\"id\":\"2d74b8e7-3c46-4aed-9a0b-683032594cbb\",\"fieldName\":\"选项值\"},{\"id\":\"fac03b7b-4845-faf4-d007-828a003af48d\",\"fieldName\":\"选项值\",\"value\":\"5\",\"type\":1},{\"value\":\"8\",\"type\":1,\"id\":\"d19c3fb7-52e4-c470-4034-c8186ace0561\",\"fieldName\":\"选项值\"},{\"value\":\"10\",\"type\":1,\"id\":\"b1792a26-fa32-1d5a-67ba-6211b25db50b\",\"fieldName\":\"选项值\"},{\"value\":\"20\",\"type\":1,\"id\":\"40d47aad-0357-0b9f-dc03-94d27154ca47\",\"fieldName\":\"选项值\"},{\"id\":\"5f5f8bbc-617a-9e2c-c5e6-b9001bd6bd15\",\"fieldName\":\"选项值\",\"value\":\"40\",\"type\":1},{\"value\":\"100\",\"type\":1,\"id\":\"58c6c8a5-ff22-014b-9cf8-4b6198af0b9d\",\"fieldName\":\"选项值\"}]', '', '3', '0', '2020-11-09 11:04:12', '0', '2020-11-09 11:04:12', '1', '2');
INSERT INTO `ppm_pro_custom_field` (`id`, `org_id`, `name`, `field_type`, `field_value`, `remark`, `is_org_field`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`) VALUES ('902', '0', '评分', '2', '[{\"id\":\"800d321d-10d5-ade3-8de1-9e4335a3121f\",\"fieldName\":\"选项值\",\"value\":\"1\",\"type\":1},{\"fieldName\":\"选项值\",\"value\":\"2\",\"type\":1,\"id\":\"b97d6d69-7c0f-11eb-733e-40d6e55e66e8\"},{\"value\":\"3\",\"type\":1,\"id\":\"2ec0b186-a501-0159-02f6-1501c3107f8c\",\"fieldName\":\"选项值\"},{\"value\":\"4\",\"type\":1,\"id\":\"0d823c33-934d-c5c2-e65b-d98de28c4c7a\",\"fieldName\":\"选项值\"},{\"type\":1,\"id\":\"6d055ee9-cfd7-9682-3e0b-1af0fbc84c78\",\"fieldName\":\"选项值\",\"value\":\"5\"},{\"value\":\"6\",\"type\":1,\"id\":\"cdc9df11-d929-3f11-d7f7-9b61c9695f1f\",\"fieldName\":\"选项值\"},{\"id\":\"6b24cf4b-0d2b-bb7c-c201-bcb42df2ce40\",\"fieldName\":\"选项值\",\"value\":\"7\",\"type\":1},{\"type\":1,\"id\":\"b2ca6d6d-6d16-375b-0934-04fa53941df1\",\"fieldName\":\"选项值\",\"value\":\"8\"},{\"value\":\"9\",\"type\":1,\"id\":\"534bd785-4a53-201b-dfe8-1d22fb4abc00\",\"fieldName\":\"选项值\"},{\"fieldName\":\"选项值\",\"value\":\"10\",\"type\":1,\"id\":\"14cc9325-ee3d-1c87-0db0-5b0388dff02d\"}]', '', '3', '0', '2020-11-09 11:05:25', '0', '2020-11-09 11:05:25', '1', '2');

ALTER TABLE `ppm_pro_project_relation`
ADD COLUMN `project_object_type_id`  bigint NOT NULL DEFAULT 0 AFTER `project_id`;
