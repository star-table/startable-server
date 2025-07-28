ALTER TABLE `ppm_pri_issue` ADD COLUMN `property_id` bigint(20) NOT NULL DEFAULT '0' AFTER `source_id`;

ALTER TABLE `ppm_prs_issue_object_type` ADD COLUMN `project_id` bigint(20) NOT NULL DEFAULT '0' AFTER `org_id`;
ALTER TABLE `ppm_prs_issue_source` ADD COLUMN `project_id` bigint(20) NOT NULL DEFAULT '0' AFTER `org_id`;

ALTER TABLE `ppm_pri_issue` change `issus_object_type_id` `issue_object_type_id` bigint(20) NOT NULL DEFAULT '0';

-- 任务性质（严重程度）
CREATE TABLE `ppm_prs_issue_property` (
  `id` bigint(20) NOT NULL,
  `org_id` bigint(20) NOT NULL DEFAULT '0',
  `project_id` bigint(20) NOT NULL DEFAULT '0',
  `lang_code` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `name` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `sort` int(11) NOT NULL DEFAULT '0',
  `project_object_type_id` bigint(20) NOT NULL DEFAULT '0',
  `remark` varchar(512) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `creator` bigint(20) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint(20) NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int(11) NOT NULL DEFAULT '1',
  `is_delete` tinyint(4) NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_prs_issue_property_org_id` (`org_id`),
  KEY `index_ppm_prs_issue_property_project_id` (`project_id`),
  KEY `index_ppm_prs_issue_property_project_object_type_id` (`project_object_type_id`),
  KEY `index_ppm_prs_issue_property_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
