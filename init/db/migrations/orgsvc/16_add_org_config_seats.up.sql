ALTER TABLE `ppm_orc_config`
ADD COLUMN `seats` int NOT NULL DEFAULT '0' COMMENT '实际购买人数' AFTER `pay_level`;
