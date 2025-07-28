ALTER TABLE `ppm_pri_issue` ADD COLUMN `path` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '0,' AFTER `title`;

update ppm_pri_issue set path = CONCAT("0,",parent_id, ",") where parent_id != 0
