-- 项目（应用）群聊关联表
-- 一个项目可以关联多个群聊，一个群聊可以绑定多个项目
CREATE TABLE `ppm_pro_project_chat` (
    `id` bigint NOT NULL,
    `org_id` bigint NOT NULL DEFAULT '0',
    `project_id` bigint NOT NULL DEFAULT '0' COMMENT '项目的 id',
    `table_id` bigint NOT NULL DEFAULT '0' COMMENT '项目下表的 id',
    `chat_id` varchar(50) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '第三方群聊 id',
    `chat_type` tinyint NOT NULL DEFAULT '0' COMMENT '1创建项目时创建的群聊（自带群聊）；2机器人被拉入已存在的群聊',
    `chat_settings` json DEFAULT NULL COMMENT '群聊配置 json。其中包含多个表对应的配置',
    `is_enable` tinyint NOT NULL DEFAULT '1' COMMENT '1启用群聊；2不启用',
    `creator` bigint NOT NULL DEFAULT '0',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updator` bigint NOT NULL DEFAULT '0',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `version` int NOT NULL DEFAULT '1',
    `is_delete` tinyint NOT NULL DEFAULT '2',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `chat_id_IDX` (`chat_id`) USING BTREE,
    KEY `project_chat_org_id_IDX` (`org_id`,`project_id`,`table_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='项目（应用）、群聊配置关联表';
