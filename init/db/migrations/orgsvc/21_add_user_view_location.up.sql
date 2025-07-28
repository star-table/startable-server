CREATE TABLE `ppm_org_user_view_location`
(
    `id`          bigint   NOT NULL,
    `org_id`      bigint   NOT NULL                 DEFAULT '0',
    `user_id`     bigint   NOT NULL                 DEFAULT '0',
    `app_id`      bigint   NOT NULL                 DEFAULT '0',
    `config`      varchar(1024) COLLATE utf8mb4_bin DEFAULT '',
    `creator`     bigint                            DEFAULT '0',
    `create_time` datetime NOT NULL                 DEFAULT CURRENT_TIMESTAMP,
    `updator`     bigint   NOT NULL                 DEFAULT '0',
    `update_time` datetime NOT NULL                 DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `version`     int      NOT NULL                 DEFAULT '1',
    `is_delete`   tinyint  NOT NULL                 DEFAULT '2',
    PRIMARY KEY (`id`),
    KEY           `ppm_org_user_view_location_org_id_IDX` (`org_id`,`user_id`,`app_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
