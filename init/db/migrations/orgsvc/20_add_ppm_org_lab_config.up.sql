CREATE TABLE `ppm_org_lab_config`
(
    `id`          bigint   NOT NULL AUTO_INCREMENT,
    `org_id`      bigint   NOT NULL DEFAULT '0',
    `config`      json              DEFAULT NULL,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `version`     int      NOT NULL DEFAULT '1',
    `is_delete`   tinyint  NOT NULL DEFAULT '2',
    PRIMARY KEY (`id`) USING BTREE,
    KEY           `ppm_org_lab_config_org_id_IDX` (`org_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
