ALTER TABLE `ppm_ord_price_plan_fs`
ADD COLUMN `expire_days`  int NOT NULL DEFAULT 0 COMMENT '有效天数' AFTER `trial_days`,
ADD COLUMN `end_date`  datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '到期日期' AFTER `expire_days`;
