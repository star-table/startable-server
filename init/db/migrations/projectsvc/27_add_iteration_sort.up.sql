ALTER TABLE `ppm_pri_iteration`
ADD COLUMN `sort`  bigint NOT NULL DEFAULT 0 COMMENT '排序' AFTER `owner`;
