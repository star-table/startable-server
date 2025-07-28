INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 166, 0, 21, 'Permission.Pro.Config.ModifyField', ' 新增/管理自定义字段', 'ModifyField', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 167, 0, 12, 'Permission.Org.Project.ModifyField', ' 新增/管理自定义字段', 'ModifyField', '', 1, 1, now(), 1, now(), 1, 2);

INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (27, 'ModifyField', '管理字段', '管理字段', 1, 1, now(), 1, now(), 1, 2);
