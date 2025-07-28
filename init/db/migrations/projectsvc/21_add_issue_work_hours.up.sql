-- 工时记录表 ppm_pri_issue_work_hours
CREATE TABLE IF NOT EXISTS `ppm_pri_issue_work_hours` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `org_id` bigint NOT NULL DEFAULT 0 COMMENT '组织id',
  `issue_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '关联的任务id',
  `type` tinyint(1) NOT NULL DEFAULT 2 COMMENT '记录类型：1预估工时记录，2实际工时记录，3详细预估工时',
  `worker_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '工作者id',
  `need_time` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '所需工时时间，单位分钟',
  `remain_time_cal_type` tinyint(1) unsigned NOT NULL DEFAULT 1 COMMENT '剩余工时计算方式：1动态计算；2手动填写',
  `remain_time`  int(11) unsigned NOT NULL DEFAULT 0 COMMENT '手动填写工时时的剩余工时，单位分钟',
  `start_time` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '开始时间，时间戳',
  `end_time` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '工时记录的结束时间，时间戳',
  `desc` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '工时记录的内容，工作内容',
  `creator` bigint(20) NOT NULL DEFAULT 0 COMMENT '工时记录创建者id',
  `updator` bigint(20) NOT NULL DEFAULT 0 COMMENT '工时记录更新者的id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `org_issue_id_idx` (`org_id`, `issue_id`) USING BTREE,
  KEY `worker_id_idx` (`worker_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='工时记录表';

-- 新增工时记录
-- 编辑工时记录
-- 删除工时记录
-- 工时记录列表
-- 开启/关闭工时功能