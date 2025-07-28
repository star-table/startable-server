-- 无码数据库表的迁移，备份
ALTER TABLE lc_app
    ADD `mirror_table_id` bigint NOT NULL DEFAULT '0' COMMENT '镜像表格id';


CREATE TABLE `lc_table_collaborator_permissions`
(
    `id`         bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `org_id`     bigint                                                        DEFAULT '0' COMMENT '组织id',
    `app_id`     bigint                                                        DEFAULT '0' COMMENT '应用id',
    `table_id`   bigint                                                        DEFAULT '0' COMMENT '表id',
    `field_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '表列名id',
    `data_type`  tinyint                                                       DEFAULT NULL COMMENT '1:用户 2.部门',
    `data_id`    bigint                                                        DEFAULT NULL COMMENT 'data_type 1: 用户id 2:部门id',
    `row_count`  int                                                           DEFAULT '1' COMMENT '列协作累加数量',
    `tag`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT 'tag',
    `creator`    bigint                                                        DEFAULT '0' COMMENT '创建人',
    `updator`    bigint                                                        DEFAULT '0' COMMENT '更新人',
    `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) COMMENT '建立时间',
    `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    unique KEY `idx_tag` (`tag`) USING BTREE,
    KEY          `idx_org_app_table` (`org_id`,`app_id`,`data_id`,`data_type`,`table_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci


ALTER TABLE lc_per_app_permission_group
    ADD `table_auth` json DEFAULT NULL COMMENT '表格权限';
ALTER TABLE lc_per_app_permission_group
    ADD `column_auth` json DEFAULT NULL COMMENT '字段权限(新)';
