-- 调整字段
alter table ppm_pro_project modify column pre_code varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '';
alter table ppm_pri_issue modify column `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '';