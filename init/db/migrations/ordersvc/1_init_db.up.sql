CREATE TABLE `ppm_ord_order` (
  `id` bigint NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `out_order_no`bigint NOT NULL DEFAULT '0' COMMENT '外部订单id',
  `status` bigint NOT NULL DEFAULT '0' COMMENT '状态（1未支付2已支付3已取消4已过期5部分支付6已退款）',
  `order_create_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '订单创建时间',
  `paid_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '支付时间',
  `effective_time` int NOT NULL DEFAULT '0' COMMENT '购买时长（月）',
  `buy_count` int NOT NULL DEFAULT '0' COMMENT '购买数量',
  `seats` int NOT NULL DEFAULT '0' COMMENT '购买人数',
  `total_price` bigint NOT NULL DEFAULT '0' COMMENT '总金额（分）',
  `order_pay_price` bigint NOT NULL DEFAULT '0' COMMENT '支付金额（分）',
  `source_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `buy_type` varchar(128) NOT NULL DEFAULT '' COMMENT '购买类型，"buy" - 普通购买;"upgrade"-为升级购买;"renew" - 续费购买',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_ord_order_org_id` (`org_id`),
  KEY `index_ppm_ord_order_paid_time` (`paid_time`),
  KEY `index_ppm_ord_order_order_pay_price` (`order_pay_price`),
  KEY `index_ppm_ord_order_out_order_no` (`out_order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;



CREATE TABLE `ppm_ord_order_fs` (
  `id` bigint NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `order_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '订单编号id',
  `price_plan_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '价格方案id',
  `price_plan_type` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '价格方案类型',
  `seats` int NOT NULL DEFAULT '0' COMMENT '实际购买人数',
  `buy_count` int NOT NULL DEFAULT '0' COMMENT '购买数量',
  `paid_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '支付时间',
  `order_create_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '订单创建时间',
  `status` varchar(20) NOT NULL DEFAULT '' COMMENT '订单当前状态，"normal" -正常；"refund"-已退款',
  `buy_type` varchar(128) NOT NULL DEFAULT '' COMMENT '购买类型，"buy" - 普通购买;"upgrade"-为升级购买;"renew" - 续费购买',
  `src_order_id` varchar(128) NOT NULL DEFAULT '' COMMENT '源订单ID',
  `dst_order_id` varchar(128) NOT NULL DEFAULT '' COMMENT '升级后的新订单ID',
  `order_pay_price` bigint NOT NULL DEFAULT '0' COMMENT '订单实际支付金额, 单位分',
  `tenant_key` varchar(128) NOT NULL DEFAULT '' COMMENT '租户唯一标识',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_ord_order_fs_org_id` (`org_id`),
  KEY `index_ppm_ord_order_fs_order_id` (`order_id`),
  KEY `index_ppm_ord_order_fs_paid_time` (`paid_time`),
  KEY `index_ppm_ord_order_fs_order_pay_price` (`order_pay_price`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `ppm_ord_price_plan_fs` (
  `id` bigint NOT NULL,
  `out_plan_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '外部付费方案id',
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '方案名称',
  `price_plan_type` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '价格方案类型',
  `seats` int NOT NULL DEFAULT '0' COMMENT '人数限制',
  `trial_days` int NOT NULL DEFAULT '0' COMMENT '试用天数',
  `month_price` bigint NOT NULL DEFAULT '0' COMMENT '月费金额, 单位分',
  `year_price` bigint NOT NULL DEFAULT '0' COMMENT '年费金额, 单位分',
  `level` int NOT NULL DEFAULT '0' COMMENT '等级',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_ord_price_plan_fs_out_plan_id` (`out_plan_id`),
  KEY `index_ppm_ord_price_plan_fs_month_price` (`month_price`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `ppm_ord_function` (
  `id` bigint NOT NULL,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '名称',
  `code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '功能项code',
  `type` int NOT NULL DEFAULT '0' COMMENT '类型',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_ord_function_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `ppm_ord_function_level` (
  `id` bigint NOT NULL,
  `function_id` bigint NOT NULL DEFAULT '0' COMMENT '功能id',
  `level` int NOT NULL DEFAULT '0' COMMENT '等级',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_ord_function_level_function_id` (`function_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
