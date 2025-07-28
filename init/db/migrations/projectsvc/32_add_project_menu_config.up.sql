CREATE TABLE `ppm_pro_project_menu_config`
(
    `id`          bigint   NOT NULL AUTO_INCREMENT,
    `org_id`      bigint   NOT NULL DEFAULT '0',
    `app_id`      bigint unsigned NOT NULL DEFAULT '0',
    `config`      longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `version`     int      NOT NULL DEFAULT '1',
    `is_delete`   tinyint  NOT NULL DEFAULT '2',
    PRIMARY KEY (`id`) USING BTREE,
    KEY           `ppm_pro_project_menu_config_org_id_IDX` (`org_id`,`app_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
