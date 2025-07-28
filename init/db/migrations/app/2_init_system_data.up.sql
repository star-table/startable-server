INSERT INTO ppm_org_organization( `id`, `name`, `web_site`, `industry_id`, `scale`, `source_channel`, `country_id`, `province_id`, `city_id`, `address`, `logo_url`, `resource_id`, `owner`, `is_authenticated`, `remark`, `init_status`, `status`, `is_show`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, '极星组织', '', 0, '', 'web', 0, 0, 0, '', '', 0, 1, 1, '', 3, 1, 0, 1, now(), 1, now(), 1, 1);
INSERT INTO ppm_org_user( `id`, `org_id`, `name`, `login_name`, `login_name_edit_count`, `email`, `mobile`, `birthday`, `sex`, `password`, `password_salt`, `source_channel`, `language`, `motto`, `last_login_ip`, `last_login_time`, `login_fail_count`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, 1, '管理员', 'admin', 99, 'admin@admin.com', '12345678901', '1970-01-01 00:00:00', 99, '2eab55be9dfa1ba036d9024baf2c60c8', '8a66572aac6911e99080784f439212b0', 'web', 'zh-CN', '', '', '1970-01-01 00:00:00', 0, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_org_user_organization( `id`, `org_id`, `user_id`, `check_status`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES ( 1, 1, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_tem_team( `id`, `org_id`, `name`, `nick_name`, `owner`, `department_id`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES ( 1, 1, '极星组织团队', '极星组织团队', 1, 0, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_tem_user_team( `id`, `org_id`, `team_id`, `user_id`, `relation_type`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES ( 1, 1, 1, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_source_channel ( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, 'web', 'web主站', '', 1, 1, now( ), 1, now( ), 1, 2 );
INSERT INTO ppm_bas_source_channel ( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 2, 'ding', '钉钉', '', 1, 1, now( ), 1, now( ), 1, 2 );
INSERT INTO ppm_bas_source_channel ( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 3, 'weixin', '微信', '', 1, 1, now( ), 1, now( ), 1, 2 );
INSERT INTO `ppm_bas_pay_level`( `id`, `lang_code`, `name`, `storage`, `member_count`, `price`, `member_price`, `duration`, `is_show`, `sort`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES (1, 'PayLevel.Trial', '试用级别', 1024000, 100, 0, 0, 1296000, 2, 0, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_bas_pay_level`( `id`, `lang_code`, `name`, `storage`, `member_count`, `price`, `member_price`, `duration`, `is_show`, `sort`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES (2, 'PayLevel.Free', '免费版', 1024000, 10, 0, 0, 0, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_bas_pay_level`( `id`, `lang_code`, `name`, `storage`, `member_count`, `price`, `member_price`, `duration`, `is_show`, `sort`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES (3, 'PayLevel.Professional', '专业版', 1024000000, 100, 25600, 25600, 31536000, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_bas_pay_level`( `id`, `lang_code`, `name`, `storage`, `member_count`, `price`, `member_price`, `duration`, `is_show`, `sort`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES (4, 'PayLevel.Vip', 'VIP版', 1024000000, 1000, 102400, 25600, 31536000, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0,  'ppm_bas_app_info', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0,  'ppm_bas_dictionary', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0,  'ppm_bas_object_id', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0,  'ppm_bas_pay_level', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0,  'ppm_bas_source_channel', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0,  'ppm_org_organization', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0,  'ppm_org_organization_out_info', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0,  'ppm_org_user', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0,  'ppm_org_user_organization', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 10, 0,  'ppm_org_user_out_info', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 11, 0,  'ppm_pri_issue', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 12, 0,  'ppm_pri_issue_detail', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 13, 0,  'ppm_pri_issue_relation', 1000, 1000, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 14, 0,  'ppm_pro_project', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 15, 0,  'ppm_pro_project_member', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 16, 0,  'ppm_prs_priority', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 17, 0,  'ppm_prs_process', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 18, 0,  'ppm_prs_process_process_status', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 19, 0,  'ppm_prs_process_status', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 20, 0,  'ppm_prs_process_step', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 21, 0,  'ppm_prs_project_object_type', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 22, 0,  'ppm_prs_project_object_type_process', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 23, 0,  'ppm_prs_project_type_project_object_type', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 24, 0,  'ppm_res_resource', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 25, 0,  'ppm_rol_operation', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 26, 0,  'ppm_rol_permission', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 27, 0,  'ppm_rol_permission_operation', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 28, 0,  'ppm_rol_role', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 29, 0,  'ppm_rol_role_group', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 30, 0,  'ppm_rol_role_permission_operation', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 31, 0,  'ppm_rol_role_user', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 32, 0,  'ppm_sys_config', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 33, 0,  'ppm_tak_message', 1000, 500, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 34, 0,  'ppm_tem_team', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 35, 0,  'ppm_tem_user_team', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 36, 0,  'ppm_tre_trends', 1000, 2000, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 37, 0,  'ppm_tre_comment', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 38, 0,  'ppm_prs_issue_source', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 39, 0,  'ppm_pri_iteration', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 40, 0,  'ppm_pro_project_version', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 41, 0,  'ppm_pro_project_module', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 42, 0,  'ppm_sys_message_config', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 43, 0,  'ppm_org_user_config', 1000, 200, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 44, 0,  'ppm_prs_project_type', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 45, 0,  'ppm_bas_change_log', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 46, 0,  'ppm_tst_case', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 47, 0,  'ppm_tst_case_attachment', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 48, 0,  'ppm_tst_case_detail', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 49, 0,  'ppm_tst_case_group', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 50, 0,  'ppm_tst_case_step', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 51, 0,  'ppm_tst_plan', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 52, 0,  'ppm_tst_plan_case', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 53, 0,  'ppm_tst_plan_case_issue', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 54, 0,  'ppm_tst_plan_case_step', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 55, 0,  'ppm_pro_application', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 56, 0,  'ppm_pro_application_version', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 57, 0,  'ppm_orc_config', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_object_id( `id`, `org_id`, `code`, `max_id`, `step`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 58, 0,  'ppm_orc_message_config', 1000, 100, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (1, 'View', '查看', '查看对象', 1    , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (2, 'Modify', '更新', '更新对象', 1  , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (3, 'Delete', '删除', '删除对象', 1  , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (4, 'Create', '创建', '创建对象', 1  , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (5, 'Check', '审核', '审核对象', 1   , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (6, 'Invite', '邀请', '邀请对象', 1  , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (7, 'Bind', '绑定', '绑定对象', 1    , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (8, 'Unbind', '解绑', '解绑对象', 1  , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (9, 'Attention', '关注', '关注对象', 1       , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (10, 'UnAttention', '解除关注', '解除关注对象', 1    , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (11, 'ModifyStatus', '更新状态', '更新对象状态（流程流转）', 1       , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (12, 'Comment', '评论', '评论对象', 1        , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (13, 'Transfer', '转让', '转让对象', 1       , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (14, 'Init', '初始化', '初始化对象', 1       , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (15, 'Drop', '注销', '注销组织', 1   , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (16, 'Filing', '归档', '归档对象', 1 , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status` , `creator`, `create_time`, `updator`, `update_time`, `version` , `is_delete`)
    VALUES (17, 'UnFiling', '解除归档', '解除归档对象', 1       , 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 1, 1, 'Permission.Sys.Sys', 'Sys', '系统管理', 0, 1, '/Sys', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 2, 1, 'Permission.Sys.Dic', 'Dic', '数据字典管理', 1, 1, '/Sys/Dic', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 3, 1, 'Permission.Sys.Source', 'Source', '来源渠道管理', 1, 1, '/Sys/Source', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 4, 1, 'Permission.Sys.PayLevel', 'PayLevel', '订购级别管理', 1, 1, '/Sys/PayLevel', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 1, 1, 2, 'PermissionOperation.Sys.Dic.View', '查看数据字典', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 2, 1, 2, 'PermissionOperation.Sys.Dic.Create', '创建数据字典', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 3, 1, 2, 'PermissionOperation.Sys.Dic.Modify', '编辑数据字典', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 4, 1, 2, 'PermissionOperation.Sys.Dic.Delete', '删除数据字典', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 5, 1, 3, 'PermissionOperation.Sys.Source.View', '查看来源渠道', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 6, 1, 3, 'PermissionOperation.Sys.Source.Create', '创建来源渠道', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 7, 1, 3, 'PermissionOperation.Sys.Source.Modify', '编辑来源渠道', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 8, 1, 3, 'PermissionOperation.Sys.Source.Delete', '删除来源渠道', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 9, 1, 4, 'PermissionOperation.Sys.PayLevel.View', '查看订购级别', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 10, 1, 4, 'PermissionOperation.Sys.PayLevel.Create', '创建订购级别', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 11, 1, 4, 'PermissionOperation.Sys.PayLevel.Modify', '编辑订购级别', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 12, 1, 4, 'PermissionOperation.Sys.PayLevel.Delete', '删除订购级别', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_group( `id`, `org_id`, `lang_code`, `name`, `remark`, `type`, `is_readonly`, `is_show`, `is_default`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 1, 'RoleGroup.Sys', '系统角色组', '', 1, 1, 1, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 1, 'Role.Sys.Admin', '系统超级管理员', '', 1, 2, 2, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 1, 'Role.Sys.Manager', '系统管理员', '', 1, 1, 2, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 1, 'Role.Sys.Member', '系统成员', '', 1, 1, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 1, 1, 0, 1, '/Sys', '*', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 1, 2, 0, 2, '/Sys/Dic', '(View)|(Create)|(Modify)', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 1, 2, 0, 3, '/Sys/Source', '(View)|(Create)|(Modify)', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 1, 2, 0, 4, '/Sys/PayLevel', 'View', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 1, 3, 0, 2, '/Sys/Dic', 'View', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 1, 3, 0, 3, '/Sys/Source', 'View', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 1, 3, 0, 4, '/Sys/PayLevel', 'View', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_tre_trends( `id`, `org_id`, `uuid`, `module1`, `module2_id`, `module2`, `module3_id`, `module3`, `oper_code`, `oper_obj_id`, `oper_obj_type`, `oper_obj_property`, `relation_obj_id`, `relation_obj_type`, `relation_type`, `new_value`, `old_value`, `ext`, `creator`, `create_time`, `is_delete` )
        VALUES ( 1, 1, '156b9addb527bb91a7c95cd0e6f526f9', 'Sys', 0, '', 0, '', 'Init', 1, 'Sys', '', 0, '', '', '', '', '', 1, now(), 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 'ProcessStatus.Project.NotStart', '未开始', 1, '', '', 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 'ProcessStatus.Project.Running', '进行中', 2, '', '', 2, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 'ProcessStatus.Project.Complete', '已完成', 3, '', '', 3, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 'ProcessStatus.Iteration.NotStart', '未开始', 1, '', '', 1, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 'ProcessStatus.Iteration.Running', '进行中', 2, '', '', 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 'ProcessStatus.Iteration.Complete', '已完成', 3, '', '', 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 'ProcessStatus.Issue.NotStart', '未完成', 1, '', '', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 'ProcessStatus.Issue.WaitEvaluate', '待评估', 2, '', '', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0, 'ProcessStatus.Issue.WaitConfirmBug', '待确认缺陷', 3, '', '', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 10, 0, 'ProcessStatus.Issue.ConfirmedBug', '已确认缺陷', 4, '', '', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 11, 0, 'ProcessStatus.Issue.ReOpen', '重新打开', 5, '', '', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 12, 0, 'ProcessStatus.Issue.Evaluated', '已评估', 6, '', '', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 13, 0, 'ProcessStatus.Issue.Planning', '计划中', 7, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 14, 0, 'ProcessStatus.Issue.Design', '设计中', 8, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 15, 0, 'ProcessStatus.Issue.Processing', '处理中', 9, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 16, 0, 'ProcessStatus.Issue.Development', '研发中', 10, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 17, 0, 'ProcessStatus.Issue.WaitTest', '待测试', 11, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 18, 0, 'ProcessStatus.Issue.Testing', '测试中', 12, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 19, 0, 'ProcessStatus.Issue.WaitRelease', '待发布', 13, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 20, 0, 'ProcessStatus.Issue.Wait', '挂起', 14, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 21, 0, 'ProcessStatus.Issue.Repair', '修复中', 15, '', '', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 22, 0, 'ProcessStatus.Issue.Success', '成功', 16, '', '', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 23, 0, 'ProcessStatus.Issue.Fail', '失败', 17, '', '', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 24, 0, 'ProcessStatus.Issue.Released', '已发布', 18, '', '', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 25, 0, 'ProcessStatus.Issue.Closed', '已关闭', 19, '', '', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 26, 0, 'ProcessStatus.Issue.Complete', '已完成', 20, '', '', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 27, 0, 'ProcessStatus.Issue.Confirmed', '已确认', 21, '', '', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 'Process.DefaultProject', '默认项目流程', 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 'Process.DefaultIteration', '默认迭代流程', 1, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 'Process.Issue.DefaultFeature', '默认特性流程', 2, 3, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 'Process.Issue.DefaultDemand', '默认需求流程', 2, 3, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 'Process.Issue.DefaultTask', '默认任务流程', 1, 3, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 'Process.Issue.DefaultAgileTask', '默认敏捷项目任务流程', 2, 3, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 'Process.Issue.DefaultBug', '默认缺陷流程', 2, 3, 5, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 'Process.Issue.DefaultTestTask', '默认测试任务流程', 2, 3, 6, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 1, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 1, 2, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 1, 3, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 2, 4, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 2, 5, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 2, 6, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 3, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 3, 13, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0, 3, 12, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 10, 0, 3, 16, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 11, 0, 3, 26, 2, 6, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 12, 0, 3, 25, 2, 7, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 13, 0, 4, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 14, 0, 4, 13, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 15, 0, 4, 14, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 16, 0, 4, 16, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 17, 0, 4, 17, 2, 5, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 18, 0, 4, 18, 2, 6, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 19, 0, 4, 19, 2, 7, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 20, 0, 4, 24, 2, 8, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 21, 0, 4, 25, 2, 9, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 22, 0, 5, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 23, 0, 5, 15, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 24, 0, 5, 26, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 25, 0, 6, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 26, 0, 6, 15, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 27, 0, 6, 20, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 28, 0, 6, 26, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 29, 0, 7, 9, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 30, 0, 7, 10, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 31, 0, 7, 21, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 32, 0, 7, 17, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 33, 0, 7, 11, 2, 5, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 34, 0, 7, 19, 2, 6, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 35, 0, 7, 24, 2, 7, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 36, 0, 7, 25, 2, 8, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 37, 0, 8, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 38, 0, 8, 11, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 39, 0, 8, 20, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 40, 0, 8, 22, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 41, 0, 8, 23, 2, 5, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 1, 'ProcessStep.Project.Start', '项目开始', 1, 2, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 1, 'ProcessStep.Project.Complete', '项目完成', 1, 3, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 1, 'ProcessStep.Project.Complete', '项目完成', 2, 3, 3, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 1, 'ProcessStep.Project.NotStart', '项目未开始', 2, 1, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 1, 'ProcessStep.Project.ReStart', '重启项目', 3, 1, 5, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 2, 'ProcessStep.Iteration.Start', '迭代开始', 4, 5, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 2, 'ProcessStep.Iteration.Complete', '迭代完成', 4, 6, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 2, 'ProcessStep.Iteration.Complete', '迭代完成', 5, 6, 3, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0, 2, 'ProcessStep.Iteration.NotStart', '迭代未开始', 5, 4, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 10, 0, 2, 'ProcessStep.Iteration.ReStart', '重启迭代', 6, 4, 5, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 11, 0, 3, 'ProcessStep.Feature.StartPlan', '开始计划', 7, 13, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 12, 0, 3, 'ProcessStep.Feature.Evaluated', '已评估', 7, 12, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 13, 0, 3, 'ProcessStep.Feature.StartDevelopment', '进入研发', 7, 16, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 14, 0, 3, 'ProcessStep.Feature.Released', '发布', 7, 26, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 15, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 7, 25, 5, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 16, 0, 3, 'ProcessStep.Feature.Evaluated', '已评估', 13, 12, 6, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 17, 0, 3, 'ProcessStep.Feature.StartDevelopment', '进入研发', 13, 16, 7, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 18, 0, 3, 'ProcessStep.Feature.Released', '发布', 13, 26, 8, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 19, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 13, 25, 9, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 20, 0, 3, 'ProcessStep.Feature.NotStart', '未开始', 13, 7, 10, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 21, 0, 3, 'ProcessStep.Feature.StartDevelopment', '进入研发', 12, 16, 11, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 22, 0, 3, 'ProcessStep.Feature.Released', '发布', 12, 26, 12, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 23, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 12, 25, 13, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 24, 0, 3, 'ProcessStep.Feature.NotStart', '未开始', 12, 7, 14, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 25, 0, 3, 'ProcessStep.Feature.ReStartPlan', '重新计划', 12, 13, 15, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 26, 0, 3, 'ProcessStep.Feature.Released', '发布', 16, 26, 16, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 27, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 16, 25, 17, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 28, 0, 3, 'ProcessStep.Feature.NotStart', '未开始', 16, 7, 18, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 29, 0, 3, 'ProcessStep.Feature.ReStartPlan', '重新计划', 16, 13, 19, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 30, 0, 3, 'ProcessStep.Feature.Evaluated', '已评估', 16, 12, 20, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 31, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 26, 25, 21, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 32, 0, 3, 'ProcessStep.Feature.NotStart', '未开始', 26, 7, 22, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 33, 0, 3, 'ProcessStep.Feature.ReStartPlan', '重新计划', 26, 13, 23, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 34, 0, 3, 'ProcessStep.Feature.Evaluated', '已评估', 26, 12, 24, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 35, 0, 3, 'ProcessStep.Feature.StartDevelopment', '进入研发', 26, 16, 25, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 36, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 7, 13, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 37, 0, 4, 'ProcessStep.Demand.Evaluated', '已评估', 7, 14, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 38, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 7, 16, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 39, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 7, 25, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 40, 0, 4, 'ProcessStep.Demand.Evaluated', '开始设计', 13, 14, 5, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 41, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 13, 16, 6, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 42, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 13, 25, 7, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 43, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 13, 7, 8, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 44, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 14, 16, 9, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 45, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 14, 25, 10, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 46, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 14, 7, 11, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 47, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 14, 13, 12, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 48, 0, 4, 'ProcessStep.Demand.WaitTest', '等待测试', 16, 17, 13, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 49, 0, 4, 'ProcessStep.Demand.StartTest', '进入测试', 16, 18, 14, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 50, 0, 4, 'ProcessStep.Demand.WaitReleased', '待发布', 16, 19, 15, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 51, 0, 4, 'ProcessStep.Demand.Released', '发布', 16, 24, 16, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 52, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 16, 25, 17, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 53, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 16, 7, 18, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 54, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 16, 13, 19, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 55, 0, 4, 'ProcessStep.Demand.StartDesign', '开始设计', 16, 14, 20, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 56, 0, 4, 'ProcessStep.Demand.StartTest', '进入测试', 17, 18, 21, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 57, 0, 4, 'ProcessStep.Demand.WaitReleased', '待发布', 17, 19, 22, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 58, 0, 4, 'ProcessStep.Demand.Released', '发布', 17, 24, 23, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 59, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 17, 25, 24, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 60, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 17, 7, 25, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 61, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 17, 13, 26, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 62, 0, 4, 'ProcessStep.Demand.StartDesign', '开始设计', 17, 14, 27, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 63, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 17, 16, 28, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 64, 0, 4, 'ProcessStep.Demand.WaitReleased', '待发布', 18, 19, 29, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 65, 0, 4, 'ProcessStep.Demand.Released', '发布', 18, 24, 30, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 66, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 18, 25, 31, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 67, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 18, 7, 32, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 68, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 18, 13, 33, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 69, 0, 4, 'ProcessStep.Demand.StartDesign', '开始设计', 18, 14, 34, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 70, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 18, 16, 35, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 71, 0, 4, 'ProcessStep.Demand.WaitTest', '等待测试', 18, 17, 36, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 72, 0, 4, 'ProcessStep.Demand.Released', '发布', 19, 24, 37, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 73, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 19, 25, 38, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 74, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 19, 16, 39, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 75, 0, 4, 'ProcessStep.Demand.WaitTest', '等待测试', 19, 17, 40, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 76, 0, 4, 'ProcessStep.Demand.StartTest', '进入测试', 19, 18, 41, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 77, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 24, 25, 42, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 78, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 24, 7, 43, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 79, 0, 4, 'ProcessStep.Demand.WaitReleased', '待发布', 24, 19, 44, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 80, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 25, 7, 45, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 81, 0, 5, 'ProcessStep.Task.Start', '任务开始', 7, 15, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 82, 0, 5, 'ProcessStep.Task.Complete', '任务完成', 7, 26, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 83, 0, 5, 'ProcessStep.Task.Complete', '任务完成', 15, 26, 3, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 84, 0, 5, 'ProcessStep.Task.NotStart', '任务未开始', 15, 7, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 85, 0, 5, 'ProcessStep.Task.ReStart', '重启任务', 26, 7, 5, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 86, 0, 6, 'ProcessStep.AgileTask.Start', '任务开始', 7, 15, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 87, 0, 6, 'ProcessStep.AgileTask.HangUp', '任务挂起', 7, 20, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 88, 0, 6, 'ProcessStep.AgileTask.Complete', '任务完成', 7, 26, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 89, 0, 6, 'ProcessStep.AgileTask.Complete', '任务完成', 15, 26, 4, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 90, 0, 6, 'ProcessStep.AgileTask.HangUp', '任务挂起', 15, 20, 5, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 91, 0, 6, 'ProcessStep.AgileTask.NotStart', '任务未开始', 15, 7, 6, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 92, 0, 6, 'ProcessStep.AgileTask.ReStart', '重启任务', 20, 7, 7, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 93, 0, 6, 'ProcessStep.AgileTask.ReStart', '重启任务', 26, 7, 8, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 94, 0, 7, 'ProcessStep.Bug.Confirm', '确认', 9, 10, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 95, 0, 7, 'ProcessStep.Bug.Repair', '修复', 9, 21, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 96, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 9, 17, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 97, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 9, 19, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 98, 0, 7, 'ProcessStep.Bug.Released', '已发布', 9, 24, 5, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 99, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 9, 25, 6, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 100, 0, 7, 'ProcessStep.Bug.Repair', '修复', 10, 21, 7, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 101, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 10, 17, 8, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 102, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 10, 19, 9, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 103, 0, 7, 'ProcessStep.Bug.Released', '已发布', 10, 24, 10, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 104, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 10, 25, 11, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 105, 0, 7, 'ProcessStep.Bug.WaitConfirm', '待确认', 10, 9, 12, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 106, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 21, 17, 13, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 107, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 21, 19, 14, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 108, 0, 7, 'ProcessStep.Bug.Released', '已发布', 21, 24, 15, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 109, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 21, 25, 16, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 110, 0, 7, 'ProcessStep.Bug.Confirm', '确认', 21, 10, 17, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 111, 0, 7, 'ProcessStep.Bug.ReOpen', '待发布', 17, 19, 18, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 112, 0, 7, 'ProcessStep.Bug.WaitReleased', '重新打开', 17, 11, 19, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 113, 0, 7, 'ProcessStep.Bug.Released', '已发布', 17, 24, 20, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 114, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 17, 25, 21, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 115, 0, 7, 'ProcessStep.Bug.Repair', '修复', 17, 21, 22, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 116, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 11, 21, 23, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 117, 0, 7, 'ProcessStep.Bug.Released', '待测试', 11, 17, 24, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 118, 0, 7, 'ProcessStep.Bug.Closed', '已发布', 11, 24, 25, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 119, 0, 7, 'ProcessStep.Bug.WaitConfirm', '关闭', 11, 25, 26, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 120, 0, 7, 'ProcessStep.Bug.Repair', '确认', 11, 10, 27, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 121, 0, 7, 'ProcessStep.Bug.Released', '已发布', 19, 24, 28, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 122, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 19, 25, 29, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 123, 0, 7, 'ProcessStep.Bug.Repair', '修复', 19, 21, 30, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 124, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 19, 17, 31, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 125, 0, 7, 'ProcessStep.Bug.ReOpen', '重新打开', 19, 11, 32, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 126, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 24, 25, 33, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 127, 0, 7, 'ProcessStep.Bug.Repair', '修复', 24, 21, 34, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 128, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 24, 17, 35, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 129, 0, 7, 'ProcessStep.Bug.ReOpen', '重新打开', 24, 11, 36, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 130, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 24, 19, 37, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 131, 0, 7, 'ProcessStep.Bug.ReOpen', '重新打开', 25, 11, 38, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 132, 0, 8, 'ProcessStep.TestTask.Success', '成功', 7, 22, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 133, 0, 8, 'ProcessStep.TestTask.Fail', '失败', 7, 23, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 134, 0, 8, 'ProcessStep.TestTask.Wait', '挂起', 7, 20, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 135, 0, 8, 'ProcessStep.TestTask.Fail', '失败', 22, 23, 4, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 136, 0, 8, 'ProcessStep.TestTask.Wait', '挂起', 22, 20, 5, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 137, 0, 8, 'ProcessStep.TestTask.ReOpen', '重新打开', 22, 11, 6, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 138, 0, 8, 'ProcessStep.TestTask.NotStart', '未开始', 22, 7, 7, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 139, 0, 8, 'ProcessStep.TestTask.Success', '成功', 23, 22, 8, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 140, 0, 8, 'ProcessStep.TestTask.Wait', '挂起', 23, 20, 9, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 141, 0, 8, 'ProcessStep.TestTask.ReOpen', '重新打开', 23, 11, 10, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 142, 0, 8, 'ProcessStep.TestTask.NotStart', '未开始', 23, 7, 11, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 143, 0, 8, 'ProcessStep.TestTask.ReOpen', '重新打开', 20, 11, 12, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 144, 0, 8, 'ProcessStep.TestTask.Success', '成功', 20, 22, 13, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 145, 0, 8, 'ProcessStep.TestTask.Fail', '失败', 20, 23, 14, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 146, 0, 8, 'ProcessStep.TestTask.Success', '成功', 11, 22, 15, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 147, 0, 8, 'ProcessStep.TestTask.Fail', '失败', 11, 23, 16, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 148, 0, 8, 'ProcessStep.TestTask.Wait', '挂起', 11, 20, 17, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type( `id`, `org_id`, `lang_code`, `name`, `sort`, `default_process_id`, `category`, `mode`, `is_readonly`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 'ProjectType.NormalTask', '普通任务项目', 1, 1, 0, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type( `id`, `org_id`, `lang_code`, `name`, `sort`, `default_process_id`, `category`, `mode`, `is_readonly`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 'ProjectType.Agile', '敏捷研发项目', 2, 1, 0, 2, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 'Project.ObjectType.Iteration', 'I', '迭代', 1, '', '', '', 1, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 'Project.ObjectType.Feature', 'F', '特性', 2, '', '', '', 2, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 'Project.ObjectType.Demand', 'D', '需求', 2, '', '', '', 3, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 'Project.ObjectType.Task', 'T', '任务', 2, '', '', '', 4, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 'Project.ObjectType.Bug', 'B', '缺陷', 2, '', '', '', 5, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 'Project.ObjectType.TestTask', 'TT', '测试任务', 2, '', '', '', 6, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 1, 0, '', 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 1, 4, '', 5, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 2, 0, '', 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 2, 1, '', 2, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 2, 2, '', 3, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 2, 3, '', 4, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 2, 4, '', 6, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 2, 5, '', 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0, 2, 6, '', 8, 1, 1, 1, now(), 1, now(), 1, 2);