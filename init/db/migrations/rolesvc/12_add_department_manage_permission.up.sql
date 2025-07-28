INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (25, 'Watch', '查看', '查看', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (26, 'ModifyDepartment', '编辑部门用户', '编辑部门用户', 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 42, 0, 'Permission.Org.Department', 'Department', '组织架构管理', 5, 2, '/Org/{org_id}/Department', '', 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 161, 0, 42, 'Permission.Org.Department.Create', '创建部门', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 162, 0, 42, 'Permission.Org.Department.Modify', '编辑部门', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 163, 0, 42, 'Permission.Org.Department.Delete', '删除部门', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 164, 0, 8, 'Permission.Org.User.Watch', '查看成员列表', 'Watch', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 165, 0, 8, 'Permission.Org.User.ModifyDepartment', '修改组织成员部门', 'ModifyDepartment', '', 1, 1, now(), 1, now(), 1, 2);
